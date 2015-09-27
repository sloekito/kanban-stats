name=kanban-stats
account=docker.cloud.nlab.io
release=0.0.9
target_os=linux
target_arch=amd64
influx_release = 0.9.4.1
download_url = http://s3.amazonaws.com/influxdb/influxdb_$(influx_release)_amd64.deb

build/container: build dist/influxdb_$(influx_release)_amd64.deb
	docker build --no-cache -t $(name) .
	touch build/container

.PHONY: release	
release: build/container
	docker tag -f $(name) $(account)/$(name):$(release)
	docker push $(account)/$(name):$(release)
	
.PHONY: clean
clean:
	rm -rf build
	rm -rf dist
	
.PHONY: build
build:
	mkdir -p build
	GOOS=$(target_os) GOARCH=$(target_arch) go build -o build/$(name)
	
dist:
	mkdir -p dist
	
dist/influxdb_$(influx_release)_amd64.deb: dist
	curl -sLo dist/influxdb_$(influx_release)_amd64.deb $(download_url)

.PHONY: toolchain
toolchain:
	go get -v github.com/mitchellh/gox
	gox -osarch=$(target_os)/$(target_arch) -build-toolchain