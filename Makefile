all: gofmt vet build-binary-in-docker run-test-in-docker clean cross

override REGISTRY := $(or ,$(REGISTRY), "eu.gcr.io/managed-certs-gke")

TAG ?= $(USER)-dev
ARTIFACTS ?= /tmp/artifacts

# Latest commit hash for current branch
GIT_COMMIT ?= $(shell git rev-parse HEAD)
# This version-strategy uses git tags to set the version string
VERSION ?= $(shell git describe --tags --always --dirty)

name=managed-certificate-controller
runner_image=$(name)-runner
runner_path=/gopath/src/github.com/GoogleCloudPlatform/gke-managed-certs/

auth-configure-docker:
	test -f /etc/service-account/service-account.json && \
		gcloud auth activate-service-account --key-file=/etc/service-account/service-account.json && \
		gcloud auth configure-docker || true

# Builds the managed certs controller binary
build-binary: clean
	pkg=github.com/GoogleCloudPlatform/gke-managed-certs; \
	ld_flags="-X $${pkg}/pkg/version.Version=$(VERSION) -X $${pkg}/pkg/version.GitCommit=$(GIT_COMMIT)"; \
	go build -o $(name) -ldflags "$${ld_flags}"

# Builds the managed-certificate-controller binary using a docker runner image
build-binary-in-docker: docker-runner-builder
	docker run -v `pwd`:$(runner_path) $(runner_image):latest bash -c 'cd $(runner_path) && make build-binary GIT_COMMIT=$(GIT_COMMIT) VERSION=$(VERSION)'

clean:
	rm -f $(name)

# Checks if Google criteria for releasing code as OSS are met
cross:
	if [ -e /google/data/ro/teams/opensource/cross ]; then /google/data/ro/teams/opensource/cross .; fi

# Builds and pushes a docker image with managed-certificate-controller binary
docker: auth-configure-docker
	until docker build --pull -t $(REGISTRY)/$(name):$(TAG) -t $(REGISTRY)/$(name):$(VERSION) .; do \
		echo "Building managed-cetrificate-controller image failed, retrying in 10 seconds..." && sleep 10; \
	done
	until docker push $(REGISTRY)/$(name):$(TAG); do \
		echo "Pushing managed-certificate-controller image failed, retrying in 10 seconds..." && sleep 10; \
	done
	until docker push $(REGISTRY)/$(name):$(VERSION); do \
		echo "Pushing managed-certificate-controller image failed, retrying in 10 seconds..." && sleep 10; \
	done

# Builds a runner image, i. e. an image used to build a managed-certificate-controller binary and to run its tests.
docker-runner-builder:
	until docker build -t $(runner_image) runner; do \
		echo "Building runner image failed, retrying in 10 seconds..." && sleep 10; \
	done

e2e:
	dest=/tmp/artifacts; \
	rm -rf $${dest}/* && mkdir -p $${dest} && \
	{ \
		CLOUD_SDK_ROOT=$(CLOUD_SDK_ROOT) \
		KUBECONFIG=$(HOME)/.kube/config \
		KUBERNETES_PROVIDER=$(KUBERNETES_PROVIDER) \
		PROJECT=$(PROJECT) \
		DNS_ZONE=$(DNS_ZONE) \
		DOMAIN=$(DOMAIN) \
		PLATFORM=$(PLATFORM) \
		TAG=$(TAG) \
		REGISTRY=$(REGISTRY) \
		go test ./e2e/... -test.timeout=60m \
			-logtostderr=false -alsologtostderr=true -v -log_dir=$${dest} \
			> $${dest}/e2e.out.txt && exitcode=$${?} || exitcode=$${?} ; \
	} && cat $${dest}/e2e.out.txt | go-junit-report > $${dest}/junit_01.xml && exit $${exitcode}

# Formats go source code with gofmt
gofmt:
	find . -type f -name '*.go' | grep -v '/vendor/' | xargs gofmt -w

# Builds the managed certs controller binary, then a docker image with this binary, and pushes the image, for dev
release: release-ci clean

# Builds the managed certs controller binary, then a docker image with this binary, and pushes the image, for continuous integration
release-ci: build-binary-in-docker run-test-in-docker docker

run-e2e-in-docker: docker-runner-builder auth-configure-docker
	$(eval cloud_config := $(or ,$(CLOUD_CONFIG), $(shell gcloud info --format="value(config.paths.global_config_dir)")))
	$(eval cloud_sdk_root := $(or ,$(CLOUD_SDK_ROOT), $(shell gcloud info --format="value(installation.sdk_root)")))
	$(eval dns_zone := $(or ,$(DNS_ZONE),"managedcertsgke"))
	$(eval kubeconfig := $(or ,$(KUBECONFIG),$(HOME)/.kube/config))
	$(eval kubernetes_provider := $(or ,$(KUBERNETES_PROVIDER),"gke"))
	$(eval platform := $(or ,$(PLATFORM), "gcp"))
	$(eval project := $(or ,$(PROJECT), $(shell gcloud config list --format="value(core.project)")))
	docker run -v `pwd`:$(runner_path) \
		-v $(cloud_sdk_root):$(cloud_sdk_root) \
		-v $(cloud_config):/root/.config/gcloud \
		-v $(cloud_config):/root/.config/gcloud-staging \
		-v $(kubeconfig):/root/.kube/config \
		-v $(ARTIFACTS):/tmp/artifacts \
		$(runner_image):latest bash -c 'cd $(runner_path) && make e2e \
		DNS_ZONE=$(dns_zone) DOMAIN=$(DOMAIN) CLOUD_SDK_ROOT=$(cloud_sdk_root) \
		KUBERNETES_PROVIDER=$(kubernetes_provider) PROJECT=$(project) \
		PLATFORM=$(platform) REGISTRY=$(REGISTRY) TAG=$(TAG)'

run-test-in-docker: docker-runner-builder
	docker run -v `pwd`:$(runner_path) $(runner_image):latest bash -c 'cd $(runner_path) && make test'

test:
	go test ./pkg/... -cover

vet:
	go vet ./...

.PHONY: all auth-configure-docker build-binary build-binary-in-docker build-dev clean cross docker docker-runner-builder e2e release release-ci run-e2e-in-docker run-test-in-docker test vet
