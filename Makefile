.DEFAULT_GOAL := all

CMD_DOCKER := $(shell command -v docker 2> /dev/null)
CMD_KUBECTL := $(shell command -v kubectl 2> /dev/null)
CMD_HELM := $(shell command -v helm 2> /dev/null)
CMD_GO := $(shell command -v go 2> /dev/null)
CMD_DEP := $(shell command -v dep 2> /dev/null)
CMD_OPERATOR_SDK := $(shell command -v operator-sdk 2> /dev/null)

.PHONY: requirements
requirements:
ifndef CMD_DOCKER
	$(error "docker" not found)
endif
ifndef CMD_KUBECTL
	$(error "kubectl" not found)
endif
ifndef CMD_HELM
	$(error "helm" not found)
endif
ifndef CMD_GO
	$(error "go" not found)
endif
ifndef CMD_DEP
	$(error "dep" not found)
endif
ifndef CMD_OPERATOR_SDK
	$(error "operator-sdk" not found)
endif

.PHONY: all
all: requirements
