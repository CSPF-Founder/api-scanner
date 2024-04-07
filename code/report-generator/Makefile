#!/usr/bin/make
SHELL = /bin/sh

IMAGE_NAME = "projects/api-security/scanner:latest"

CURRENT_UID := $(shell id -u)
CURRENT_GID := $(shell id -g)
CURRENT_USERNAME = $(shell whoami)

test:
	python -m pytest ./tests/

test_all:
	INCLUDE_SKIPPED=1 python -m pytest ./tests/

build:
	docker build -t $(IMAGE_NAME) --build-arg USER=$(CURRENT_USERNAME) .


deps:
	# Create requirements.txt file from poetry.lock
	poetry export -f requirements.txt --output requirements.txt --without-hashes 