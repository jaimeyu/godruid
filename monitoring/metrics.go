package monitoring

import (
	"time"

	"github.com/accedian/adh-gather/logger"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// GatherMetricPrefix - prefix used for all metrics emmitted from gather
	GatherMetricPrefix = "gather"

	TenantStr                = "tenant"
	AdminUserStr             = "admin_user"
	IngestionDictionaryStr   = "ing_dict"
	IngestionProfileStr      = "ing_prf"
	TenantUserStr            = "tenant_user"
	TenantDomainStr          = "domain"
	ThresholdProfileStr      = "thr_prf"
	MonitoredObjectStr       = "mon_obj"
	ThrCrossStr              = "thr_cross"
	ThrCrossStrTopN          = "thr_cross_topn"
	HistogramStr             = "histogram"
	RawMetricStr             = "raw_metric"
	GenSLAReportStr          = "gen_sla_report"
	SLAReportStr             = "sla_report"
	TenantMetaStr            = "tenant_meta"
	ReportSchedConfigStr     = "report_sched_conf"
	AdminViewsStr            = "admin_views"
	ValidTypesStr            = "valid_types"
	AggMetricsStr            = "aggr_metrics"
	TenantConnectorConfigStr = "connector_config"

	// OPCreateStr - metric constant for a create operation
	OPCreateStr = "create"
	// OPUpdateStr - metric constant for an update operation
	OPUpdateStr = "update"
	// OPPatchStr - metric constant for a patch operation
	OPPatchStr = "patch"
	// OPGetStr - metric constant for a get operation
	OPGetStr = "get"
	// OPDeleteStr - metric constant for a delete operation
	OPDeleteStr = "delete"
	// OPGetAllStr - metric constant for a get operation
	OPGetAllStr = "get_all"
	// OPGetActiveStr - metric constant for a get operation
	OPGetActiveStr = "get_active"
	OPAddStr       = "add"
	OPBulkInsert   = "bulk_insert"
	OPBulkUpdate   = "bulk_update"

	// TimeStr - metric constant for a time metric
	TimeStr    = "time"
	MapStr     = "map"
	IDStr      = "id"
	SummaryStr = "summary"

	// UnitMilliStr - metric constant for a metric measured in milliseconds
	UnitMilliStr = "ms"

	metricNameDelimiter = "_"

	CreateTenantStr = TenantStr + metricNameDelimiter + OPCreateStr
	UpdateTenantStr = TenantStr + metricNameDelimiter + OPUpdateStr
	GetTenantStr    = TenantStr + metricNameDelimiter + OPGetStr
	PatchTenantStr  = TenantStr + metricNameDelimiter + OPPatchStr
	DeleteTenantStr = TenantStr + metricNameDelimiter + OPDeleteStr
	GetAllTenantStr = TenantStr + metricNameDelimiter + OPGetAllStr

	CreateAdminUserStr = AdminUserStr + metricNameDelimiter + OPCreateStr
	UpdateAdminUserStr = AdminUserStr + metricNameDelimiter + OPUpdateStr
	GetAdminUserStr    = AdminUserStr + metricNameDelimiter + OPGetStr
	DeleteAdminUserStr = AdminUserStr + metricNameDelimiter + OPDeleteStr
	GetAllAdminUserStr = AdminUserStr + metricNameDelimiter + OPGetAllStr

	CreateIngDictStr = IngestionDictionaryStr + metricNameDelimiter + OPCreateStr
	UpdateIngDictStr = IngestionDictionaryStr + metricNameDelimiter + OPUpdateStr
	GetIngDictStr    = IngestionDictionaryStr + metricNameDelimiter + OPGetStr
	DeleteIngDictStr = IngestionDictionaryStr + metricNameDelimiter + OPDeleteStr

	CreateIngPrfStr    = IngestionProfileStr + metricNameDelimiter + OPCreateStr
	UpdateIngPrfStr    = IngestionProfileStr + metricNameDelimiter + OPUpdateStr
	GetIngPrfStr       = IngestionProfileStr + metricNameDelimiter + OPGetStr
	PatchIngPrfStr     = IngestionProfileStr + metricNameDelimiter + OPPatchStr
	GetActiveIngPrfStr = IngestionProfileStr + metricNameDelimiter + OPGetActiveStr
	DeleteIngPrfStr    = IngestionProfileStr + metricNameDelimiter + OPDeleteStr

	CreateTenantUserStr = TenantUserStr + metricNameDelimiter + OPCreateStr
	UpdateTenantUserStr = TenantUserStr + metricNameDelimiter + OPUpdateStr
	GetTenantUserStr    = TenantUserStr + metricNameDelimiter + OPGetStr
	PatchTenantUserStr  = TenantUserStr + metricNameDelimiter + OPPatchStr
	DeleteTenantUserStr = TenantUserStr + metricNameDelimiter + OPDeleteStr
	GetAllTenantUserStr = TenantUserStr + metricNameDelimiter + OPGetAllStr

	CreateTenantDomainStr = TenantDomainStr + metricNameDelimiter + OPCreateStr
	UpdateTenantDomainStr = TenantDomainStr + metricNameDelimiter + OPUpdateStr
	GetTenantDomainStr    = TenantDomainStr + metricNameDelimiter + OPGetStr
	PatchTenantDomainStr  = TenantDomainStr + metricNameDelimiter + OPPatchStr
	DeleteTenantDomainStr = TenantDomainStr + metricNameDelimiter + OPDeleteStr
	GetAllTenantDomainStr = TenantDomainStr + metricNameDelimiter + OPGetAllStr

	CreateTenantConnectorConfigStr = TenantConnectorConfigStr + metricNameDelimiter + OPCreateStr
	UpdateTenantConnectorConfigStr = TenantConnectorConfigStr + metricNameDelimiter + OPUpdateStr
	GetTenantConnectorConfigStr    = TenantConnectorConfigStr + metricNameDelimiter + OPGetStr
	DeleteTenantConnectorConfigStr = TenantConnectorConfigStr + metricNameDelimiter + OPDeleteStr
	GetAllTenantConnectorConfigStr = TenantConnectorConfigStr + metricNameDelimiter + OPGetAllStr

	CreateThrPrfStr = ThresholdProfileStr + metricNameDelimiter + OPCreateStr
	UpdateThrPrfStr = ThresholdProfileStr + metricNameDelimiter + OPUpdateStr
	GetThrPrfStr    = ThresholdProfileStr + metricNameDelimiter + OPGetStr
	PatchThrPrfStr  = ThresholdProfileStr + metricNameDelimiter + OPPatchStr
	GetAllThrPrfStr = ThresholdProfileStr + metricNameDelimiter + OPGetAllStr
	DeleteThrPrfStr = ThresholdProfileStr + metricNameDelimiter + OPDeleteStr

	CreateMonObjStr      = MonitoredObjectStr + metricNameDelimiter + OPCreateStr
	UpdateMonObjStr      = MonitoredObjectStr + metricNameDelimiter + OPUpdateStr
	GetMonObjStr         = MonitoredObjectStr + metricNameDelimiter + OPGetStr
	PatchMonObjStr       = MonitoredObjectStr + metricNameDelimiter + OPPatchStr
	GetAllMonObjStr      = MonitoredObjectStr + metricNameDelimiter + OPGetAllStr
	DeleteMonObjStr      = MonitoredObjectStr + metricNameDelimiter + OPDeleteStr
	GetMonObjToDomMapStr = MonitoredObjectStr + metricNameDelimiter + TenantDomainStr + metricNameDelimiter + MapStr + metricNameDelimiter + OPGetStr

	GetThrCrossStr             = ThrCrossStr + metricNameDelimiter + OPGetStr
	GetThrCrossByMonObjStr     = ThrCrossStrTopN + metricNameDelimiter + MonitoredObjectStr + metricNameDelimiter + OPGetStr
	GetThrCrossByMonObjTopNStr = ThrCrossStr + metricNameDelimiter + MonitoredObjectStr + metricNameDelimiter + OPGetStr
	GetHistogramObjStr         = HistogramStr + metricNameDelimiter + OPGetStr
	GetRawMetricStr            = RawMetricStr + metricNameDelimiter + OPGetStr
	GenerateSLAReportStr       = GenSLAReportStr + metricNameDelimiter + OPGetStr

	QueryAggregatedMetricsStr = AggMetricsStr + metricNameDelimiter + OPGetStr

	CreateTenantMetaStr = TenantMetaStr + metricNameDelimiter + OPCreateStr
	UpdateTenantMetaStr = TenantMetaStr + metricNameDelimiter + OPUpdateStr
	GetTenantMetaStr    = TenantMetaStr + metricNameDelimiter + OPGetStr
	PatchTenantMetaStr  = TenantMetaStr + metricNameDelimiter + OPPatchStr
	DeleteTenantMetaStr = TenantMetaStr + metricNameDelimiter + OPDeleteStr

	CreateReportScheduleConfigStr = ReportSchedConfigStr + metricNameDelimiter + OPCreateStr
	UpdateReportScheduleConfigStr = ReportSchedConfigStr + metricNameDelimiter + OPUpdateStr
	GetReportScheduleConfigStr    = ReportSchedConfigStr + metricNameDelimiter + OPGetStr
	GetAllReportScheduleConfigStr = ReportSchedConfigStr + metricNameDelimiter + OPGetAllStr
	DeleteReportScheduleConfigStr = ReportSchedConfigStr + metricNameDelimiter + OPDeleteStr

	GetSLAReportStr    = SLAReportStr + metricNameDelimiter + OPGetStr
	GetAllSLAReportStr = SLAReportStr + metricNameDelimiter + OPGetAllStr

	CreateValidTypesStr      = ValidTypesStr + metricNameDelimiter + OPCreateStr
	UpdateValidTypesStr      = ValidTypesStr + metricNameDelimiter + OPUpdateStr
	GetValidTypesStr         = ValidTypesStr + metricNameDelimiter + OPGetStr
	GetSpecificValidTypesStr = ValidTypesStr + metricNameDelimiter + OPGetStr + "_spec"
	DeleteValidTypesStr      = ValidTypesStr + metricNameDelimiter + OPDeleteStr

	GetTenantIDByAliasStr      = IDStr + metricNameDelimiter + "_by_alais" + metricNameDelimiter + OPGetStr
	GetTenantSummaryByAliasStr = SummaryStr + metricNameDelimiter + "_by_alais" + metricNameDelimiter + OPGetStr
	AddAdminViewsStr           = AdminViewsStr + metricNameDelimiter + OPAddStr

	BulkInsertMonObjStr = MonitoredObjectStr + metricNameDelimiter + OPBulkInsert
	BulkUpdateMonObjStr = MonitoredObjectStr + metricNameDelimiter + OPBulkUpdate
)

type MetricCounterType string

const (
	APIRecieved              MetricCounterType = "APIRecieved"
	APICompleted             MetricCounterType = "APICompleted"
	AdminAPIRecieved         MetricCounterType = "AdminAPIRecieved"
	AdminAPICompleted        MetricCounterType = "AdminAPICompleted"
	TenantAPIRecieved        MetricCounterType = "TenantAPIRecieved"
	TenantAPICompleted       MetricCounterType = "TenantAPICompleted"
	PouchAPIRecieved         MetricCounterType = "PouchAPIRecieved"
	PouchAPICompleted        MetricCounterType = "PouchAPICompleted"
	PouchChangesAPIRecieved  MetricCounterType = "PouchChangesAPIRecieved"
	PouchChangesAPICompleted MetricCounterType = "PouchChangesAPICompleted"
	MetricAPIRecieved        MetricCounterType = "MetricAPIRecieved"
	MetricAPICompleted       MetricCounterType = "MetricAPICompleted"
)

var (
	// APICallDuration - Time it takes to complete a call.
	APICallDuration prometheus.SummaryVec

	// RecievedAPICalls - the number of API calls gather has recieved since startup
	RecievedAPICalls prometheus.Counter

	// CompletedAPICalls - number of API calls gather has completed since startup
	CompletedAPICalls prometheus.Counter

	// CompletedAdminServiceAPICalls - number of API calls the admin service has completed since startup
	CompletedAdminServiceAPICalls prometheus.Counter

	// RecievedAdminServiceAPICalls - the number of API calls the admin service has recieved since startup
	RecievedAdminServiceAPICalls prometheus.Counter

	// CompletedTenantServiceAPICalls - number of API calls the tenant service has completed since startup
	CompletedTenantServiceAPICalls prometheus.Counter

	// RecievedTenantServiceAPICalls - the number of API calls the tenant service has recieved since startup
	RecievedTenantServiceAPICalls prometheus.Counter

	// CompletedPouchServiceAPICalls - number of API calls the pouch service has completed since startup
	CompletedPouchServiceAPICalls prometheus.Counter

	// RecievedPouchChangesAPICalls - the number of API calls pouch service has recieved since startup
	RecievedPouchChangesAPICalls prometheus.Counter

	// CompletedPouchChangesAPICalls - number of API calls the pouch service has completed since startup
	CompletedPouchChangesAPICalls prometheus.Counter

	// RecievedPouchServiceAPICalls - the number of API calls pouch service has recieved since startup
	RecievedPouchServiceAPICalls prometheus.Counter

	// CompletedMetricServiceAPICalls - number of API calls the metric service has completed since startup
	CompletedMetricServiceAPICalls prometheus.Counter

	// RecievedMetricServiceAPICalls - the number of API calls metric service has recieved since startup
	RecievedMetricServiceAPICalls prometheus.Counter
)

// InitMetrics - registers all metrics to be collected for Gather.
func InitMetrics() {
	APICallDuration = *prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name: "gather_api_call_duration",
		Help: "Time taken to execute an API call",
		Objectives: map[float64]float64{
			0.5:  0.1,
			0.9:  0.5,
			0.99: 1.0},
	}, []string{"code", "name"})

	RecievedAPICalls = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "gather_api_call_received_since_startup",
		Help: "Number of API calls recieved by Gather since startup"})

	CompletedAPICalls = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "gather_api_call_completed_since_startup",
		Help: "Number of API calls completed by Gather since startup"})

	CompletedAdminServiceAPICalls = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "gather_admin_service_api_completed_since_startup",
		Help: "Number of API calls completed by the Admin Service since startup"})

	RecievedAdminServiceAPICalls = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "gather_admin_service_api_recieved_since_startup",
		Help: "Number of API calls recieved by the Admin Service since startup"})

	RecievedTenantServiceAPICalls = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "gather_tenant_service_api_call_received_since_startup",
		Help: "Number of API calls recieved by the Tenant Service since startup"})

	CompletedTenantServiceAPICalls = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "gather_tenant_service_api_call_completed_since_startup",
		Help: "Number of API calls completed by the Tenant Service since startup"})

	RecievedPouchServiceAPICalls = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "gather_pouch_service_api_call_received_since_startup",
		Help: "Number of API calls recieved by the Pouch Service since startup"})

	CompletedPouchServiceAPICalls = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "gather_pouch_service_api_call_completed_since_startup",
		Help: "Number of API calls completed by the Pouch Service since startup"})

	RecievedPouchChangesAPICalls = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "gather_pouch_changes_api_call_received_since_startup",
		Help: "Number of API calls recieved for the _changes API since startup"})

	CompletedPouchChangesAPICalls = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "gather_pouch_changes_api_call_completed_since_startup",
		Help: "Number of API calls completed for the _changes API since startup"})

	RecievedMetricServiceAPICalls = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "gather_metric_service_api_call_received_since_startup",
		Help: "Number of API calls recieved by the Metric Service since startup"})

	CompletedMetricServiceAPICalls = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "gather_metric_service_api_call_completed_since_startup",
		Help: "Number of API calls completed by the Metric Service since startup"})

	prometheus.MustRegister(APICallDuration)
	prometheus.MustRegister(RecievedAPICalls)
	prometheus.MustRegister(CompletedAPICalls)
	prometheus.MustRegister(CompletedAdminServiceAPICalls)
	prometheus.MustRegister(RecievedAdminServiceAPICalls)
	prometheus.MustRegister(RecievedTenantServiceAPICalls)
	prometheus.MustRegister(CompletedTenantServiceAPICalls)
	prometheus.MustRegister(RecievedPouchServiceAPICalls)
	prometheus.MustRegister(CompletedPouchServiceAPICalls)
	prometheus.MustRegister(RecievedPouchChangesAPICalls)
	prometheus.MustRegister(CompletedPouchChangesAPICalls)
	prometheus.MustRegister(RecievedMetricServiceAPICalls)
	prometheus.MustRegister(CompletedMetricServiceAPICalls)
}

// TrackAPITimeMetricInSeconds - helper function to track metrics related to API call duration.
func TrackAPITimeMetricInSeconds(startTime time.Time, labels ...string) {
	duration := time.Since(startTime).Seconds()

	logger.Log.Infof("%v: %f", labels, duration)
	APICallDuration.WithLabelValues(labels...).Observe(duration)
}

// IncrementCounter - increments the value of a counter.
func IncrementCounter(counterType MetricCounterType) {
	switch counterType {
	case APIRecieved:
		RecievedAPICalls.Inc()
	case APICompleted:
		CompletedAPICalls.Inc()
	case AdminAPICompleted:
		CompletedAdminServiceAPICalls.Inc()
	case AdminAPIRecieved:
		RecievedAdminServiceAPICalls.Inc()
	case TenantAPICompleted:
		CompletedTenantServiceAPICalls.Inc()
	case TenantAPIRecieved:
		RecievedTenantServiceAPICalls.Inc()
	case PouchAPICompleted:
		CompletedPouchServiceAPICalls.Inc()
	case PouchAPIRecieved:
		RecievedPouchServiceAPICalls.Inc()
	case PouchChangesAPICompleted:
		CompletedPouchChangesAPICalls.Inc()
	case PouchChangesAPIRecieved:
		RecievedPouchChangesAPICalls.Inc()
	case MetricAPICompleted:
		CompletedMetricServiceAPICalls.Inc()
	case MetricAPIRecieved:
		RecievedMetricServiceAPICalls.Inc()
	default:
		logger.Log.Debugf("Unable to increment counter type %v", counterType)
	}

}
