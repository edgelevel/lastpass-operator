#!/bin/bash

DOCKER_TAG=${1:?"Missing DOCKER_TAG"}
DOCKER_PASSWORD=${2:?"Missing DOCKER_PASSWORD"}

make docker-push tag=$DOCKER_TAG docker-password=$DOCKER_PASSWORD
