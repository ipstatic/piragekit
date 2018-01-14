BINARY = piragekit

VERSION?=?
COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS = -ldflags "-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.branch=${BRANCH}"

# Build the project
all: test fmt arm5

arm5:
	GOOS=linux GOARCH=arm GOARM=5 go build ${LDFLAGS} -o ${BINARY}-arm5

test:
	go test -v

fmt:
	go fmt $$(go list ./... | grep -v /vendor/)

.PHONY: arm5 fmt test
