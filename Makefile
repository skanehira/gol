# Go parameters
GOBUILD=go build
GOCLEAN=go clean
BINARY_NAME=gol
DOCKER_BINARY_NAME=gol-docker

export GO111MODULE=on

all: build

clean:
	$(GOCLEAN)

build: clean
	$(GOBUILD) -o $(BINARY_NAME)

# copy to $GOBIN
install: build
	cp -f $(BINARY_NAME) $(GOBIN)/

# build release binary
release: clean
	GOOS=darwin GOARCH=amd64 $(GOBUILD) && zip MacOS.zip $(BINARY_NAME) && rm -rf $(BINARY_NAME)
