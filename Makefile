name=daily_trello
account=docker.cloud.nlab.io
release=0.0.4
target_os=linux
target_arch=amd64

build/container: build
	docker build --no-cache -t $(name) .
	touch build/container

.PHONY: release	
release: build/container
	docker tag -f $(name) $(account)/$(name):$(release)
	docker push $(account)/$(name):$(release)
	
.PHONY: clean
clean:
	rm -rf build
	
.PHONY: build
build:
	mkdir -p build
	GOOS=$(target_os) GOARCH=$(target_arch) go build -o build/$(name)
	
dist:
	mkdir -p dist

.PHONY: toolchain
toolchain:
	go get -v github.com/mitchellh/gox
	gox -osarch=$(target_os)/$(target_arch) -build-toolchain