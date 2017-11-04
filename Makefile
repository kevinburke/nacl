SHELL = /bin/bash -o pipefail

BAZEL_VERSION := 0.7.0
BAZEL_DEB := bazel_$(BAZEL_VERSION)_amd64.deb

BENCHSTAT := $(GOPATH)/bin/benchstat
BUMP_VERSION := $(GOPATH)/bin/bump_version
GODOCDOC := $(GOPATH)/bin/godocdoc
MEGACHECK := $(GOPATH)/bin/megacheck

test: vet
	@# this target should always be listed first so "make" runs the tests.
	bazel test --test_output=errors //...

$(MEGACHECK):
	go get honnef.co/go/tools/cmd/megacheck

install-travis:
	wget "https://storage.googleapis.com/bazel-apt/pool/jdk1.8/b/bazel/$(BAZEL_DEB)"
	sudo dpkg --force-all -i $(BAZEL_DEB)
	sudo apt-get install moreutils -y

check: $(MEGACHECK)
	go list ./... | grep -v vendor | xargs $(MEGACHECK) --ignore='github.com/kevinburke/nacl/*/*.go:S1002'

race-test: check vet
	bazel test --test_output=errors --features=race //...

ci:
	bazel --batch --host_jvm_args=-Dbazel.DigestFunction=SHA1 test \
		--experimental_repository_cache="$$HOME/.bzrepos" \
		--spawn_strategy=remote \
		--test_output=errors \
		--noshow_progress --noshow_loading_progress \
		--strategy=Javac=remote \
		--features=race //... 2>&1 | ts '[%Y-%m-%d %H:%M:%.S]'

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
	go get github.com/Shyp/bump_version

release: check race-test | $(BUMP_VERSION)
	$(BUMP_VERSION) minor nacl.go
