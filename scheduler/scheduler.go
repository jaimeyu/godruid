package scheduler

import (
	"fmt"
	"strconv"

	"sync"
	"sync/atomic"

	metrics "github.com/accedian/adh-gather/models/metrics"
	"github.com/getlantern/deepcopy"
	"github.com/robfig/cron"

	//adhh "github.com/accedian/adh-gather/handlers"
	"time"

	"github.com/accedian/adh-gather/datastore/couchDB"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/models"
)

// SLAConfig - container for the ReportScheduleConfig and its associated job
type SLAConfig struct {
	Config *metrics.ReportScheduleConfig
	Job    func()
}

/* Structure to pass to the worker pool through a channel
 * request - This is the report request that will get sent to druid
 * result  - This is a channel where the worker will send its results to
 */
type workerSubmitChan struct {
	request func() //*metrics.SLAReportRequest
}

/* This structure holds all the necessary info for the scheduler
 */
type schedulercfg struct {
	// HTTP Requests are handle in goroutines which means we need to additions/deletions atomic
	mux sync.Mutex

	// cron doesn't have a delete function so we destroy it
	// and overwrite with a new one when we add/delete jobs
	curCron *cron.Cron

	// This is the metrics service handler
	// TODO: Maybe we shouldn't use functions in the msh and
	// create a new 'metrics scheduled handler'
	msh metrics.MetricServiceHandler

	// DB access to stored schedules.  For now, this is actually the tenantDB
	db metrics.ScheduleDB

	// This is an interim solution because I need to figure out a way to
	// get all the current tenants in order to drill down and get all the jobs per tenant.
	adminDb metrics.AdminInterface

	// How much time to buffer operations before timing out.
	// @TODO Not used in logic!
	gracePeriod int32

	// Maximum workers set when struct is initialized
	curMaxWorkers uint

	// worker's job channel
	// Workers listen to this channel to pull jobs in
	// Note that workers do NOT send data into this channel
	druidRequestsQ chan workerSubmitChan

	cronRestarts int32
}

// Maximum number of workers that can access the druid db concurrently
// TODO We don't know what the limits are right now, 25 is just a random number
const maxWorkers = 25

// Scheduler restart backoff
const retryDelay = 30

// Max restarts
const maxRestarts = 30

// Let jobs wait for 10 minutes before the job gets dropped for being late.
const gracePeriod = 10 * 60

// Minimize package's globals
var schedulecfg schedulercfg

// Takes in a metrics.ReportScheduleConfig and converts into a SLAReportRequest
func convertScheduleConfigToRequest(s metrics.ReportScheduleConfig) (*metrics.SLAReportRequest, error) {
	request := metrics.SLAReportRequest{}
	request.TenantID = s.TenantID
	request.ThresholdProfileID = s.ThresholdProfile
	request.Meta = s.Meta
	request.Granularity = s.Granularity
	request.Timeout = s.Timeout
	request.Timezone = "UTC"
	request.SLAScheduleConfig = s.ID

	// iso interval yyyymmddThhmmssfff/yyyymmddThhmmssfff
	// see http://support.sas.com/documentation/cdl/en/lrdict/64316/HTML/default/viewer.htm#a003169814.htm

	today := time.Now().UTC()
	tz, err := time.LoadLocation("UTC")
	if err != nil {
		msg := fmt.Sprintf("Could not load timezone: %s. Error: %s", models.AsJSONString(s), err.Error())
		logger.Log.Errorf(msg)
		return nil, err
	}

	h, err := strconv.Atoi(s.Hour)
	m, err := strconv.Atoi(s.Minute)

	sec, err := strconv.Atoi(s.Second)

	timeEnd := time.Date(today.Year(), today.Month(), today.Day(), h, m, sec, 0, tz)

	// @TODO Do we know of a ISO8601 validator for the Duration notation? Right now, we're just copying it in wholesale.
	// Here is a regex example: https://stackoverflow.com/a/30592146 but should we do it ourselves...
	duration := s.TimeRangeDuration
	timeStart := duration

	// Cheat by using RFC3339, see https://stackoverflow.com/questions/522251/whats-the-difference-between-iso-8601-and-rfc-3339-date-formats
	// As long as we don't use interval or period specific notation, we should be fine
	isoTimeEnd := timeEnd.Format(time.RFC3339)

	intervalStr := fmt.Sprintf("%s/%s", timeStart, isoTimeEnd)

	msg := fmt.Sprintf("Interval for the druid query is: %s ", intervalStr)
	logger.Log.Debugf(msg)
	request.Interval = intervalStr

	return &request, nil
}

func pendWork(request metrics.SLAReportRequest) error {
	timeStart := time.Now()
	// Create a channel for when the report is returned on
	result := make(chan *metrics.SLAReport)
	w := func() {
		timeEnd := time.Now()

		delay := timeEnd.Sub(timeStart)
		if delay > time.Duration(schedulecfg.gracePeriod)*time.Second {
			logger.Log.Errorf("Could not run the job in a timely manner, dropping it")
			result <- nil
			return
		}

		// Now get the Report
		report, err := schedulecfg.msh.GetInternalSLAReportV1(&request)
		if err != nil {
			msg := fmt.Sprintf("Unable to get Scheduled SLA Report Configuration: %s. Error: %s", models.AsJSONString(request), err.Error())
			logger.Log.Errorf(msg)
			result <- nil
			return
		}
		logger.Log.Debugf("REPORT: %+v", report)

		result <- report
	}

	work := workerSubmitChan{
		request: w,
	}

	schedulecfg.druidRequestsQ <- work
	if request.Timeout < 5000 {
		request.Timeout = 5000
	}
	t := time.Duration(request.Timeout)

	timeout := t * time.Millisecond

	select {
	case r := <-result:
		if r == nil {
			msg := fmt.Sprintf("Received nil report for request %s", models.AsJSONString(request))
			logger.Log.Errorf(msg)
			break
		} else {
			logger.Log.Debugf("Report received! ", models.AsJSONString(r))
			_, err := schedulecfg.db.CreateSLAReport(r)
			if err != nil {
				logger.Log.Errorf("Not store SLA Report: %s", err)
			}

			break
		}
	case <-time.After(time.Duration(timeout)):
		logger.Log.Errorf("Report request timed out %s", models.AsJSONString(request))
	}
	return nil
}

/*createWork - Creates the function that converts a
  SLA Scheduled Config into SLA Report Request format.
  The reason is the Report Request requires a time range and we need to generate them.
*/
func createWork(origSchedConfig *metrics.ReportScheduleConfig) (func(), error) {
	var s metrics.ReportScheduleConfig
	// Copy the file
	err := deepcopy.Copy(&s, origSchedConfig)
	if err != nil {
		msg := fmt.Sprintf("Could not deep copy the SLA schedule configuration %s. Error: %s", models.AsJSONString(s), err.Error())
		logger.Log.Errorf(msg)
		return nil, err
	}
	logger.Log.Debugf("Creating work for %s", s.Name)

	functor := func() {
		logger.Log.Debugf("Job %s's Work is exec'ing", s.Name)
		request, err := convertScheduleConfigToRequest(s)
		if err != nil {
			msg := fmt.Sprintf("Unable to convert Scheduled SLA Report Configuration: %s. Error: %s", models.AsJSONString(s), err.Error())
			logger.Log.Errorf(msg)
			return
		}

		pendWork(*request)

	}

	return functor, nil
}

/*addJob - Adds a job to the package's job list.
 * After the job is successfully added, it restarts the cron scheduler
 * so it will have up to date jobs.
 * Note that the function is thread safe. I expect callers to be from goroutines.
 */
func addJob(config SLAConfig) error {
	var err error
	s := config.Config
	logger.Log.Infof("Adding job %s", models.AsJSONString(s))

	if len(s.Second) == 0 {
		logger.Log.Debugf("Seconds was empty, setting to 0")
		s.Second = "0"
	}
	if _, err := strconv.Atoi(s.Second); err != nil {
		logger.Log.Debugf("Seconds was not a valid character, setting to 0")
		if s.Second != "*" {
			s.Second = "0"
		}
	}

	spec := fmt.Sprintf("%s %s %s %s %s %s", s.Second, s.Minute, s.Hour, s.DayOfMonth, s.Month, s.DayOfWeek)
	//spec := fmt.Sprintf("0 * * * * *") // runs every 1 minute schedule
	//spec := fmt.Sprintf("0 0 3 * * 0") // runs every week on sunday at 3am EDT
	//loc := time.LoadLocation("Europe/London")
	logger.Log.Infof("Setting cron job %s spec to '%s'", s.Name, spec)

	config.Job, err = createWork(config.Config)
	if err != nil {
		return err
	}
	err = schedulecfg.curCron.AddFunc(spec, config.Job)
	if err != nil {
		logger.Log.Errorf("Could not create cron job: '%s'", err)

		return err
	}

	return nil
}

// dbgDumpCronJobs - Dumps the cron jobs
func dbgDumpCronJobs() {
	logger.Log.Debugf("Dumping cron jobs")
	for i, j := range GetScheduledJobs() {
		logger.Log.Debugf("%d sched:%+v", i, j)
	}

}

// AddJob - External caller uses this to add in Jobs to the scheduler
func AddJob(config SLAConfig) error {
	err := addJob(config)
	if err != nil {
		return err
	}

	// If externally called, restart CRON
	RebuildCronJobs()
	return nil
}

/*RemoveJob - Remove a job to the package's job list.
 *  After the job is successfully removed, it restarts the cron scheduler
 *  so it will have up to date jobs.
 *  Note that the function is thread safe. I expect callers to be from goroutines.
 */
func RemoveJob(config *metrics.ReportScheduleConfig) error {

	//	logger.Log.Debugf("Removing job:%s", config.Name)
	//	result, err := schedulecfg.db.DeleteReportScheduleConfig(config.TenantID, config.ID)
	//	if err != nil {
	//		logger.Log.Errorf("Could not delete job: %s ", err)
	//		return result, err
	//	}

	RebuildCronJobs()
	return nil
}

//getAllStoredJobsFromDB - Gets all the stored jobs from couchdb and returns it
func getAllStoredJobsFromDB() ([]*metrics.ReportScheduleConfig, error) {

	logger.Log.Debugf("Getting stored jobs from db")
	tenants, err := schedulecfg.adminDb.GetAllTenantDescriptors()
	if err != nil {
		logger.Log.Errorf("Could not get all known tenants")
	}
	logger.Log.Debugf("Got %d tenants", len(tenants))

	var dbJobs []*metrics.ReportScheduleConfig
	for _, tenant := range tenants {
		logger.Log.Debugf("Getting ReportConfigs from tenant %s", tenant.ID)
		jobs, err := schedulecfg.db.GetAllReportScheduleConfigs(tenant.ID)
		logger.Log.Debugf("%s jobs:%+v", tenant.ID, jobs)
		if err != nil {
			logger.Log.Errorf("Could not access tenant %s' db", tenant.ID)
		} else {
			dbJobs = append(dbJobs, jobs...)
			logger.Log.Debugf("%s dbjobs:%+v", tenant.ID, dbJobs)
		}
	}
	logger.Log.Debugf("Got all jobs dbjobs:%+v", dbJobs)
	// Mask the error... we get a 404 error if there are no configs in the tenant's table which is valid.
	return dbJobs, nil
}

// rebuildCronJobsHelper - Cron cannot start, this will help get the cron jobs to restart
func rebuildCronJobsHelper() {

	atomic.AddInt32(&schedulecfg.cronRestarts, 1)
	if atomic.LoadInt32(&schedulecfg.cronRestarts) > maxRestarts {
		logger.Log.Fatal("Scheduler could not rebuild the schedules due to lack of access to DB")
	}

	RebuildCronJobs()
}

// RebuildCronJobs - Restarts the cron plugin so it can repopulate its job list
// This is a work around until the cron library supports 'remove' semantics to
// simplify updating jobs.
func RebuildCronJobs() error {

	schedulecfg.mux.Lock()
	defer schedulecfg.mux.Unlock()

	logger.Log.Debugf("Scheduler is rebuilding CRON jobs")

	configs := []*metrics.ReportScheduleConfig{}
	var err error

	logger.Log.Debugf("Getting stored jobs")
	configs, err = getAllStoredJobsFromDB()
	if err != nil {
		logger.Log.Errorf("Could not start CRON because could not get stored scheduleds, retrying in %d seconds: %s ", retryDelay, err)
		// TODO Hey, should crap out or hammer the database more?
		time.AfterFunc(time.Second*retryDelay, rebuildCronJobsHelper)
		return err
	}
	atomic.StoreInt32(&schedulecfg.cronRestarts, 0)

	oldcron := schedulecfg.curCron
	loc := time.UTC
	schedulecfg.curCron = cron.NewWithLocation(loc)

	schedulecfg.curCron.Start()
	for _, config := range configs {
		req := SLAConfig{
			Config: config,
		}
		// Don't run jobs that are not active.
		if config.Active == true {
			err := addJob(req)
			if err != nil {
				logger.Log.Errorf("Could not add job")
				return err
			}
		} else {
			logger.Log.Debugf("Job %s is not active", req.Config.Name)
		}

	}

	if oldcron != nil {
		logger.Log.Debugf("Stopping old CRON service (new CRON service already started)")
		oldcron.Stop()
	}
	// Dump it for debugging
	dbgDumpCronJobs()
	return nil
}

// Deinit - Deinits the scheduler
func Deinit() {
	logger.Log.Debugf("De-init'ing Scheduler")
	schedulecfg.curCron.Stop()
	// Cleans up the workers
	close(schedulecfg.druidRequestsQ)
}

// Initialize Initialize the scheduler
func Initialize(m metrics.MetricServiceHandler, scheduleDB metrics.ScheduleDB, adminDB metrics.AdminInterface, workers uint) {
	logger.Log.Debugf("Scheduler is initializing")
	schedulecfg = schedulercfg{}
	schedulecfg.mux = sync.Mutex{}
	schedulecfg.msh = m
	schedulecfg.db = scheduleDB
	schedulecfg.adminDb = adminDB
	schedulecfg.gracePeriod = gracePeriod
	schedulecfg.cronRestarts = 0

	var dberr error

	if adminDB == nil {

		logger.Log.Debugf("ScheduleDB is nil, creating new instance")
		schedulecfg.db, dberr = couchDB.CreateTenantServiceDAO()
		// TODO How to handle fata error?
		if dberr != nil {
			logger.Log.Fatalf("Could not create tenant service")
		}
	}

	if scheduleDB == nil {
		logger.Log.Debugf("AdminDB is nil, creating new instance")
		schedulecfg.adminDb, dberr = couchDB.CreateAdminServiceDAO()
		// TODO How to handle fata error?
		if dberr != nil {
			logger.Log.Fatalf("Could not create admin service")
		}
	}

	if workers >= maxWorkers {
		workers = maxWorkers
	}
	schedulecfg.curMaxWorkers = workers
	logger.Log.Debugf("Set max workers to %d", workers)

	if schedulecfg.druidRequestsQ != nil {
		logger.Log.Debugf("Closing old worker channel")
		// If there is an existing channel, close it so the workers clean up
		close(schedulecfg.druidRequestsQ)
	}

	schedulecfg.druidRequestsQ = make(chan workerSubmitChan, workers*3)
	var i uint
	for i = 0; i < schedulecfg.curMaxWorkers; i++ {
		go worker(i, schedulecfg.druidRequestsQ)
	}

	RebuildCronJobs()

	logger.Log.Debugf("Scheduler is initialized")
}

// SLAGetSchedules - Gets the current list of cron jobs
func GetScheduledJobs() []*cron.Entry {
	return schedulecfg.curCron.Entries()
}

// Stop - Stops the current cron instance
func Stop() error {
	schedulecfg.curCron.Stop()
	return nil
}

// worker - the workers are goroutines that executes the druid jobs and waits for completion.
// We don't want to hit druid too hard so we can rate limit by limiting the number of concurrent workers
func worker(id uint, jobs <-chan workerSubmitChan) {

	logger.Log.Debugf("Report Worker #%d created and waiting for jobs", id)
	for j := range jobs {
		j.request()
	}
}
