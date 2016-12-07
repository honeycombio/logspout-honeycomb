# Insturctions

## Running logspout-honeycomb
Run the logspout-honeycomb as follows

```
$ docker run -e "ROUTE_URIS=honeycomb://localhost" \
             -e "HONEYCOMB_WRITE_KEY=<api-key>" \
             -e "HONEYCOMB_DATASET=<dataset>" \
             -e "HONEYCOMB_SAMPLE_RATE=10" \
             -e "HONEYCOMB_API_URL=https://api.honeycomb.io" \
             --volume=/var/run/docker.sock:/var/run/docker.sock \
             --publish=127.0.0.1:8000:80 \
             logspout-honeycomb
```

* <api-key>: your api key
* <dataset>: dataset to log to

All the other environment variables are optional and have a default.

There are 2 requirements for a container to be logged:
* logging should be done to stdout and stderr
* container is started without the -t option (pseudo-TTY).

This is pretty straight forward and can be tested by running:
```
$ docker run mongo:latest
```
Voila! You just spun up a mongo container and you can see the output flowing into your dashboard. A lot of other docker images already push to stdout and stderr by default, such as:

* redis
* mysql
* nginx
* elasticsearch
* mongodb

## Ignoring containers

You can ignore specific containers by tell by setting an environment variable when starting your container:

```
$ docker run -d -e 'LOGSPOUT=ignore' image
```

or you can add a label which you define by setting an environment variable when running logspout:

```
$ docker run -e "ROUTE_URIS=honeycomb://localhost" \
             -e "HONEYCOMB_WRITE_KEY=<api-key>" \
             -e "HONEYCOMB_DATASET=<dataset>" \
             -e "HONEYCOMB_SAMPLE_RATE=10" \
             -e "HONEYCOMB_API_URL=https://api.honeycomb.io" \
             -e "EXCLUDE_LABEL=logspout.exclude" \
             --volume=/var/run/docker.sock:/var/run/docker.sock \
             --publish=127.0.0.1:8000:80 \
             logspout-honeycomb
$ docker run -d --label logspout.exclude=true mongo:latest
```

Since logspout-honeycomb is forked from the original logspout all the settings from [their instructions](https://github.com/gliderlabs/logspout) apply.