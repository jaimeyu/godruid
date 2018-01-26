package monitoring

import (
	"time"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// GatherMetricPrefix - prefix used for all metrics emmitted from gather
	GatherMetricPrefix = "gather"

	TenantStr = "tenant"
	AdminUserStr = "admin_user"
	IngestionDictionaryStr = "ing_dict"
	IngestionProfileStr = "ing_prf"
	TenantUserStr = "tenant_user"
	TenantDomainStr = "domain"
	ThresholdProfileStr = "thr_prf"
	MonitoredObjectStr = "mon_obj"
	ThrCrossStr = "thr_cross"
	HistogramStr = "histogram"
	TenantMetaStr = "tenant_meta"
	AdminViewsStr = "admin_views"
	ValidTypesStr = "valid_types"

	// OPCreateStr - metric constant for a create operation
	OPCreateStr = "create"
	// OPUpdateStr - metric constant for an update operation
	OPUpdateStr = "update"
	// OPGetStr - metric constant for a get operation
	OPGetStr = "get"
	// OPDeleteStr - metric constant for a delete operation
	OPDeleteStr = "delete"
	// OPGetAllStr - metric constant for a get operation
	OPGetAllStr = "get_all"
	// OPGetActiveStr - metric constant for a get operation
	OPGetActiveStr = "get_active"
	OPAddStr = "add"

	// TimeStr - metric constant for a time metric
	TimeStr = "time"
	MapStr = "map"
	IDStr = "id"

	// UnitMilliStr - metric constant for a metric measured in milliseconds
	UnitMilliStr = "ms"

	metricNameDelimiter = "_"

	CreateTenantStr = TenantStr + metricNameDelimiter + OPCreateStr
	UpdateTenantStr = TenantStr + metricNameDelimiter + OPUpdateStr
	GetTenantStr = TenantStr + metricNameDelimiter + OPGetStr
	DeleteTenantStr = TenantStr + metricNameDelimiter + OPDeleteStr
	GetAllTenantStr = TenantStr + metricNameDelimiter + OPGetAllStr

	CreateAdminUserStr = AdminUserStr + metricNameDelimiter + OPCreateStr
	UpdateAdminUserStr = AdminUserStr + metricNameDelimiter + OPUpdateStr
	GetAdminUserStr = AdminUserStr + metricNameDelimiter + OPGetStr
	DeleteAdminUserStr = AdminUserStr + metricNameDelimiter + OPDeleteStr
	GetAllAdminUserStr = AdminUserStr + metricNameDelimiter + OPGetAllStr

	CreateIngDictStr = IngestionDictionaryStr + metricNameDelimiter + OPCreateStr
	UpdateIngDictStr = IngestionDictionaryStr + metricNameDelimiter + OPUpdateStr
	GetIngDictStr = IngestionDictionaryStr + metricNameDelimiter + OPGetStr
	DeleteIngDictStr = IngestionDictionaryStr + metricNameDelimiter + OPDeleteStr

	CreateIngPrfStr = IngestionProfileStr + metricNameDelimiter + OPCreateStr
	UpdateIngPrfStr = IngestionProfileStr + metricNameDelimiter + OPUpdateStr
	GetIngPrfStr = IngestionProfileStr + metricNameDelimiter + OPGetStr
	GetActiveIngPrfStr = IngestionProfileStr + metricNameDelimiter + OPGetActiveStr
	DeleteIngPrfStr = IngestionProfileStr + metricNameDelimiter + OPDeleteStr

	CreateTenantUserStr = TenantUserStr + metricNameDelimiter + OPCreateStr
	UpdateTenantUserStr = TenantUserStr + metricNameDelimiter + OPUpdateStr
	GetTenantUserStr = TenantUserStr + metricNameDelimiter + OPGetStr
	DeleteTenantUserStr = TenantUserStr + metricNameDelimiter + OPDeleteStr
	GetAllTenantUserStr = TenantUserStr + metricNameDelimiter + OPGetAllStr

	CreateTenantDomainStr = TenantDomainStr + metricNameDelimiter + OPCreateStr
	UpdateTenantDomainStr = TenantDomainStr + metricNameDelimiter + OPUpdateStr
	GetTenantDomainStr = TenantDomainStr + metricNameDelimiter + OPGetStr
	DeleteTenantDomainStr = TenantDomainStr + metricNameDelimiter + OPDeleteStr
	GetAllTenantDomainStr = TenantDomainStr + metricNameDelimiter + OPGetAllStr

	CreateThrPrfStr = ThresholdProfileStr + metricNameDelimiter + OPCreateStr
	UpdateThrPrfStr = ThresholdProfileStr + metricNameDelimiter + OPUpdateStr
	GetThrPrfStr = ThresholdProfileStr + metricNameDelimiter + OPGetStr
	GetAllThrPrfStr = ThresholdProfileStr + metricNameDelimiter + OPGetAllStr
	DeleteThrPrfStr = ThresholdProfileStr + metricNameDelimiter + OPDeleteStr

	CreateMonObjStr = MonitoredObjectStr + metricNameDelimiter + OPCreateStr
	UpdateMonObjStr = MonitoredObjectStr + metricNameDelimiter + OPUpdateStr
	GetMonObjStr = MonitoredObjectStr + metricNameDelimiter + OPGetStr
	GetAllMonObjStr = MonitoredObjectStr + metricNameDelimiter + OPGetAllStr
	DeleteMonObjStr = MonitoredObjectStr + metricNameDelimiter + OPDeleteStr
	GetMonObjToDomMapStr = MonitoredObjectStr + metricNameDelimiter + TenantDomainStr + metricNameDelimiter + MapStr + metricNameDelimiter + OPGetStr

	GetThrCrossStr = ThrCrossStr + metricNameDelimiter + OPGetStr
	GetThrCrossByMonObjStr = ThrCrossStr + metricNameDelimiter + MonitoredObjectStr + metricNameDelimiter + OPGetStr
	GetHistogramObjStr = HistogramStr + metricNameDelimiter + OPGetStr

	CreateTenantMetaStr = TenantMetaStr + metricNameDelimiter + OPCreateStr
	UpdateTenantMetaStr = TenantMetaStr + metricNameDelimiter + OPUpdateStr
	GetTenantMetaStr = TenantMetaStr + metricNameDelimiter + OPGetStr
	DeleteTenantMetaStr = TenantMetaStr + metricNameDelimiter + OPDeleteStr

	CreateValidTypesStr = ValidTypesStr + metricNameDelimiter + OPCreateStr
	UpdateValidTypesStr = ValidTypesStr + metricNameDelimiter + OPUpdateStr
	GetValidTypesStr = ValidTypesStr + metricNameDelimiter + OPGetStr
	GetSpecificValidTypesStr = ValidTypesStr + metricNameDelimiter + OPGetStr + "_spec"
	DeleteValidTypesStr = ValidTypesStr + metricNameDelimiter + OPDeleteStr

	GetTenantIDByAliasStr = IDStr + metricNameDelimiter + "_by_alais" + metricNameDelimiter + OPGetStr
	AddAdminViewsStr = AdminViewsStr + metricNameDelimiter + OPAddStr
)
	
var (
	// APICallDuration - Time it takes to create a Tenant in Gather.
	APICallDuration prometheus.HistogramVec
)

// InitMetrics - registers all metrics to be collected for Gather.
func InitMetrics() {
	APICallDuration = *prometheus.NewHistogramVec(prometheus.HistogramOpts{
        Name: "gather_api_call_duration",
        Help: "Time taken to execute an API call",
	}, []string{"name","code"})

	prometheus.MustRegister(APICallDuration)
}

// GenerateMetricName - used to generate a properly formatted name for a metric.
// func GenerateMetricName(nameParts ...string) string {
// 	partsWithPrefix := []string{}
// 	partsWithPrefix = append(partsWithPrefix, nameParts...)

// 	return strings.Join(partsWithPrefix, "_")
// }

// TrackAPITimeMetricInMilli - helper function to track metrics related to API call duration.
func TrackAPITimeMetricInMilli(startTime time.Time, labels ...string) {
	duration := time.Since(startTime).Seconds() * float64(time.Millisecond)
	APICallDuration.WithLabelValues(labels...).Observe(duration)
}
