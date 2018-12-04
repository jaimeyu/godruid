## Current Release 
### 0.272.0 
**Release Date:** Tue Dec  4 18:44:19 UTC 2018     
## Previous Releases 
### 0.271.0 
**Release Date:** Fri Nov 30 17:21:15 UTC 2018     
* Fix - Guards to ensure that badly formed threshold crossing top N requests do not affect gather
### 0.270.0 
**Release Date:** Fri Nov 30 03:19:14 UTC 2018     
### 0.269.0 
**Release Date:** Wed Nov 28 22:37:52 UTC 2018     
### 0.268.0 
**Release Date:** Tue Nov 27 19:53:39 UTC 2018     
* Fix - Fixes issue where topn was not considering 0 valued metrics
### 0.267.0 
**Release Date:** Tue Nov 27 17:21:33 UTC 2018     
### 0.266.0 
**Release Date:** Tue Nov 27 14:39:19 UTC 2018     
* Feature - Adding a flag to all the v2 APIs to control data cleaning
### 0.265.0 
**Release Date:** Tue Nov 27 11:05:39 UTC 2018     
* Fix - Found bug in Connector's run.sh that did not include the host file override for the login function. 
### 0.264.0 
**Release Date:** Tue Nov 27 02:46:02 UTC 2018     
* Metric Baseline Postgres Impl
### 0.263.0 
**Release Date:** Fri Nov 23 21:41:44 UTC 2018     
### 0.262.0 
**Release Date:** Tue Nov 20 16:41:27 UTC 2018     
* Fix - Add a toggle to add/remove data cleaning to aggregate request 
### 0.261.0 
**Release Date:** Tue Nov 20 14:47:32 UTC 2018     
* Feature - Initial checkin of metric query integration tests
### 0.260.0 
**Release Date:** Tue Nov 20 02:58:26 UTC 2018     
### 0.259.0 
**Release Date:** Mon Nov 19 14:28:33 UTC 2018     
* Feature - Adding V2 API for  bulk meta
### 0.258.0 
**Release Date:** Fri Nov 16 16:08:32 UTC 2018     
* Fix - Alpine image is now pointing to alpine latest instead of edge
### 0.257.0 
**Release Date:** Thu Nov 15 19:50:03 UTC 2018     
* Fix - Removed excessive logging from doGetDownloadRoadrunner and doGetIngestionDictionaryV2
### 0.256.0 
**Release Date:** Thu Nov 15 17:56:17 UTC 2018     
* Fix - Remove unnecessary debug log
### 0.255.0 
**Release Date:** Wed Nov 14 15:14:03 UTC 2018     
* Fix - Updated run.sh to remove the old running aod_connector container before starting a new one.
### 0.254.0 
**Release Date:** Tue Nov 13 22:44:24 UTC 2018     
* Fix - Making the dowloaded Roadrunner tar.gz archive and actual gz archive
### 0.253.0 
**Release Date:** Tue Nov 13 19:35:01 UTC 2018     
* Feature - Multitenancy support for RoadRunner configuration. 
### 0.252.0 
**Release Date:** Tue Nov 13 16:04:51 UTC 2018     
* Fix - Only return threshold topn entries that have a non-zero score against them
### 0.251.0 
**Release Date:** Mon Nov 12 21:35:24 UTC 2018     
Feature - Adding ability to sort topn responses
### 0.250.0 
**Release Date:** Mon Nov 12 20:52:19 UTC 2018     
### 0.249.0 
**Release Date:** Mon Nov 12 18:25:42 UTC 2018     
* Fix - Workaround for Alpine OS broken dependencies. We're getting conflicting SSL library errors so I retargetted the images to use our last good builds until either we fix the images or maintainer fixes the dependency problem.
### 0.248.0 
**Release Date:** Mon Nov 12 15:15:32 UTC 2018     
* Fix - /distribution/sign-csr was whitelisting only Skylight Admin access and was blocking tenant admins. This caused issues with Road Runner trying to connect and setup its certificates because road runner was not using skylight admin credentials anymore. 
### 0.247.0 
**Release Date:** Fri Nov  9 01:16:41 UTC 2018     
### 0.246.0 
**Release Date:** Thu Nov  8 20:33:30 UTC 2018     
### 0.245.0 
**Release Date:** Thu Nov  8 16:48:17 UTC 2018     
### 0.244.0 
**Release Date:** Thu Nov  8 04:33:34 UTC 2018     
* Fix - Roadrunner package updated so it can be started multiple times on a guest for multitenancy deployments. 
### 0.243.0 
**Release Date:** Wed Nov  7 18:57:27 UTC 2018     
### 0.242.0 
**Release Date:** Tue Nov  6 20:08:54 UTC 2018     
Feature - Adding code to process baseline threshold profile queries
Feature - Upgrading docker-go-sdk reference to 0.11.0 which includes Go v1.11.1.
Feature - Ran dep ensure to pull in missing dependencies
Fix     - Build issue
### 0.241.0 
**Release Date:** Fri Nov  2 16:51:51 UTC 2018     
### 0.240.0 
**Release Date:** Thu Nov  1 18:40:27 UTC 2018     
### 0.239.0 
**Release Date:** Wed Oct 31 19:06:11 UTC 2018     
### 0.238.0 
**Release Date:** Wed Oct 24 18:30:39 UTC 2018     
* Fix - Adding metric API guards
### 0.237.0 
**Release Date:** Wed Oct 24 18:14:13 UTC 2018     
### 0.236.0 
**Release Date:** Wed Oct 24 17:27:41 UTC 2018     
* Fix - Adding metric API guards
### 0.235.0 
**Release Date:** Wed Oct 24 16:35:29 UTC 2018     
### 0.234.0 
**Release Date:** Tue Oct 23 18:08:51 UTC 2018     
* Feature - V2 Metric APIs
### 0.233.0 
**Release Date:** Fri Oct 19 15:23:19 UTC 2018     
* Feature - adding tenant meatadata configuration APIs
### 0.232.0 
**Release Date:** Wed Oct 17 16:08:00 UTC 2018     
### 0.231.0 
**Release Date:** Mon Oct 15 14:16:58 UTC 2018     
* Fix - Adding even more debug logs to event handling
### 0.230.0 
**Release Date:** Sat Oct 13 15:42:24 UTC 2018     
* Fix - Restructuring V2 card API
### 0.229.0 
**Release Date:** Sat Oct 13 11:49:31 UTC 2018     
* Fix - Restructuring V2 card API
### 0.228.0 
**Release Date:** Fri Oct 12 18:04:48 UTC 2018     
* Feature - Restructuring V2 card API
### 0.227.0 
**Release Date:** Thu Oct 11 15:25:30 UTC 2018     
### 0.226.0 
**Release Date:** Wed Oct 10 15:52:44 UTC 2018     
* Feature - Adding additional logging to event handler 
### 0.225.0 
**Release Date:** Fri Oct  5 15:53:46 UTC 2018     
* Fix - Removing restriction on data cleaning profile vendor
* Feature - Adding logging to event handler startup
### 0.224.0 
**Release Date:** Thu Sep 27 17:41:43 UTC 2018     
### 0.223.0 
**Release Date:** Wed Sep 26 15:30:57 UTC 2018     
* Fix - handling multiple MOs with the same objectname for bulk meta update
### 0.222.0 
**Release Date:** Thu Sep 20 17:09:17 UTC 2018     
### 0.221.0 
**Release Date:** Thu Sep 20 15:57:58 UTC 2018     
### 0.220.0 
**Release Date:** Wed Sep 19 18:47:44 UTC 2018     
### 0.219.0 
**Release Date:** Thu Sep 13 02:28:03 UTC 2018     
### 0.218.0 
**Release Date:** Thu Sep 13 01:42:47 UTC 2018     
* Fix - handling relationships for the report schedule config APIs
### 0.217.0 
**Release Date:** Thu Sep 13 00:21:23 UTC 2018     
* Fix - handling relationship data on patch for Dashboards
### 0.216.0 
**Release Date:** Wed Sep 12 15:51:47 UTC 2018     
* Feature - adding regulator metric type to the ingestion dictionary
### 0.215.0 
**Release Date:** Tue Sep 11 16:41:39 UTC 2018     
* Feature - adding allowed origin support based on suffix matching
* Fix - always setting the TenantId in a 'patch' scenario to ensure legacy models will work with V2 APIs
### 0.214.0 
**Release Date:** Mon Sep 10 18:37:59 UTC 2018     
* Fix - Metadata value filtering with case insensitivity
### 0.213.0 
**Release Date:** Mon Sep 10 13:39:24 UTC 2018     
* Fix - removing requirement for the 'thresholds' property on Threshold Profile update
### 0.212.0 
**Release Date:** Fri Sep  7 19:19:45 UTC 2018     
* Fix - sort monitored object ids before issuing a query when using metadata
* Feature - APIs for branding and locale
### 0.211.0 
**Release Date:** Wed Sep  5 15:19:28 UTC 2018     
### 0.210.0 
**Release Date:** Tue Sep  4 14:28:35 UTC 2018     
* Feature - data cleaning history api
### 0.209.0 
**Release Date:** Fri Aug 31 02:33:32 UTC 2018     
* Fix - Adding numbers to the allowed metadata key validation.
### 0.208.0 
**Release Date:** Tue Aug 28 00:34:21 UTC 2018     
* Fix - Re-order timeseries data for raw query in chronological order
### 0.207.0 
**Release Date:** Sat Aug 25 11:23:50 UTC 2018     
* Fix - Resolve inconsistency issue for raw query when more than a single monitored object is requested
### 0.206.0 
**Release Date:** Fri Aug 24 19:55:36 UTC 2018     
### 0.205.0 
**Release Date:** Thu Aug 23 20:12:31 UTC 2018     
* Fix - Modify raw query to retrieve all monitored objects if no metadata is specified
### 0.204.0 
**Release Date:** Thu Aug 23 15:32:47 UTC 2018     
* Refactor - adding flattened structure to metric related models in preparation of deprecation of the hierarchical models
### 0.203.0 
**Release Date:** Tue Aug 21 18:21:04 UTC 2018     
### 0.202.0 
**Release Date:** Tue Aug 21 17:29:36 UTC 2018     
* Fix - Issue where couchdb views were not getting generated anymore.
### 0.201.0 
**Release Date:** Mon Aug 20 16:24:35 UTC 2018     
* Fix - Adding metadata validators to CRUID APIs
* Fix - Updated views to be handle lowercase and uppercase values the same
* Fix - Changed metadata validators to only use lowercase letters a...z and the _
### 0.200.0 
**Release Date:** Fri Aug 17 20:34:29 UTC 2018     
* Fix - making sure Druid queries are built using sorted keys to ensure cache hits on subsequent queries
### 0.199.0 
**Release Date:** Fri Aug 17 19:23:12 UTC 2018     
* Feature - Adding endpoint to query monitored objects by metadata
### 0.198.0 
**Release Date:** Wed Aug 15 16:16:21 UTC 2018     
* Feature - Adding configurable druid timeouts

### 0.196.0 
**Release Date:** Tue Aug 14 17:20:09 UTC 2018     
* Feature - dashboard and Card API support
### 0.195.0 
**Release Date:** Tue Aug 14 15:48:43 UTC 2018     
* Fix - Added auth headers to bulk request
* Fix - Removed pollChanges() function from firing periodically. This eats up performance when deployment is scaled up. 
* Fix - Monitored Object swagger definition was missing metadata
### 0.194.0 
**Release Date:** Fri Aug 10 20:25:52 UTC 2018     
### 0.193.0 
**Release Date:** Fri Aug 10 16:50:26 UTC 2018     
* Fix - Added debugging code and monitored object ID mapping fix
* Fix - Bug in TriggerBuildCouchIndex that wasn't properly handling error
* Fix - Added better validator for monitored object metadata
### 0.192.0 
**Release Date:** Fri Aug 10 05:57:39 UTC 2018     
* Fix - Raw metrics API with filtering
* Fix - Couchdb View had a typo
* Fix - 0 data cleaning profiles returns error instead of empty set
### 0.191.0 
**Release Date:** Fri Aug 10 01:01:41 UTC 2018     
* Fix - Raw metrics API with filtering
### 0.190.0 
**Release Date:** Thu Aug  9 20:02:16 UTC 2018     
### 0.189.0 
**Release Date:** Thu Aug  9 19:03:00 UTC 2018     
* Feature - SLA Report V2 APIs
### 0.188.0 
**Release Date:** Thu Aug  9 17:38:04 UTC 2018     
### 0.187.0 
**Release Date:** Thu Aug  9 14:41:00 UTC 2018     
* Feature - Report Schedule Configuration V2 APIs
### 0.186.0 
**Release Date:** Wed Aug  8 20:02:07 UTC 2018     
* Feature - Monitored Object V2 APIs
### 0.185.0 
**Release Date:** Tue Aug  7 15:10:44 UTC 2018     
* Feature - V2 API support for Threshold Profiles
### 0.184.0 
**Release Date:** Wed Aug  1 16:36:22 UTC 2018     
* Feature - V2 APIs for Admin Service 
### 0.183.0 
**Release Date:** Wed Aug  1 00:32:06 UTC 2018     
### 0.182.0 
**Release Date:** Tue Jul 31 13:45:47 UTC 2018     
### 0.181.0 
**Release Date:** Mon Jul 30 19:17:50 UTC 2018     
### 0.180.0 
**Release Date:** Fri Jul 27 19:28:38 UTC 2018     
### 0.179.0 
**Release Date:** Fri Jul 27 14:27:31 UTC 2018     
### 0.178.0 
**Release Date:** Thu Jul 26 18:59:59 UTC 2018     
### 0.177.0 
**Release Date:** Thu Jul 26 15:18:15 UTC 2018     
* Feature - pagination for v2 GetAllMonitoredObjects API
### 0.176.0 
**Release Date:** Wed Jul 25 14:09:56 UTC 2018     
* Fix raw metrics api to include unclean by default
### 0.175.0 
**Release Date:** Mon Jul 23 20:29:41 UTC 2018     
### 0.174.0 
**Release Date:** Fri Jul 20 17:57:47 UTC 2018     
### 0.173.0 
**Release Date:** Fri Jul 20 16:43:39 UTC 2018     
### 0.172.0 
**Release Date:** Fri Jul 20 14:29:20 UTC 2018     
### 0.171.0 
**Release Date:** Fri Jul 20 13:59:02 UTC 2018     
### 0.170.0 
**Release Date:** Thu Jul 19 19:50:34 UTC 2018     
### 0.169.0 
**Release Date:** Thu Jul 19 14:00:31 UTC 2018     
### 0.168.0 
**Release Date:** Tue Jul 17 16:40:38 UTC 2018     
### 0.167.0 
**Release Date:** Tue Jul 17 15:21:47 UTC 2018     
### 0.166.0 
**Release Date:** Tue Jul 17 15:04:37 UTC 2018     
### 0.165.0 
**Release Date:** Mon Jul 16 15:43:02 UTC 2018     
### 0.164.0 
**Release Date:** Sat Jul 14 14:23:34 UTC 2018     
* Feat - include a query filter to query only clean metrics
### 0.163.0 
**Release Date:** Thu Jul 12 14:52:46 UTC 2018     
### 0.162.0 
**Release Date:** Thu Jul 12 00:41:01 UTC 2018     
* Refactor - now using swagger file to generate the server
### 0.161.0 
**Release Date:** Wed Jul  4 15:49:24 UTC 2018     
* Fix - Cannot create domains for other tenant admin users that are not skylight admin. We did provisionned the roles properly, but the for loop that checked for the roles wasn't really looping. [Fixes Issue 47](https://app.trackduck.com/project/5af4415b55b8b593751d5527/issue/5b3bd4b994cf14512d43feb7?utm_source=integration&utm_medium=app_slack&utm_content=task)
### 0.160.0 
**Release Date:** Wed Jun 27 18:23:46 UTC 2018     
### 0.159.0 
**Release Date:** Wed Jun 27 18:02:30 UTC 2018     
* Feature - adding metric instrumentation for Druid queries to track individual time of query from the entire API time.
### 0.158.0 
**Release Date:** Mon Jun 25 18:58:20 UTC 2018     
* Fix - Changed default log level on gather to turn off debug logs
### 0.157.0 
**Release Date:** Mon Jun 25 12:52:02 UTC 2018     
### 0.156.0 
**Release Date:** Mon Jun 25 12:45:42 UTC 2018     
Feature - Developer feature. Develerper's local gather can now access remote druid with real datasets and avoid setting up GCE instances.
Fix - GRP calls to tenant thresholdprofiles did not log tenant ID which made it difficult to debug, https://app.asana.com/0/682255795640582/716246867279196
### 0.155.0 
**Release Date:** Wed Jun 20 21:21:57 UTC 2018     
* Fix - batching Monitored Object name updates to reduce API strain.
### 0.154.0 
**Release Date:** Wed Jun 20 19:15:16 UTC 2018     
* Fix - don't put empty keys in druid lookup
### 0.153.0 
**Release Date:** Wed Jun 20 14:53:20 UTC 2018     
* Fix - Changing notification handler to push notification on kafka using Async in fire and forget mode. This fixes issue [Monitored Object Domain doesn't seem to work on the dashboards](https://app.asana.com/0/710429663603017/716274045857257/f)
### 0.152.0 
**Release Date:** Thu Jun 14 20:15:22 UTC 2018     
### 0.151.0 
**Release Date:** Thu Jun 14 18:25:25 UTC 2018     
### 0.149.0 
**Release Date:** Thu Jun 14 18:01:18 UTC 2018     
* Feature - Adding endpoint to a TopN operation for a metric.
### 0.148.0 
**Release Date:** Wed Jun 13 21:21:18 UTC 2018     
### 0.147.0 
**Release Date:** Wed Jun 13 20:41:40 UTC 2018     
* feat - adding support for custom bucket histogram queries
### 0.146.0 
**Release Date:** Wed Jun 13 20:35:58 UTC 2018     
* Feature - a new query endpoint for threshold crossing that includes violation time
### 0.145.0 
**Release Date:** Wed Jun 13 14:57:43 UTC 2018     
* Fix - support 'All' granularity for metric queries
### 0.144.0 
**Release Date:** Tue Jun 12 02:13:32 UTC 2018     
### 0.143.0 
**Release Date:** Tue Jun 12 01:28:22 UTC 2018     
* Fix - changing the ingestion dictionary based on twamp needs
### 0.142.0 
**Release Date:** Mon Jun 11 15:32:29 UTC 2018     
* Fix - temporarily removing UI section of threshold profiles that is not formatted according to the backend model on UI created threshold profiles. Test commented out as well.
### 0.141.0 
**Release Date:** Fri Jun  8 20:58:19 UTC 2018     
### 0.140.0 
**Release Date:** Fri Jun  8 16:58:11 UTC 2018     
* feat - adding a system level auth role for API calls made by internal services to datahub
* fix - correcting issue with auth being called by internal services for APIs that  o not require auth.
### 0.139.0 
**Release Date:** Fri Jun  8 15:16:39 UTC 2018     
Feature - Adding Tests for Skylight-AAA authorizations utilities.
### 0.138.0 
**Release Date:** Thu Jun  7 17:01:26 UTC 2018     
* Fix - SLA Report generation was not consistent between the immediate generate and the scheduled generation.
### 0.137.0 
**Release Date:** Wed Jun  6 18:49:43 UTC 2018     
Fix - Aligning report scheduler config attribute names with UI model
### 0.136.0 
**Release Date:** Tue Jun  5 19:04:55 UTC 2018     
* fix - correcting issue with auth being called by internal services for APIs that  do not require auth.
### 0.135.0 
**Release Date:** Tue Jun  5 02:48:22 UTC 2018     
* Fix - Exporting GetName to return the proper report type
### 0.134.0 
**Release Date:** Mon Jun  4 20:05:36 UTC 2018     
* Fix - Changing the name of SLAReport and SLASummary to Report and ReportSummary
### 0.133.0 
**Release Date:** Fri Jun  1 18:02:13 UTC 2018     
Feature - Add Authorization Functionality to Gather Endpoints using the new headers from Skylight-AAA
### 0.132.0 
**Release Date:** Fri Jun  1 13:34:36 UTC 2018     
* Feature - CRUD functionality for report scheduling config and SLA report
* Feature - Adding scheduler to schedule SLA Reports
### 0.131.0 
**Release Date:** Thu May 31 14:46:02 UTC 2018     
* Fix - for Druid lookups, use ID attribute if MonitoredObjectID attribute is missing
### 0.130.0 
**Release Date:** Thu May 31 11:56:22 UTC 2018     
### 0.129.0 
**Release Date:** Wed May 30 20:43:20 UTC 2018     
* Feat - new api to bulk update monitored objects
### 0.128.0 
**Release Date:** Tue May 29 15:40:01 UTC 2018     
* Feature - block deletion of domain and threshold profile if used by other resources
### 0.127.0 
**Release Date:** Tue May 29 14:41:51 UTC 2018     
* Fix - Removed 'createdTimestamp' requirement for model validation and updated validation messages to remove the requirement in the response

### 0.126.0 
**Release Date:** Mon May 28 23:58:41 UTC 2018     
* Feature - CRUD functionality for report scheduling config and SLA report
### 0.125.0 
**Release Date:** Mon May 28 23:17:48 UTC 2018     
* Feature - use timezone for SLA bucketing
### 0.124.0 
**Release Date:** Mon May 28 19:20:51 UTC 2018     
### 0.123.0 
**Release Date:** Mon May 28 17:40:29 UTC 2018     
Fix - changed the name of the monitored object datastore
### 0.122.0 
**Release Date:** Mon May 28 15:49:27 UTC 2018     
### 0.121.0 
* Feature - move reports and monitored objects to separate datastore
**Release Date:** Fri May 25 19:24:42 UTC 2018     
* Fix - cleanup swagger file, update dependencies
### 0.120.0 
**Release Date:** Fri May 25 15:03:51 UTC 2018     
### 0.119.0 
**Release Date:** Thu May 24 23:55:47 UTC 2018     
### 0.118.0 
**Release Date:** Thu May 24 21:26:46 UTC 2018     
* Feature - aggregated metrics query
### 0.117.0 
**Release Date:** Thu May 24 16:28:07 UTC 2018     
* Fix - Removed 'createdTimestamp' requirement for model validation in Admin Tenant, Tenant User, Tenant Domain, Tenant Ingestion Profile, Tenant Threshold Profile, Tenant Metadata.
### 0.116.0 
**Release Date:** Wed May 23 18:25:07 UTC 2018     
### 0.115.0 
**Release Date:** Wed May 23 14:32:49 UTC 2018     
* Fix - remove ThresholdProfileSet from the domain model
### 0.114.0 
**Release Date:** Wed May 23 13:14:24 UTC 2018     
### 0.113.0 
**Release Date:** Wed May 16 16:41:53 UTC 2018     
* Fix - handle lookup not found errors from druid
### 0.112.0 
**Release Date:** Fri May 11 18:47:21 UTC 2018     
### 0.111.0 
**Release Date:** Fri May 11 17:53:36 UTC 2018     
* Fix - Adding validity checks to Threshold query by TopN to avoid invalid Druid queries
### 0.110.0 
**Release Date:** Wed May  9 16:25:23 UTC 2018     
### 0.109.0 
**Release Date:** Sun May  6 21:05:18 UTC 2018     
* Feature - Use Druid lookups for domain filtering.
### 0.108.0 
**Release Date:** Wed Apr 25 19:40:48 UTC 2018     
* Fix - SLA report when no metric rows exist for time range
### 0.107.0 
**Release Date:** Wed Apr 25 19:21:16 UTC 2018     
* Fix - adding 'changeNotifications' flag to be able to bypass change notification sub-routine if not needed.
### 0.106.0 
**Release Date:** Wed Apr 25 14:44:33 UTC 2018     
### 0.105.0 
**Release Date:** Tue Apr 24 19:56:43 UTC 2018     
* Fix - add missing kafka config for local builds to start.
### 0.104.0 
**Release Date:** Fri Apr 20 15:09:28 UTC 2018     
* Fix - improvements to SLA report
### 0.103.0 
**Release Date:** Thu Apr 19 20:11:20 UTC 2018     
* Fix - add domain validation to metric API
### 0.102.0 
**Release Date:** Thu Apr 19 19:16:40 UTC 2018     
* Feature - SLA Report
### 0.101.0 
**Release Date:** Mon Apr 16 11:14:53 UTC 2018     
* Feature - send change notification for insert/update monitored object 
### 0.100.0 
**Release Date:** Fri Apr 13 19:13:58 UTC 2018     
* Support domains in metrics queries
### 0.99.0 
**Release Date:** Thu Apr 12 18:49:20 UTC 2018     
* Adding missing dependencies and adding in a very necessary return statement.
### 0.98.0 
**Release Date:** Thu Apr 12 16:27:55 UTC 2018     
* Fix - addressing issue that caused all metrics service APIs to return 404
### 0.97.0 
**Release Date:** Thu Apr 12 13:34:28 UTC 2018     
* Feature - Added poller for notifying changes in monitored objects.
* Feature - increase default CouchDB query results to 1000
* Fix - typos in logs

### 0.96.0 
**Release Date:** Tue Apr 10 18:12:19 UTC 2018     
### 0.95.0 
**Release Date:** Thu Apr  5 14:06:56 UTC 2018     
### 0.94.0 
**Release Date:** Fri Mar 23 15:12:24 UTC 2018     
### 0.93.0 
**Release Date:** Wed Mar 21 20:15:31 UTC 2018     
### 0.92.0 
**Release Date:** Wed Mar 21 20:03:40 UTC 2018     
### 0.91.0 
**Release Date:** Fri Mar 16 16:06:38 UTC 2018     
### 0.90.0 
**Release Date:** Thu Mar  1 17:47:33 UTC 2018     
* Refactor - changing how the tests run so that code coverage will be used

### 0.89.0 
**Release Date:** Mon Feb 26 22:08:31 UTC 2018     
### 0.88.0 
**Release Date:** Mon Feb 26 21:27:49 UTC 2018     
### 0.87.0 
**Release Date:** Mon Feb 26 21:21:23 UTC 2018     
### 0.86.0 
**Release Date:** Mon Feb 26 14:56:37 UTC 2018     
### 0.85.0 
**Release Date:** Sun Feb 25 22:08:41 UTC 2018  
* Added - DBs used by couch to manage the service now created on startup if they are missing.   
### 0.84.0 
**Release Date:** Fri Feb 23 16:40:13 UTC 2018     
### 0.83.0 
**Release Date:** Thu Feb 22 19:14:52 UTC 2018     
### 0.82.0 
**Release Date:** Thu Feb 22 16:49:33 UTC 2018     
### 0.81.0 
**Release Date:** Thu Feb 22 16:43:12 UTC 2018     
### 0.80.0 
**Release Date:** Wed Feb 21 19:59:53 UTC 2018     
### 0.79.0 
**Release Date:** Wed Feb 21 18:01:57 UTC 2018   
* update default threshold profile  
### 0.78.0 
**Release Date:** Tue Feb 20 20:31:41 UTC 2018     
### 0.77.0 
**Release Date:** Sat Feb 17 01:52:52 UTC 2018   
* Change the base color of domain created using the test-data api. 
### 0.76.0 
**Release Date:** Fri Feb 16 17:56:43 UTC 2018     
### 0.75.0 
**Release Date:** Thu Feb 15 20:08:24 UTC 2018     
### 0.74.0 
**Release Date:** Thu Feb 15 16:21:42 UTC 2018   
* Added - fix for intermittent test failure on build machine.  
### 0.73.0 
**Release Date:** Thu Feb 15 15:13:11 UTC 2018   
* Added - new default threshold for MWC  
### 0.72.0 
**Release Date:** Wed Feb 14 19:29:30 UTC 2018     
* Added - testing for tenant service started and addressed load shedding for pouch issue
### 0.71.0 
**Release Date:** Wed Feb 14 16:16:44 UTC 2018     
### 0.70.0 
**Release Date:** Tue Feb 13 17:51:40 UTC 2018     
* Added more tests for the Admin serivice
### 0.69.0 
**Release Date:** Tue Feb 13 17:21:17 UTC 2018     
### 0.68.0 
**Release Date:** Mon Feb 12 20:17:59 UTC 2018     
### 0.67.0 
**Release Date:** Mon Feb 12 17:56:27 UTC 2018     
### 0.66.0 
**Release Date:** Fri Feb  9 19:14:20 UTC 2018   
* Added - API to serve swagger file  
### 0.65.0 
**Release Date:** Fri Feb  9 14:29:44 UTC 2018     
### 0.64.0 
**Release Date:** Tue Feb  6 19:53:05 UTC 2018     
* Added - load shedding configuration for each API grouping.
### 0.63.0 
**Release Date:** Tue Feb  6 18:38:37 UTC 2018  
* Added - metrics to track number of recieved and completed API calls to gather for each service and globally.   
### 0.62.0 
**Release Date:** Tue Feb  6 17:08:38 UTC 2018     
### 0.61.0 
**Release Date:** Tue Feb  6 16:11:51 UTC 2018  
* Added - flowmeter metrics to the default threshold profile
* Added - logging to the queries for the metrics service
* Added - metrics tracking for APIs in the metrics service   
### 0.60.0 
* Updating build process to utilize integration tests
**Release Date:** Thu Feb  1 13:27:22 UTC 2018  
* Added a test framework to gather for unit testing   
### 0.59.0 
**Release Date:** Mon Jan 29 18:39:39 UTC 2018     
### 0.58.0 
**Release Date:** Mon Jan 29 18:19:02 UTC 2018     
### 0.57.0 
**Release Date:** Mon Jan 29 16:57:53 UTC 2018     
### 0.56.0 
**Release Date:** Mon Jan 29 15:49:25 UTC 2018     
### 0.55.0 
**Release Date:** Fri Jan 26 21:42:44 UTC 2018     
### 0.54.0 
**Release Date:** Fri Jan 26 21:02:38 UTC 2018     
### 0.53.0 
**Release Date:** Fri Jan 26 19:44:06 UTC 2018     
### 0.52.0 
**Release Date:** Fri Jan 26 18:43:03 UTC 2018     
### 0.51.0 
**Release Date:** Fri Jan 26 15:05:20 UTC 2018   
* moving monitoring metrics to its own mutex  
### 0.50.0 
**Release Date:** Fri Jan 19 21:18:31 UTC 2018 
* feat - adding metric tracking for API call duration    
### 0.49.0 
**Release Date:** Fri Jan 19 21:07:35 UTC 2018     
### 0.48.0 
**Release Date:** Fri Jan 19 20:49:45 UTC 2018     
### 0.47.0 
**Release Date:** Fri Jan 19 19:32:50 UTC 2018     
### 0.46.0 
**Release Date:** Fri Jan 19 19:10:45 UTC 2018     
### 0.45.0 
**Release Date:** Fri Jan 19 18:54:36 UTC 2018     
### 0.44.0 
**Release Date:** Fri Jan 19 13:48:31 UTC 2018     
### 0.43.0 
**Release Date:** Tue Jan 16 18:53:43 UTC 2018     
### 0.42.0 
**Release Date:** Tue Jan 16 14:50:13 UTC 2018     
### 0.41.0 
**Release Date:** Tue Jan  9 18:28:44 UTC 2018    
* add dependency management and refactor db operations. 
### 0.40.0 
**Release Date:** Tue Jan  9 17:10:33 UTC 2018     
### 0.39.0 
**Release Date:** Fri Jan  5 20:06:32 UTC 2018     
### 0.38.0 
**Release Date:** Fri Jan  5 03:24:03 UTC 2018   
* change default thresh to more reasonable sko values.  
### 0.37.0 
**Release Date:** Wed Jan  3 13:55:21 UTC 2018 
* removing prefix from relational id in tenant metadata    
### 0.36.0 
**Release Date:** Fri Dec 22 20:59:43 UTC 2017     
* remove alias from tenant model and update index to use lowercase name instead.
### 0.35.0 
**Release Date:** Fri Dec 22 20:26:18 UTC 2017 
* change to meta type.    
### 0.34.0 
**Release Date:** Fri Dec 22 20:02:26 UTC 2017    
* add getTenantByAlias API 
### 0.33.0 
**Release Date:** Fri Dec 22 17:16:48 UTC 2017     
### 0.32.0 
**Release Date:** Fri Dec 22 14:57:21 UTC 2017  
* refactor - changing Meta type to TenantMetadata    
### 0.31.0 
**Release Date:** Thu Dec 21 21:53:35 UTC 2017   
 
**Release Date:** Thu Dec 21 21:53:35 UTC 2017     
### 0.30.0 
**Release Date:** Thu Dec 21 20:52:00 UTC 2017
* bug - make sure MO always have at least 1 domain for SKO    
### 0.29.0 
**Release Date:** Thu Dec 21 20:45:05 UTC 2017     
### 0.28.0 
**Release Date:** Thu Dec 21 16:02:04 UTC 2017     
* bug - change '%' to pct in default threshold profile
### 0.27.0 
**Release Date:** Wed Dec 20 20:17:17 UTC 2017     
#### adh-gather:0.27.0
* feat - adding the Tenant Meta data model to tenant DB.
### 0.26.0 
**Release Date:** Wed Dec 20 16:02:07 UTC 2017     
### 0.26.0
* Fixes for SLA Domain Report 
### 0.25.0 
**Release Date:** Wed Dec 20 15:44:38 UTC 2017     
* Adding Domain SLA Report test data generation.

### 0.24.0 
**Release Date:** Wed Dec 20 14:03:14 UTC 2017     
### 0.23.0 
**Release Date:** Tue Dec 19 20:56:00 UTC 2017     
### 0.23.0 
* Refactor naming in ingestion profile and getting rid of relational pouch id prefixes in domains and monitored objects.
### 0.22.0 
**Release Date:** Mon Dec 18 20:23:21 UTC 2017     
### 0.21.0 
**Release Date:** Mon Dec 18 16:03:03 UTC 2017     
### 0.20.0 
**Release Date:** Fri Dec 15 21:58:47 UTC 2017     
### 0.19.0 
**Release Date:** Thu Dec 14 21:45:35 UTC 2017     
### Added
* Adding object count by domain API
### 0.18.0 
**Release Date:** Thu Dec 14 15:51:11 UTC 2017     
### Added
* Changing how IDs are constructed for data stored in Couch.
### 0.17.0 
**Release Date:** Thu Dec 14 14:26:42 UTC 2017     
### 0.16.0 
**Release Date:** Wed Dec 13 23:36:20 UTC 2017     
### Added
* Test data service as well as the linkage between Domain/ThresholdProfile/MonitoredObject in the data model.
### 0.15.0 
**Release Date:** Tue Dec 12 18:07:01 UTC 2017     
### 0.14.0 
**Release Date:** Fri Dec  8 20:21:56 UTC 2017     
### 0.13.0 
**Release Date:** Thu Dec  7 16:31:20 UTC 2017     
### 0.12.0 
**Release Date:** Tue Dec  5 16:02:18 UTC 2017     
* Adding TLS support.

### 0.11.0 
**Release Date:** Mon Nov 27 22:11:07 UTC 2017 

### 0.10.0 
**Release Date:** Tue Nov 21 17:47:47 UTC 2017     

### 0.10.0
**Release Date:** Tue Nov 21 11:37:09 UTC 2017 
* Adding Viper config as well as Docker support.

### 0.9.0 
**Release Date:** Fri Nov 17 22:17:53 UTC 2017     
### 0.9.0
**Release Date:** Tue Nov 14 11:33:09 UTC 2017 
* Changes based on PR feedback for initial ADH-Gather shell. Changes include initial protobuf definition, gRPC service implementation, REST reverse proxy generation, initial CouchDB impl, and separation of service and DAO layers.
### 0.8.0 
**Release Date:** Thu Nov  2 14:14:09 UTC 2017     
### 0.7.0 
**Release Date:** Tue Oct 31 21:04:09 UTC 2017     
* Fix the build

### 0.5.0 
**Release Date:** Tue Oct 31 19:17:16 UTC 2017     
### 0.4.0
**Release Date:** Tue Oct 31 15:13:48 EDT 2017
### 0.3.0
**Release Date:** Tue Oct 31 19:12:18 UTC 2017

