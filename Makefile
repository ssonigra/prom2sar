# Image URL to use all building/pushing image targets
IMG ?= prometheus-dump-operator:latest
NAMESPACE ?= prometheus-dump-operator

# Get the currently used golang install path
GOPATH ?= $(shell go env GOPATH)
GOBIN ?= $(GOPATH)/bin

.PHONY: all
all: build

##@ General

.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: fmt
fmt: ## Run go fmt
	go fmt ./...

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: test
test: fmt vet ## Run tests
	go test ./... -coverprofile cover.out

##@ Build

.PHONY: build
build: fmt vet ## Build operator binary
	go build -o bin/prometheus-dump-operator cmd/main.go

.PHONY: build-cli
build-cli: fmt vet ## Build standalone CLI binary (prom2sar)
	go build -o bin/prom2sar cmd/prom2sar/main.go

.PHONY: build-all
build-all: build build-cli ## Build both operator and CLI

.PHONY: install-cli
install-cli: build-cli ## Install CLI to /usr/local/bin (requires sudo)
	sudo cp bin/prom2sar /usr/local/bin/
	sudo chmod +x /usr/local/bin/prom2sar
	@echo "prom2sar installed to /usr/local/bin/"

.PHONY: run
run: fmt vet ## Run operator locally
	go run cmd/main.go

.PHONY: docker-build
docker-build: ## Build docker image
	docker build -t ${IMG} .

.PHONY: docker-push
docker-push: ## Push docker image
	docker push ${IMG}

##@ Deployment

.PHONY: install
install: ## Install CRDs
	kubectl apply -f deploy/crds/

.PHONY: uninstall
uninstall: ## Uninstall CRDs
	kubectl delete -f deploy/crds/

.PHONY: deploy
deploy: ## Deploy operator to cluster
	kubectl create namespace ${NAMESPACE} --dry-run=client -o yaml | kubectl apply -f -
	kubectl apply -f deploy/

.PHONY: undeploy
undeploy: ## Remove operator from cluster
	kubectl delete -f deploy/

##@ Examples

.PHONY: example-basic
example-basic: ## Apply basic SAR conversion example
	kubectl apply -f examples/basic-sar-conversion.yaml

.PHONY: example-cpu
example-cpu: ## Apply CPU-only example
	kubectl apply -f examples/cpu-only-sar.yaml

.PHONY: example-custom
example-custom: ## Apply custom metrics example
	kubectl apply -f examples/custom-metrics-sar.yaml

.PHONY: examples-clean
examples-clean: ## Remove all example CRs
	kubectl delete -f examples/ --ignore-not-found

##@ Utilities

.PHONY: logs
logs: ## Show operator logs
	kubectl logs -n ${NAMESPACE} -l app=prometheus-dump-operator -f

.PHONY: status
status: ## Show status of all PrometheusDumpLoader resources
	kubectl get promethedusdumploaders -A

.PHONY: clean
clean: ## Clean build artifacts
	rm -rf bin/
	go clean
