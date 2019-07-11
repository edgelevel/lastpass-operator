.DEFAULT_GOAL := all

CMD_DOCKER := $(shell command -v docker 2> /dev/null)
CMD_KUBECTL := $(shell command -v kubectl 2> /dev/null)
CMD_HELM := $(shell command -v helm 2> /dev/null)
CMD_GO := $(shell command -v go 2> /dev/null)
CMD_DEP := $(shell command -v dep 2> /dev/null)
CMD_OPERATOR_SDK := $(shell command -v operator-sdk 2> /dev/null)

DOCKER_USERNAME := edgelevel
DOCKER_IMAGE := $(DOCKER_USERNAME)/lastpass-operator

.PHONY: requirements
requirements:
ifndef CMD_DOCKER
	$(error "docker" not found)
endif
# ifndef CMD_KUBECTL
# 	$(error "kubectl" not found)
# endif
# ifndef CMD_HELM
# 	$(error "helm" not found)
# endif
ifndef CMD_GO
	$(error "go" not found)
endif
ifndef CMD_DEP
	$(error "dep" not found)
endif
ifndef CMD_OPERATOR_SDK
	$(error "operator-sdk" not found)
endif

.PHONY: test
test:
	dep ensure
	go test -v ./...

.PHONY: all
all: requirements test

.PHONY: docker-build
docker-build: all
	operator-sdk build $(DOCKER_IMAGE):${tag}

.PHONY: docker-login
docker-login:
	# prompt
	#docker login --username $(DOCKER_USERNAME)
	# without prompt
	echo ${docker-password} | docker login -u $(DOCKER_USERNAME) --password-stdin

.PHONY: docker-push
docker-push: docker-build docker-login
	docker tag $(DOCKER_IMAGE):${tag} $(DOCKER_IMAGE):latest
	docker push $(DOCKER_IMAGE):${tag}
	docker push $(DOCKER_IMAGE):latest
