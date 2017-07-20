SHELL = /bin/bash

BENCHSTAT := $(shell command -v benchstat)
MEGACHECK := $(shell command -v megacheck)

test: vet
	@# this target should always be listed first so "make" runs the tests.
	go list ./... | grep -v vendor | xargs go test -short

race-test: vet
	go list ./... | grep -v vendor | xargs go test -race

vet:
ifndef MEGACHECK
	go get -u honnef.co/go/tools/cmd/megacheck
endif
	go list ./... | grep -v vendor | xargs go vet
	go list ./... | grep -v vendor | xargs megacheck --ignore='github.com/kevinburke/nacl/*/*.go:S1002'

docs:
ifndef GODOCDOC
	go get github.com/kevinburke/godocdoc
endif
	godocdoc

bench:
ifndef BENCHSTAT
	go get rsc.io/benchstat
endif
	tmp=$$(mktemp); go list ./... | grep -v vendor | xargs go test -benchtime=2s -bench=. -run='^$$' > "$$tmp" 2>&1 && benchstat "$$tmp"
