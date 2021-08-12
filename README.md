# logspout-honeycomb

[![OSS Lifecycle](https://img.shields.io/osslifecycle/honeycombio/logspout-honeycomb)](https://github.com/honeycombio/home/blob/main/honeycomb-oss-lifecycle-and-practices.md)

Honeycomb adapter for Logspout. More documentation can be found [in Honeycomb docs](https://honeycomb.io/docs/connect/logspout/).

Expects to ingest JSON log lines, and will send JSON blobs up to Honeycomb, annotated with the current logspout stream, container, container ID, hostname, and docker image name.

If the log lines being streamed through Logspout aren't JSON, the contents of the message will be tucked under a `"message"` key in the Honeycomb payload, alongside the metadata mentioned above.

## Building

To build the Honeycomb Logspout Docker image, run:
* `make docker`

## Configuration and invocation

This module can be configured either by setting environment variables in
Docker, or by using the [Logspout routesapi](https://github.com/gliderlabs/logspout/tree/master/routesapi). The following variables are available:

Env. Variable | routesapi key | Type | Required? | Description |
| --- | --- | --- | --- | -----|
| `HONEYCOMB_WRITE_KEY` | `writeKey` | string | required | Your Honeycomb team's write key. |
| `HONEYCOMB_DATASET` | `dataset` | string | required | The name of the destination dataset in your Honeycomb account. It will be created if it does not already exist. |
| `HONEYCOMB_SAMPLE_RATE` | `sampleRate` | integer | optional | Sample your event stream: send 1 out of every N events |

### Environment variables

Configure the logspout-honeycomb image via environment variables and run the container:

    docker run \
        -e "ROUTE_URIS=honeycomb://localhost" \
        -e "HONEYCOMB_WRITE_KEY=<YOUR_WRITE_KEY>" \
        -e "HONEYCOMB_DATASET=<YOUR_DATASET>" \
        --volume=/var/run/docker.sock:/var/run/docker.sock \
        --publish=127.0.0.1:8000:80 \
        honeycombio/logspout-honeycomb:1.13

### routesapi

Configuration can be set after the logspout-honeycomb image is already running via routesapi:

    docker run \
        --volume=/var/run/docker.sock:/var/run/docker.sock \
        --publish=127.0.0.1:8000:80 \
        honeycombio/logspout-honeycomb:1.13

    curl $(docker port `docker ps -lq` 80)/routes \
        -X POST \
            -d '{"adapter": "honeycomb",
                 "address": "honeycomb://localhost",
                 "options": {"writeKey":"<YOUR_WRITE_KEY>",
                             "dataset":"<YOUR_DATASET>"}}'
