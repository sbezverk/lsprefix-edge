REGISTRY_NAME?=docker.io/sbezverk
IMAGE_VERSION?=0.0.0

.PHONY: all lsprefix-edge container push clean test

ifdef V
TESTARGS = -v -args -alsologtostderr -v 5
else
TESTARGS =
endif

all: lsprefix-edge

lsprefix-edge:
	mkdir -p bin
	$(MAKE) -C ./cmd compile-lsprefix-edge

lsprefix-edge-container: lsprefix-edge
	docker build -t $(REGISTRY_NAME)/lsprefix-edge:$(IMAGE_VERSION) -f ./build/Dockerfile.lsprefix-edge .

push: lsprefix-edge-container
	docker push $(REGISTRY_NAME)/lsprefix-edge:$(IMAGE_VERSION)

clean:
	rm -rf bin

test:
	GO111MODULE=on go test `go list ./... | grep -v 'vendor'` $(TESTARGS)
	GO111MODULE=on go vet `go list ./... | grep -v vendor`
