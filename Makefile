NAME=logspout-honeycomb
BUILD_DIR=build

# If you want to configure the Honeycomb Logspout adapter with environment
# variables, set them here. Otherwise you need to use the RoutesAPI.
HONEYCOMB_WRITE_KEY=
HONEYCOMB_DATASET=

# Builds a Docker image of LogSpout with the Honeycomb adapter included.
docker: honeycomb.go
	mkdir $(BUILD_DIR)
	# Clone Logspout code, which we need to build a Docker image.
	git clone https://github.com/gliderlabs/logspout.git $(BUILD_DIR)/logspout
	# Copy this repo's files into logspout checkout, so it can find them for
	# its Docker build. Otherwise, 'go get' fails to checkout our private
	# repos because it can't auth in the Docker container.
	mkdir $(BUILD_DIR)/logspout/build-logspout-honeycomb
	cp -v *.go $(BUILD_DIR)/logspout/build-logspout-honeycomb/.
	git clone https://github.com/honeycombio/libhoney-go $(BUILD_DIR)/logspout/build-libhoney/
	# Modify the Docker build to copy in our private repos
	patch $(BUILD_DIR)/logspout/Dockerfile < logspout-mods/docker.diff
	# Modify Logspout module file to use Honeycomb adapter
	cp -v logspout-mods/modules.go $(BUILD_DIR)/logspout/.
	docker build $(BUILD_DIR)/logspout -t $(NAME)

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
