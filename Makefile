NAME=logspout-honeycomb
BUILD_DIR=build

build:
	mkdir $(BUILD_DIR)
	git clone https://github.com/gliderlabs/logspout.git $(BUILD_DIR)/logspout
	# Copy this repo's files into logspout checkout, so it can find them for
	# its Docker build.
	mkdir $(BUILD_DIR)/logspout/build-logspout-honeycomb
	cp -v *.go $(BUILD_DIR)/logspout/build-logspout-honeycomb/.
	mkdir $(BUILD_DIR)/logspout/build-libhound
	cp -rv ../libhound-go-private/* $(BUILD_DIR)/logspout/build-libhound
	# Modify the Docker build to copy in our private repos
	patch $(BUILD_DIR)/logspout/Dockerfile < logspout-mods/docker.diff
	# Modify Logspout module file to use Honeycomb adapter
	cp -v logspout-mods/modules.go $(BUILD_DIR)/logspout/.
	docker build $(BUILD_DIR)/logspout -t logspout-honeycomb

clean:
	rm -rf $(BUILD_DIR)
