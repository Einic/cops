# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Setting SHELL to bash allows bash commands to be executed by recipes.
# This is a requirement for 'setup-envtest.sh' in the test target.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec


CWD := $(shell pwd)
VERSION := "v1.0.0"
#LDFLAGS := "-s -w -X main.version=${VERSION}"
LDFLAGS := "-s -w "
RELEASE = $(shell date +"%Y%m%d")
KERNEL_VERSION := $(shell uname -r|grep -oP '\.el(\d+)\.' | sed -n 's/\.el\([0-9]\+\)\./el\1/p')
ARCH := $(shell uname -m)

#CRON_SRC_FILE = cron_cdmasst.sh
#CRON_DEST_FILE = "./bin/cron_cdmasst.sh"
#CRON_COPY_CMD = cp $(CRON_SRC_FILE) $(CRON_DEST_FILE)
#
#ENV_SRC_FILE = install_daemon.sh
#ENV_DEST_FILE = "./bin/install_daemon.sh"
#ENV_COPY_CMD = cp $(ENV_SRC_FILE) $(ENV_DEST_FILE)

README_SRC_FILE = README.md
README_DEST_FILE = "./bin/README.md"
README_COPY_CMD = cp $(README_SRC_FILE) $(README_DEST_FILE)

.PHONY: all
all: build
    @echo "Building for Linux Kernel Version: $(KERNEL_VERSION)"
    @echo "Building for Architecture: $(ARCH)"

##@ General
ver:
	@echo ${VERSION}

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development
.PHONY: manifests
manifests: ## Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.

.PHONY: generate
generate: ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: test
test: manifests generate fmt vet envtest ## Run tests.
	go test ./... -coverprofile cover.out

##@ Build

.PHONY: linux-build
linux-build: generate fmt vet ## Build manager linux binary.
	go build -ldflags=${LDFLAGS} -o bin/cops main.go && find bin/cops -type f -executable | xargs -i upx -qq {}

.PHONY: darwin-build
darwin-build: generate fmt vet ## Build manager darwin binary.
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags=${LDFLAGS} -o bin/cops_darwin main.go && find bin/cops_darwin -type f -executable | xargs -i upx -qq {}

.PHONY: win-build
win-build: generate fmt vet ## Build manager windows binary.
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags=${LDFLAGS} -o bin/cops.exe main.go && find bin/cops.exe -type f -executable | xargs -i upx -qq {}

.PHONY: run
run: manifests generate fmt vet ## Run a controller from your host.
	go run ./main.go

.PHONY: build
build: linux-build ## Build manager linux rpm binary.
#	$(CRON_COPY_CMD)
#	$(ENV_COPY_CMD)
	$(README_COPY_CMD)
	sed -i "s/VERSION/$(VERSION)/g" rpm.json
	sed -i "s/RELEASE/$(RELEASE)/g" rpm.json
	mkdir -p rpms
	go-bin-rpm generate -f rpm.json -o ./rpms/cops-$(VERSION)-$(RELEASE).$(KERNEL_VERSION).$(ARCH).rpm

.PHONY: clean
clean:
	rm -rf pkg-build
	rm -rf logs
	rm -rf bin/README.md

##@ Deployment

ifndef ignore-not-found
  ignore-not-found = false
endif

#请使用TAB键替换空格，
#rpm 打包前安装
#1.https://github.com/mh-cbon/go-bin-rpm
#2.yum install -y gcc make rpm-build redhat-rpm-config