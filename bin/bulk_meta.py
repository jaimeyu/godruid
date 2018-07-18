import argparse
import csv
import http.client, urllib.parse
import json
import logging

def conf_logging():
    logging.basicConfig(format='%(asctime)s:%(levelname)s: %(message)s', level=logging.INFO)
    logging.info("Logging configured...")

def construct_entry_func(headers, metadatakey):
    def construct_entry(entry): 
        metadata = dict(zip(headers, entry))
        entry_struct = { "keyName":metadata[metadatakey],"MetadataKey":metadatakey,"metadata":metadata}
        return json.dumps(entry_struct)
    return construct_entry

def envoy_func(conn, auth, host, tenantid):
    def envoy(b_id, batch):
        batchlist = list(batch)
        logging.info("Sending batch %s of size %d to %s", b_id, len(batchlist), host)
        payload = "{'_id'='1','_rev'='1','items':[" + ",".join(batchlist)+"]}"
        print(payload)
    return envoy

def login(conn, host, username, password):
    conn.request("POST","/api/v1/auth/login",urllib.parse.urlencode({"username":username,"password":password}),{"Content-Type":"application/x-www-form-urlencoded"})
    login_response = conn.getresponse()
    if login_response.status != 200:
        logging.error("Could not login to host %s", host)
        return
    login_response.read()
    return login_response.getheader("Authorization")

def tenant_id(conn, auth, host, tenantname):
    logging.info("Retrieving tenant ID for %s", tenantname)
    conn.request("GET","/api/v1/tenant-by-alias/" + tenantname, body=None, headers={"Authorization":auth})
    alias_response = conn.getresponse()
    if alias_response.status != 200:
        logging.error("Could not retrieve information for tenant %s at host %s", tenantname, host)
        return
    return alias_response.read().decode("utf-8")

def process(file, batchsize, keyname, f_envoy):
    with open(file) as csvfile:
        bulkmetareader = csv.reader(csvfile, delimiter=',')
        i = 0
        headers = bulkmetareader.__next__()
        f_entry = construct_entry_func(headers, keyname)
        batch = []

        logging.info("Processing csv with headers: \n" + "\n".join(headers))
        for entry in bulkmetareader:
            batch += [entry]
            i += 1
            if i%batchsize == 0:
                f_envoy(i//batchsize, map(f_entry, batch))
                batch = []
        
        if len(batch) > 0:
            f_envoy(i//batchsize, map(f_entry, batch))

# Process the command line arguments
parser = argparse.ArgumentParser(description="Bulk insert meta information against monitored objects in datahub.")
parser.add_argument("-b", "--batchsize",type=int, help="Total size of a batch of metadata entries that should be sent to datahub")
parser.add_argument("-f", "--file", help="Absolute path to the csv file containing meta information")
parser.add_argument("-s", "--host", help="Host to send the metadata information to")
parser.add_argument("-u", "--username", help="Username to be used for logging into datahub")
parser.add_argument("-p", "--password", help="Password to be used for logging into datahub")
parser.add_argument("-t", "--tenantname", help="Name of the tenant that the monitored objects to be enriched are associated with")

args = parser.parse_args()

conf_logging()

logging.info("Loading entries from " + metafile)

conn = http.client.HTTPSConnection(host, timeout=5)

logging.info("Logging into datahub...")
auth = login(conn, host, username, password)
if auth is None:
    logging.error("Could not login. Exiting...")

tid = tenant_id(conn, auth, host, tenant)

process(metafile, batchsize, keyname, envoy_func(conn, auth, host, tid))

logging.info("Finishing processing %s", metafile)