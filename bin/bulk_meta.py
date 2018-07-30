import argparse
import calendar
import csv
import http.client, urllib.parse
import json
import logging
import os
import sys
import time

def conf_logging():
    logging.basicConfig(format='%(asctime)s:%(levelname)s: %(message)s', level=logging.INFO)
    logging.info("Logging configured...")

def open_processed_file():
    processed_headers = ["key","status"]
    processed_filename = "/tmp/processed.csv"
    logging.info("Generating processed log at %s", processed_filename)
    f = open(processed_filename,"w+")
    f.write(",".join(processed_headers))
    return f

def construct_entry_func(headers, metadatakey):
    def construct_entry(entry): 
        metadata = {k: v for k, v in dict(zip(headers,entry)).items() if v}
        entry_struct = { "keyName":metadatakey,"metadataKey":metadata[metadatakey],"metadata":metadata}
        return json.dumps(entry_struct)
    return construct_entry

def envoy_func(conn, auth, host, tenantid, processfile):
    def envoy(b_id, batch):
        batchlist = list(batch)
        logging.info("Sending batch %s of size %d to %s", b_id, len(batchlist), host)
        payload = "{\"items\":[%s]}" % ",".join(batchlist)
        print(payload)
        conn.request("POST","/api/v1/tenants/%s/bulk/upsert/monitored-objects/meta" % tenantid, payload, {"Content-Type":"application/json","X-Forwarded-User-Roles":"skylight-admin"})
        bulk_response = conn.getresponse()
        if bulk_response.status != 200:
            logging.error("Batch with id %s failed to properly apply" % b_id)
            return
        bulk_response.read()
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

        logging.info("Processing csv with headers: \n%s", "\n".join(headers))
        
        for entry in bulkmetareader:
            batch += [entry]
            i += 1
            if i%batchsize == 0:
                f_envoy(i//batchsize, map(f_entry, batch))
                batch = []
        
        if len(batch) > 0:
            f_envoy((i//batchsize)+1, map(f_entry, batch))

start = time.time()

# Process the command line arguments
parser = argparse.ArgumentParser(description="Bulk insert meta information against monitored objects in datahub.")
parser.add_argument("-b", "--batchsize",type=int, help="Total size of a batch of metadata entries that should be sent to datahub")
parser.add_argument("-f", "--file", help="Absolute path to the csv file containing meta information")
parser.add_argument("-s", "--host", help="Host to send the metadata information to")
parser.add_argument("-u", "--username", help="Username to be used for logging into datahub")
parser.add_argument("-p", "--password", help="Password to be used for logging into datahub")
parser.add_argument("-t", "--tenantname", help="Name of the tenant that the monitored objects to be enriched are associated with")

args = parser.parse_args()

batchsize = args.batchsize
metafile = args.file
host = args.host
username = args.username
password = args.password
tenant = args.tenant
keyname = "Enode B"

conf_logging()

logging.info("Loading entries from " + metafile)

conn = http.client.HTTPSConnection(host)

logging.info("Logging into datahub...")
auth = login(conn, host, username, password)
if auth is None:
    logging.error("Could not login. Exiting...")

tid = tenant_id(conn, auth, host, tenant)

pf = open_processed_file()

try:
    logging.info("Starting to process...")
    process(metafile, batchsize, keyname, envoy_func(conn, auth, host, tid, pf))
finally:
    pf.close()

logging.info("Finishing processing %s in %s seconds", metafile, (time.time() - start))