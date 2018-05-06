SHELL = /bin/bash -o pipefail

BENCHSTAT := $(GOPATH)/bin/benchstat
BUMP_VERSION := $(GOPATH)/bin/bump_version
GODOCDOC := $(GOPATH)/bin/godocdoc
MEGACHECK := $(GOPATH)/bin/megacheck

test: vet
	@# this target should always be listed first so "make" runs the tests.
	go test ./...

$(MEGACHECK):
	go get honnef.co/go/tools/cmd/megacheck

check: $(MEGACHECK)
	go list ./... | grep -v vendor | xargs $(MEGACHECK) --ignore='github.com/kevinburke/nacl/*/*.go:S1002'

race-test: check vet
	go test -race ./...

vet:
	go list ./... | grep -v vendor | xargs go vet

$(GODOCDOC):
	go get github.com/kevinburke/godocdoc

docs: $(GODOCDOC)
	$(GODOCDOC)

$(BENCHSTAT):
	go get golang.org/x/perf/cmd/benchstat

bench: $(BENCHSTAT)
	tmp=$$(mktemp); go list ./... | grep -v vendor | xargs go test -benchtime=2s -bench=. -run='^$$' > "$$tmp" 2>&1 && $(BENCHSTAT) "$$tmp"

$(BUMP_VERSION):
	go get github.com/kevinburke/bump_version

release: check race-test | $(BUMP_VERSION)
	$(BUMP_VERSION) minor nacl.go
