SHELL = /bin/bash

BENCHSTAT := $(shell command -v benchstat)
BUMP_VERSION := $(shell command -v bump_version)
GODOCDOC := $(shell command -v godocdoc)
MEGACHECK := $(shell command -v megacheck)

test: vet
	@# this target should always be listed first so "make" runs the tests.
	bazel test --test_output=errors //...

megacheck:
ifndef MEGACHECK
	go get -u honnef.co/go/tools/cmd/megacheck
endif
	go list ./... | grep -v vendor | xargs megacheck --ignore='github.com/kevinburke/nacl/*/*.go:S1002'

race-test: megacheck vet
	bazel test --test_output=errors --features=race //...

ci:
	bazel --batch --host_jvm_args=-Dbazel.DigestFunction=SHA1 test \
		--experimental_repository_cache="$$HOME/.bzrepos" \
		--spawn_strategy=remote \
		--remote_rest_cache=https://remote.rest.stackmachine.com/cache \
		--test_output=errors \
		--strategy=Javac=remote \
		--features=race //... 2>&1 | ts '[%Y-%m-%d %H:%M:%.S]'

vet:
	go list ./... | grep -v vendor | xargs go vet

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

release: megacheck race-test
ifndef BUMP_VERSION
	go get github.com/Shyp/bump_version
endif
	bump_version minor nacl.go
