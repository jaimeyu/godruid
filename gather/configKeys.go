package gather

// ConfigKey - an enumeration used to allow config strings to be
// maintained in one place to prevent erors and allow auto-completion.
type ConfigKey string

const (
	CK_server_datastore_ip                   ConfigKey = "server.datastore.ip"
	CK_server_datastore_port                 ConfigKey = "server.datastore.port"
	CK_server_datastore_batchsize            ConfigKey = "server.datastore.batchsize"
	CK_server_rest_ip                        ConfigKey = "server.rest.ip"
	CK_server_rest_port                      ConfigKey = "server.rest.port"
	CK_server_grpc_ip                        ConfigKey = "server.grpc.ip"
	CK_server_grpc_port                      ConfigKey = "server.grpc.port"
	CK_server_monitoring_port                ConfigKey = "server.monitoring.port"
	CK_server_profile_port                   ConfigKey = "server.profile.port"
	CK_server_websocket_ip                   ConfigKey = "server.websocket.ip"
	CK_server_websocket_port                 ConfigKey = "server.websocket.port"
	CK_server_cors_allowedorigins            ConfigKey = "server.cors.allowedorigins"
	CK_server_changenotif_refreshFreqSeconds ConfigKey = "server.changenotif.refreshFreqSeconds"
	CK_args_admindb_name                     ConfigKey = "args.admindb.name"
	CK_args_admindb_impl                     ConfigKey = "args.admindb.impl"
	CK_args_tenantdb_impl                    ConfigKey = "args.tenantdb.impl"
	CK_args_pouchplugindb_impl               ConfigKey = "args.pouchplugindb.impl"
	CK_args_testdatadb_impl                  ConfigKey = "args.testdatadb.impl"
	CK_args_debug                            ConfigKey = "args.debug"
	CK_args_maxConcurrentMetricAPICalls      ConfigKey = "args.maxConcurrentMetricAPICalls"
	CK_args_maxConcurrentProvAPICalls        ConfigKey = "args.maxConcurrentProvAPICalls"
	CK_args_maxConcurrentPouchAPICalls       ConfigKey = "args.maxConcurrentPouchAPICalls"
	CK_connector_maxSecondsWithoutHeartbeat  ConfigKey = "connector.maxSecondsWithoutHeartbeat"
	CK_druid_broker_server                   ConfigKey = "druid.broker.server"
	CK_druid_broker_port                     ConfigKey = "druid.broker.port"
	CK_druid_broker_table                    ConfigKey = "druid.broker.table"
	CK_druid_coordinator_server              ConfigKey = "druid.coordinator.server"
	CK_druid_coordinator_port                ConfigKey = "druid.coordinator.port"
	CK_kafka_broker                          ConfigKey = "kafka.broker"
	CK_args_authorizationAAA                 ConfigKey = "args.AuthorizationAAA"
)

func (key ConfigKey) String() string {
	return string(key)
}
