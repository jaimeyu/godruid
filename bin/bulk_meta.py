import argparse
import csv
import http.client, urllib.parse
import json
import logging

def conf_logging():
    logging.basicConfig(format='%(asctime)s:%(levelname)s: %(message)s', level=logging.INFO)
    logging.info("Logging configured...")

def construct_entry_func(headers):
    def construct_entry(entry): 
        return json.dumps(dict(zip(headers, entry)))
    return construct_entry

def envoy_func(conn, auth, host, tenantid):
    def envoy(b_id, batch):
        batchlist = list(batch)
        logging.info("Sending batch %s of size %d to %s", b_id, len(batchlist), host)
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

def mo_id(conn, auth, host, tenantid, objectname):
    logging.debug("Retrieving monitored object ID with object name %s", objectname)
    conn.request("GET","/couchdb/tenant_2_"+tenantid+"_monitored-objects/_design/indexOfobjectName/_view/byobjectName?startkey=[\""+objectname+"\"]&endkey=[\""+objectname+"\"]",body=None, headers={"Authorization":auth})
    mo_response = conn.getresponse()
    if mo_response.status != 200:
        logging.error("Could not retrieve monitored object id for object name %s", objectname)
        return
    r_data = json.loads(mo_response.read().decode("utf-8"))
    object_mo_id = r_data["rows"][0]["id"] #TODO work on this
    return object_mo_id


def process(file, batchsize, f_envoy):
    with open(file) as csvfile:
        bulkmetareader = csv.reader(csvfile, delimiter=',')
        i = 0
        headers = bulkmetareader.__next__()
        entry_func = construct_entry_func(headers)
        batch = []

        logging.info("Processing csv with headers: \n" + "\n".join(headers))
        for entry in bulkmetareader:
            i = i+1
            batch += [entry]
            if i == batchsize:
                f_envoy(i/batchsize, map(entry_func, batch))
                batch = []
        
        if len(batch) > 0:
            f_envoy(i/batchsize, map(entry_func, batch))

# Process the command line arguments
parser = argparse.ArgumentParser(description="Bulk insert meta information against monitored objects in datahub.")
parser.add_argument("-b", "--batchsize",type=int, help="Total size of a batch of metadata entries that should be sent to datahub")
parser.add_argument("-f", "--file", help="Absolute path to the csv file containing meta information")
parser.add_argument("-s", "--host", help="Host to send the metadata information to")
parser.add_argument("-u", "--username", help="Username to be used for logging into datahub")
parser.add_argument("-p", "--password", help="Password to be used for logging into datahub")
parser.add_argument("-t", "--tenantname", help="Name of the tenant that the monitored objects to be enriched are associated with")

args = parser.parse_args()

batchsize = 5
metafile = "/Users/abatosparac/go/src/github.com/accedian/adh-gather/bin/test.csv"
host = "jyu.npav.accedian.net"
username = "admin@datahub.com"
password = "AccedianPass"
tenant = "iris"

# batchsize = args.batchsize
# metafile = args.file
# host = args.host
# username = args.username
# password = args.password
# tenant = args.tenant

conf_logging()

logging.info("Loading entries from " + metafile)

conn = http.client.HTTPSConnection(host, timeout=5)

logging.info("Logging into datahub...")
auth = login(conn, host, username, password)
if auth is None:
    logging.error("Could not login. Exiting...")

tid = tenant_id(conn, auth, host, tenant)

print(mo_id(conn, auth, host, "1eda2fea-3571-4600-ae0f-b9ed2a6071e8", "00f8EdIAVC"))

#process(metafile, batchsize, envoy_func(conn, auth, host, tid))
