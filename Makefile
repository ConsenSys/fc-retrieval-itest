# Copyright (C) 2020 ConsenSys Software Inc.

# Build the Filecoin Retrieval Gateway

# Usage:
#   [VERSION=v3] [REGISTRY="gcr.io/google_containers"] make build
VERSION?=dev
IMAGE?=consensys/fc-retrieval-itest
COMPOSE_FILE?=docker-compose.yml


# This target (the first target in the build file) is the one that is executed if no 
# command line args are specified.
default: clean utest build tag

# builds a docker image that builds the app and packages it into a minimal docker image
build:
	docker build -t ${IMAGE}:${VERSION} .


# push the image to an registry
push:
	cd scripts; bash push.sh ${VERSION} ${IMAGE}:${VERSION}

tag:
	cd scripts; bash tag.sh ${VERSION} ${IMAGE}:${VERSION}

uselocal:
	cd scripts; bash use-local-repos.sh

useremote:
	cd scripts; bash use-remote-repos.sh

detectmisconfig:
	cd scripts; bash detect-pkg-misconfig.sh

# Local build: make sure the test code compiles. 
lbuild:
	go test -c github.com/ConsenSys/fc-retrieval-itest/internal/integration

itestlocal: setup-env-localtesting itestdocker

setup-env-localtesting:
	cd scripts; bash setup-env.sh

# Version that can be run on a desktop computer or in Circle CI.
# Itest run from a container.
# Run the gateway(s), provider(s), and register services in Docker. Run the 
# tests locally. Dump the go.mod file so that the precise versions of 
# Client and Gateway Admin library are recorded. 
itestdocker:
	docker network create shared || true
	docker-compose down
	docker-compose -f $(COMPOSE_FILE) up -d gateway provider register redis 
	echo *********************************************
	cat go.mod
	sleep 10
	echo REDIS STARTUP *********************************************
	docker container logs redis
	echo REGISTER STARTUP *********************************************
	docker container logs register
	echo GATEWAY STARTUP *********************************************
	docker container logs gateway
	echo PROVIDER STARTUP *********************************************
	docker container logs provider
	echo NETWORK CONFIG *********************************************
	docker network inspect shared
	echo *********************************************
	docker-compose up itest
	echo *********************************************
	echo REDIS LOGS *********************************************
	docker container logs redis
	echo REGISTER LOGS *********************************************
	docker container logs register
	echo GATEWAY LOGS *********************************************
	docker container logs gateway
	echo PROVIDER LOGS *********************************************
	docker container logs provider
	echo ITEST LOGS *********************************************
	docker container logs itest
	echo *********************************************
	docker-compose down
	
# This is the previous methodology, where the integration tests were in 
# a Docker container.
#
#	docker-compose down
#	docker-compose up --abort-on-container-exit --exit-code-from itest


clean:
	docker rm -f ${IMAGE}:${VERSION} 2> /dev/null || true
	docker rmi -f ${IMAGE}:${VERSION} || true

check-modules:
	./scripts/check-modules/check-modules

check-main-modules:
	./scripts/check-main-modules/check-main-modules

# Alays assume these targets are out of date.
.PHONY: clean itest itest-dev utest build release push detectmisconfig

