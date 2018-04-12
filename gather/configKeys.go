package gather

// ConfigKey - an enumeration used to allow config strings to be
// maintained in one place to prevent erors and allow auto-completion.
type ConfigKey string

const (
	CK_server_datastore_ip              ConfigKey = "server.datastore.ip"
	CK_server_datastore_port            ConfigKey = "server.datastore.port"
	CK_server_rest_ip                   ConfigKey = "server.rest.ip"
	CK_server_rest_port                 ConfigKey = "server.rest.port"
	CK_server_grpc_ip                   ConfigKey = "server.grpc.ip"
	CK_server_grpc_port                 ConfigKey = "server.grpc.port"
	CK_server_monitoring_port           ConfigKey = "server.monitoring.port"
	CK_server_profile_port              ConfigKey = "server.profile.port"
	CK_server_cors_allowedorigins       ConfigKey = "server.cors.allowedorigins"
	CK_args_admindb_name                ConfigKey = "args.admindb.name"
	CK_args_admindb_impl                ConfigKey = "args.admindb.impl"
	CK_args_tenantdb_impl               ConfigKey = "args.tenantdb.impl"
	CK_args_pouchplugindb_impl          ConfigKey = "args.pouchplugindb.impl"
	CK_args_testdatadb_impl             ConfigKey = "args.testdatadb.impl"
	CK_args_debug                       ConfigKey = "args.debug"
	CK_args_maxConcurrentMetricAPICalls ConfigKey = "args.maxConcurrentMetricAPICalls"
	CK_args_maxConcurrentProvAPICalls   ConfigKey = "args.maxConcurrentProvAPICalls"
	CK_args_maxConcurrentPouchAPICalls  ConfigKey = "args.maxConcurrentPouchAPICalls"
	CK_druid_server                     ConfigKey = "druid.server"
	CK_druid_port                       ConfigKey = "druid.port"
	CK_druid_table                      ConfigKey = "druid.table"
	CK_kafka_broker                     ConfigKey = "kafka.broker"
)

func (key ConfigKey) String() string {
	return string(key)
}
