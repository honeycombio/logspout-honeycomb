# logspout-honeycomb
Honeycomb adapter for LogSpout

# Building

To build the Honeycomb LogSpout Docker image, run:
* `make build`

# Configuration

This module can be configured either by setting environment variables in
Docker, or by using the [LogSpout RoutesAPI](https://github.com/gliderlabs/logspout/tree/master/routesapi).

There are 4 variables to consider:
* Write key
 * __String__. Your Honeycomb account's API key.
* Dataset
 * __String__. The name of the destination dataset in your Honeycomb account. It will be created if it does not already exist.
* Sample rate (optional)
 * __Integer__. Only send 1 out of N events
* API URL (optional)
 * __URL String__. An alternate Honeycomb API endpoint to send events to. Debugging purposes only.

## Environment variables

    docker run \
        -e "ROUTE_URIS=honeycomb://localhost" \
        -e "HONEYCOMB_WRITE_KEY=abcdefg12345678" \
        -e "HONEYCOMB_DATASET=myDataset" \
        -e "HONEYCOMB_SAMPLE_RATE=10" \
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
