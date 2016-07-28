# logspout-honeycomb
Honeycomb adapter for LogSpout.

# Building

To build the Honeycomb LogSpout Docker image, run:
* `make docker`

# Configuration

This module can be configured either by setting environment variables in
Docker, or by using the [LogSpout RoutesAPI](https://github.com/gliderlabs/logspout/tree/master/routesapi). There are 4 variables to consider:

* `WriteKey` (string) Your Honeycomb account's API key.
* `Dataset` (string) The name of the destination dataset in your Honeycomb account. It will be created if it does not already exist.
* `SampleRate` (optional, integer) Only send 1 out of N events
* `ApiUrl` (optional, URL string) An alternate Honeycomb API endpoint to send events to. Debugging purposes only.

## Environment variables

    docker run \
        -e "ROUTE_URIS=honeycomb://localhost" \
        -e "HONEYCOMB_WRITE_KEY=abcdefg12345678" \
        -e "HONEYCOMB_DATASET=myDataset" \
        -e "HONEYCOMB_SAMPLE_RATE=10" \
        -e "HONEYCOMB_API_URL=https://api.hound.sh" \
        --volume=/var/run/docker.sock:/var/run/docker.sock \
        --publish=127.0.0.1:8000:80 \
        logspout-honeycomb

## RoutesAPI

    curl $(docker port `docker ps -lq` 80)/routes \
        -X POST \
            -d '{"adapter": "honeycomb",
                 "address": "honeycomb://localhost",
                 "options": {"writeKey":"abcdefg12345678",
                             "dataset":"mydataset",
                             "sampleRate":10}}'
