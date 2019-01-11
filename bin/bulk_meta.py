import argparse
import calendar
import csv
import http.client, urllib.parse
import json
import logging
import os
import sys
import time

# Builds up a basic console logger configuration
def conf_logging():
    logging.basicConfig(format='%(asctime)s:%(levelname)s: %(message)s', level=logging.INFO)
    logging.info("Logging configured...")

# Initializes a processing file that indicates the status of each monitored object that we attempt to attach metadata to
def open_processed_file():
    processed_headers = ["key","status"]
    processed_filename = "/tmp/processed.csv"
    logging.info("Generating processed log at %s", processed_filename)
    f = open(processed_filename,"w+")
    f.write(",".join(processed_headers)+"\n")
    return f

# Returns a configured function that knows how to take the headers from the csv file 
# and build up a bulk metadata entry that the API expects
def construct_entry_func(headers):
    # Returns a json representations of the metadata structure for a single entry
    def construct_entry(entry): 
        objectName = entry.pop(0)
        metadata = {k: v for k, v in dict(zip(headers,entry)).items() if v}
        entry_struct = { "objectName":objectName,"metadata":metadata}
        return json.dumps(entry_struct)
    return construct_entry

# Returns a configured function that knows how to send a batch of monitoredobjects->metadata to the handler
def envoy_func(conn, auth, host, processfile):
    def envoy(b_id, batch):
        batchlist = list(batch)
        logging.info("Sending batch %s of size %d to %s", b_id, len(batchlist), host)
        # Build up the batch and send
        payload = '{"data":{"type":"monitoredObjectsMeta","attributes":{"metadata-entries":[%s]}}}' % ",".join(batchlist)
        conn.request("POST","/api/v2/bulk/insert/monitored-objects/meta", payload, {"Content-Type":"application/json","Authorization":auth})
        bulk_response = conn.getresponse()
        
        # Convert the response to json and process each individual response entry to track in the process file
        response = json.loads(bulk_response.read().decode("utf-8"))
        for entry in response.get("data"):
            attributes = entry.get("attributes")
            status = "200" if attributes.get("ok") == True else attributes.get("reason")
            processfile.write("%s,%s\n" % (attributes.get("id"), status))
    return envoy

# Processes login so that the script can use the bulk API
def login(conn, host, username, password):
    conn.request("POST","/api/v1/auth/login",urllib.parse.urlencode({"username":username,"password":password}),{"Content-Type":"application/x-www-form-urlencoded"})
    login_response = conn.getresponse()
    if login_response.status != 200:
        logging.error("Could not login to host %s", host)
        return
    login_response.read()
    return login_response.getheader("Authorization")

# Processes the CSV file and sends the metadata contained within to the mapped monitored object name
def process(file, batchsize, f_envoy):
    with open(file) as csvfile:
        bulkmetareader = csv.reader(csvfile, delimiter=',')
        i = 0
        headers = bulkmetareader.__next__()
        f_entry = construct_entry_func(headers)
        batch = []

        # Remove the first entry as it is assumed to be the identifying object name column
        headers.pop(0)

        logging.info("Processing csv with headers: \n%s", "\n".join(headers))
        
        for entry in bulkmetareader:
            batch += [entry]
            i += 1
            if i%batchsize == 0:
                f_envoy(i//batchsize, map(f_entry, batch))
                batch = []
        
        # Process the remaining records in their own batch
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

args = parser.parse_args()

batchsize = args.batchsize
metafile = args.file
host = args.host
username = args.username
password = args.password

conf_logging()

logging.info("Loading entries from " + metafile)

conn = http.client.HTTPSConnection(host)

logging.info("Logging into datahub...")
auth = login(conn, host, username, password)
if auth is None:
    logging.error("Could not login. Exiting...")
    sys.exit()

pf = open_processed_file()

try:
    logging.info("Starting to process...")
    process(metafile, batchsize, envoy_func(conn, auth, host, pf))
finally:
    pf.close()

logging.info("Finished processing %s in %s seconds", metafile, (time.time() - start))
