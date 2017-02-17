NAME=logspout-honeycomb
BUILD_DIR=build

# If you want to configure the Honeycomb Logspout adapter with environment
# variables, set them here. Otherwise you need to use the RoutesAPI.
HONEYCOMB_WRITE_KEY=
HONEYCOMB_DATASET=

# Builds a Docker image of LogSpout with the Honeycomb adapter included.
docker: honeycomb.go
	docker build -t $(NAME) docker

# Fire up a container with the Honeycomb Logspout adapter in it,
# configured by environment variables
run-with-env:
	docker run \
		-e "ROUTE_URIS=honeycomb://localhost" \
		-e "HONEYCOMB_WRITE_KEY=$(HONEYCOMB_WRITE_KEY)" \
		-e "HONEYCOMB_DATASET=$(HONEYCOMB_DATASET)" \
		--volume=/var/run/docker.sock:/var/run/docker.sock \
		--publish=127.0.0.1:8000:80 \
		$(NAME)

# Fire up a container with the Honeycomb Logspout adapter in it
run:
	docker run \
		--volume=/var/run/docker.sock:/var/run/docker.sock \
		--publish=127.0.0.1:8000:80 \
		$(NAME)

# Launches a Docker image that logs random JSON messages every second
random-source:
	docker run alpine /bin/sh -c \
		'while true; do echo {\"random\": `echo $$RANDOM`}; sleep 1; done'

clean:
	rm -rf $(BUILD_DIR)
	docker rmi -f $(NAME)

clean-images:
	docker images | grep none | awk '{print $3}' | xargs docker rmi -f
