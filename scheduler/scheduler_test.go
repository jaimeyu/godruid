package scheduler_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/accedian/adh-gather/logger"
	admmod "github.com/accedian/adh-gather/models/admin"
	metrics "github.com/accedian/adh-gather/models/metrics"
	"github.com/accedian/adh-gather/scheduler"
	"github.com/robfig/cron"
)

type mockDB struct {
	fn      func()
	configs []*metrics.ReportScheduleConfig
	err     error
}

func (m *mockDB) GetAllReportScheduleConfigs(tenantID string) ([]*metrics.ReportScheduleConfig, error) {
	logger.Log.Debugf("Getting stored mock jobs")

	for i, j := range m.configs {
		logger.Log.Debugf("{%d,%+v}", i, j)
	}

	return m.configs, m.err
}

func (m *mockDB) CreateSLAReport(report *metrics.SLAReport) (*metrics.SLAReport, error) {
	logger.Log.Debugf("Creating SLA REPORT")
	return report, m.err
}

type metricServiceHandler struct {
	report *metrics.SLAReport
	err    error
	result chan int
}

func (m *metricServiceHandler) GetInternalSLAReport(request *metrics.SLAReportRequest) (*metrics.SLAReport, error) {

	logger.Log.Debugf("Execute mock get a SLA Report from druid")
	// Send a 1 to signfy we actually executed the cron job
	m.result <- 1
	return m.report, m.err
}

func (m *mockDB) DeleteReportScheduleConfig(tenantID string, configID string) (*metrics.ReportScheduleConfig, error) {
	logger.Log.Debugf("Deleting mock get a SLA Report from duid")

	m.configs = m.configs[:len(m.configs)-1]

	return nil, nil
}

type mockAdminDB struct {
}

func (a *mockAdminDB) GetAllTenantDescriptors() ([]*admmod.Tenant, error) {
	var list []*admmod.Tenant
	list = append(list, &admmod.Tenant{})
	return list, nil
}

func TestSchedulerBasics(t *testing.T) {

	defer scheduler.Stop()

	c := make(chan int)
	logger.Log.Debug("Start TestSchedulerBasic")

	var configs []*metrics.ReportScheduleConfig

	mockdb := mockDB{
		err:     nil,
		configs: configs,
	}

	report := metrics.SLAReport{
		ID:                   "000001",
		ReportCompletionTime: "000010",
		TenantID:             "123456",
		ReportTimeRange:      "asdkj?",
		ReportSummary:        metrics.ReportSummary{},
		TimeSeriesResult:     nil,
		ByHourOfDayResult:    nil,
		ByDayOfWeekResult:    nil,
	}

	mockmsh := metricServiceHandler{
		err:    nil,
		report: &report,
		result: c,
	}

	mockadmindb := mockAdminDB{}

	scheduler.Initialize(&mockmsh, &mockdb, &mockadmindb, 1)

	var testPayload = metrics.ReportScheduleConfig{}

	// Execute every 1 second
	testPayload.DayOfWeek = "*"
	testPayload.DayOfMonth = "*"
	testPayload.Month = "*"
	testPayload.Hour = "*"
	testPayload.Minute = "*"
	testPayload.Second = "*"
	testPayload.Name = "Mock Job 1"
	testPayload.Timeout = 5000
	testPayload.TimeRangeDuration = "P1Y2M10DT2H30M"
	testPayload.Active = true

	var passed bool = false

	job := scheduler.SLAConfig{
		Config: &testPayload,
	}

	mockdb.configs = append(mockdb.configs, &testPayload)
	for _, cf := range configs {
		logger.Log.Debug("Dumping configs:", cf)
	}

	err := scheduler.AddJob(job)
	if err != nil {
		t.FailNow()
	}

	endTest := func() {
		logger.Log.Debug("Test timing out")
		c <- 0
	}

	logger.Log.Debug("Getting scheduled")
	entries := scheduler.GetScheduledJobs()
	for _, entry := range entries {
		logger.Log.Debug("Dumping cron entries: %+v", entry)
	}

	go func() {
		logger.Log.Debug("Starting test timeout")
		time.AfterFunc(120*time.Second, endTest)
	}()

	res := <-c
	if res == 1 {
		passed = true
	}
	entries = scheduler.GetScheduledJobs()
	for _, entry := range entries {
		logger.Log.Debug("Test done, Dumping cron entries: %+v", entry)
	}

	if passed == false {
		t.Fail()
	}

}

func disabledTestSchedulerTimeout(t *testing.T) {
	defer scheduler.Stop()

	logger.Log.Debug("Starting Test Timeouts")
	var configs []*metrics.ReportScheduleConfig
	mockdb := mockDB{
		err:     nil,
		configs: configs,
	}

	c := make(chan int)
	mockmsh := metricServiceHandler{
		err:    nil,
		report: nil,
		result: c,
	}

	mockadmindb := mockAdminDB{}

	scheduler.Initialize(&mockmsh, &mockdb, &mockadmindb, 1)

	var testPayload = metrics.ReportScheduleConfig{}

	// Execute every 1 second
	testPayload.DayOfWeek = "*"
	testPayload.DayOfMonth = "*"
	testPayload.Month = "*"
	testPayload.Hour = "*"
	testPayload.Minute = "*"
	testPayload.Second = "*"
	testPayload.Name = "Mock Job 1"
	testPayload.Timeout = 5000
	testPayload.TimeRangeDuration = "P1Y2M10DT2H30M"
	testPayload.Active = true

	var passed bool = false

	mockdb.configs = append(mockdb.configs, &testPayload)
	job := scheduler.SLAConfig{
		Config: &testPayload,
	}

	err := scheduler.AddJob(job)
	if err != nil {
		t.FailNow()
	}

	endTest := func() {
		logger.Log.Debug("Test timing out")
		c <- 0
	}

	logger.Log.Debug("Getting scheduled")
	entries := scheduler.GetScheduledJobs()
	for _, entry := range entries {
		logger.Log.Debug("Dumping cron entries: %+v", entry)
	}

	go func() {
		logger.Log.Debug("Starting test timeout")
		time.AfterFunc(120*time.Second, endTest)
	}()

	res := <-c
	if res == 1 {
		passed = true
	}
	entries = scheduler.GetScheduledJobs()
	for _, entry := range entries {
		logger.Log.Debug("Test done, Dumping cron entries: %+v", entry)
	}

	if passed == false {
		t.Fail()
	}

}

func disabledTestSchedulerNils(t *testing.T) {
	defer scheduler.Stop()

	logger.Log.Debug("Starting Test Nils")
	var configs []*metrics.ReportScheduleConfig
	mockdb := mockDB{
		err:     nil,
		configs: configs,
	}

	c := make(chan int)
	mockmsh := metricServiceHandler{
		err:    nil,
		report: nil,
		result: c,
	}

	mockadmindb := mockAdminDB{}

	scheduler.Initialize(&mockmsh, &mockdb, &mockadmindb, 1)

	var testPayload = metrics.ReportScheduleConfig{}

	// Execute every 1 second
	testPayload.DayOfWeek = "*"
	testPayload.DayOfMonth = "*"
	testPayload.Month = "*"
	testPayload.Hour = "*"
	testPayload.Minute = "*"
	testPayload.Second = "*"
	testPayload.Name = "Mock Job 1"
	testPayload.Timeout = 5000
	testPayload.TimeRangeDuration = "P1Y2M10DT2H30M"
	testPayload.Active = true

	var passed bool = false

	mockdb.configs = append(mockdb.configs, &testPayload)
	job := scheduler.SLAConfig{
		Config: &testPayload,
	}

	err := scheduler.AddJob(job)
	if err != nil {
		t.FailNow()
	}

	endTest := func() {
		logger.Log.Debug("Test timing out")
		c <- 0
	}

	logger.Log.Debug("Getting scheduled")
	entries := scheduler.GetScheduledJobs()
	for _, entry := range entries {
		logger.Log.Debug("Dumping cron entries: %+v", entry)
	}

	go func() {
		logger.Log.Debug("Starting test timeout")
		time.AfterFunc(120*time.Second, endTest)
	}()

	res := <-c
	if res == 1 {
		passed = true
	}
	entries = scheduler.GetScheduledJobs()
	for _, entry := range entries {
		logger.Log.Debug("Test done, Dumping cron entries: %+v", entry)
	}

	if passed == false {
		t.Fail()
	}

}

func TestSchedulerMultiples(t *testing.T) {

	defer scheduler.Stop()

	var configs []*metrics.ReportScheduleConfig

	for i := 0; i < 5; i++ {
		config := metrics.ReportScheduleConfig{
			TimeRangeDuration: "P1",
			ThresholdProfile:  "1000",
			Granularity:       "1000",
			Timeout:           1000,
			Name:              fmt.Sprintf("Test Report: %d", i),
			Second:            "0",
			Minute:            "0",
			Hour:              "3",
			DayOfMonth:        "*",
			Month:             "*",
			DayOfWeek:         "*",
		}
		configs = append(configs, &config)
	}

	mockdb := mockDB{
		err:     nil,
		configs: configs,
	}

	mockmsh := metricServiceHandler{
		err:    nil,
		report: nil,
	}

	mockadmindb := mockAdminDB{}

	scheduler.Initialize(&mockmsh, &mockdb, &mockadmindb, 5)
	var passed bool = false
	var reports = make([]scheduler.SLAConfig, 5)

	for i := 0; i < 5; i++ {
		var testPayload = metrics.ReportScheduleConfig{}

		testPayload.DayOfWeek = "*"
		testPayload.DayOfMonth = "*"
		testPayload.Month = "*"
		testPayload.Hour = "0"
		testPayload.Minute = "0"
		testPayload.Second = "0"

		testPayload.Name = fmt.Sprintf("Test_%d", i)
		testPayload.Active = true

		job := scheduler.SLAConfig{
			Config: &testPayload,
			Job:    nil,
		}
		reports[i] = job

		mockdb.configs = append(mockdb.configs, &testPayload)

		err := scheduler.AddJob(job)
		if err != nil {
			t.FailNow()
		}
	}
	scheduler.RebuildCronJobs()

	var entries []*cron.Entry
	entries = scheduler.GetScheduledJobs()
	for _, entry := range entries {
		logger.Log.Debug("Dumping cron entries: %+v", entry)
	}

	if len(entries) != 5 {
		logger.Log.Errorf("Storage failed:%d ", len(entries))
		t.FailNow()
	}

	logger.Log.Debug("***** Delete 1 item from the configs and check deletion and restart was correct")
	mockdb.DeleteReportScheduleConfig("0", "0")

	scheduler.RebuildCronJobs()

	entries = scheduler.GetScheduledJobs()

	if len(entries) != 4 {
		logger.Log.Errorf("Removed a job but still have %d jobs in cron list", len(entries))
		for _, entry := range entries {
			logger.Log.Debug("Dumping cron entries: %+v", entry)
		}

		t.FailNow()
	} else {
		passed = true
	}

	entries = scheduler.GetScheduledJobs()
	logger.Log.Debug("Test done, Dumping cron entries ", len(entries))
	for i, entry := range entries {
		logger.Log.Debug("%d = %+v", i, entry)
	}

	if passed == false {
		//t.Fail()
	}

}

func TestSpec(t *testing.T) {

	job := metrics.ReportScheduleConfig{
		Minute:           "0",
		Hour:             "0",
		DayOfWeek:        "22",
		DayOfMonth:       "44",
		Month:            "*",
		TenantID:         "0",
		ThresholdProfile: "1234",
	}
	err := job.Validate(false)
	if err == nil {
		logger.Log.Error("Failed 2 invalid values", err)
		t.Fail()
	}

	job = metrics.ReportScheduleConfig{
		Minute:           "0",
		Hour:             "0",
		DayOfWeek:        "22",
		DayOfMonth:       "0",
		Month:            "*",
		TenantID:         "0",
		ThresholdProfile: "1234",
	}
	err = job.Validate(false)
	if err == nil {
		logger.Log.Error("Failed invalid value", err)
		t.Fail()
	}
	job = metrics.ReportScheduleConfig{
		Minute:           "*",
		Hour:             "*",
		DayOfWeek:        "*",
		DayOfMonth:       "*",
		Month:            "*",
		TenantID:         "0",
		ThresholdProfile: "1234",
	}
	err = job.Validate(false)
	if err == nil {
		logger.Log.Error("Failed good config", err)
		t.Fail()
	}

	job = metrics.ReportScheduleConfig{
		Minute:           "4",
		Hour:             "*",
		DayOfWeek:        "*",
		DayOfMonth:       "*",
		Month:            "*",
		TenantID:         "0",
		ThresholdProfile: "1234",
	}
	err = job.Validate(false)
	if err != nil {
		logger.Log.Error("Failed good config", err)
		t.Fail()
	}

	job = metrics.ReportScheduleConfig{
		Minute:           "0",
		Hour:             "0",
		DayOfWeek:        "2",
		DayOfMonth:       "1",
		Month:            "*",
		TenantID:         "0",
		ThresholdProfile: "1234",
	}
	err = job.Validate(false)
	if err != nil {
		logger.Log.Error("Failed good config", err)
		t.Fail()
	}

}
