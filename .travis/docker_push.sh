#!/bin/bash

# defined in travis ui
DOCKER_PASSWORD=${1:?"Missing DOCKER_PASSWORD"}

make docker-push tag=$TRAVIS_TAG docker-password=$DOCKER_PASSWORD
