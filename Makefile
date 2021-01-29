# Copyright (C) 2020 ConsenSys Software Inc.

# Build the Filecoin Retrieval Gateway

# Usage:
#   [VERSION=v3] [REGISTRY="gcr.io/google_containers"] make build
VERSION?=dev
REGISTRY?=

# This target (the first target in the build file) is the one that is executed if no 
# command line args are specified.
release: clean utest build

# builds a docker image that builds the app and packages it into a minimal docker image
build:
#	docker build -t ${REGISTRY}fc-retrieval-itest:${VERSION} .


# push the image to an registry
push:
#	gcloud docker -- push ${REGISTRY}/fc-retrieval-client:${VERSION}

uselocal:
	cd scripts; bash use-local-repos.sh

useremote:
	cd scripts; bash use-remote-repos.sh

detectmisconfig:
	cd scripts; bash detect-pkg-misconfig.sh

# Local build: make sure the test code compiles. 
lbuild:
	go test -c github.com/ConsenSys/fc-retrieval-itest/internal/integration


# Run the gateway(s), provider(s), and register services in Docker. Run the 
# tests locally. Dump the go.mod file so that the precise versions of 
# Client and Gateway Admin library are recorded. 
itest:
	docker-compose up -d
	echo *********************************************
	cat go.mod
	sleep 2
	echo *********************************************
	go test ./...
	docker-compose down
	
# This is the previous methodology, where the integration tests were in 
# a Docker container.
#
#	docker-compose down
#	docker-compose up --abort-on-container-exit --exit-code-from itest



# remove previous images and containers
clean:
	docker rmi -f "${REGISTRY}fc-retrieval-itest:${VERSION}" || true

# Alays assume these targets are out of date.
.PHONY: clean itest utest build release push detectmisconfig

