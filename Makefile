GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
IMAGE=qkrqjadn/sumwhere
DOCKER_PASSWORD=1q2w3e4r
DOCKER_USERNAME=qkrqjadn
VERSION=1.0
GITCOMMITCOUNT:=$$(git rev-list HEAD | wc -l | tr -d ' ')
GITHASH:=$$(git rev-parse --short HEAD)
DATETIME:=$$(date "+%Y%m%d")
VERSIONS:=$(VERSION).$(GITCOMMITCOUNT)-$(GITHASH)-$(DATETIME)

.PHONY: auth api test clean
run_auth:
	$(GOCMD) build -o $@ cmd/auth/auth.go

run_api:
	$(GOCMD) build -i $@ cmd/api/api.go

test:
	$(GOCMD) test -race ./...

clean:
	rm -f run_auth run_api
	$(GOCLEAN)