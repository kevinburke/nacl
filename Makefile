SHELL = /bin/bash -o pipefail

BENCHSTAT := $(GOPATH)/bin/benchstat
BUMP_VERSION := $(GOPATH)/bin/bump_version
GODOCDOC := $(GOPATH)/bin/godocdoc
STATICCHECK := $(GOPATH)/bin/staticcheck

test: vet
	@# this target should always be listed first so "make" runs the tests.
	go test -trimpath ./...

$(STATICCHECK):
	GO111MODULE=on go install honnef.co/go/tools/cmd/staticcheck@latest

check: $(STATICCHECK)
	$(STATICCHECK) ./...

race-test:
	go test -trimpath -race ./...

vet:
	go vet -trimpath ./...

$(GODOCDOC):
	go get github.com/kevinburke/godocdoc

docs: $(GODOCDOC)
	$(GODOCDOC)

$(BENCHSTAT):
	go get golang.org/x/perf/cmd/benchstat

bench: $(BENCHSTAT)
	go test -trimpath -count=3 -benchtime=2s -bench=. -run='^$$' ./... | $(BENCHSTAT) /dev/stdin

$(BUMP_VERSION):
	go get github.com/kevinburke/bump_version

release: check race-test | $(BUMP_VERSION)
	$(BUMP_VERSION) minor nacl.go
