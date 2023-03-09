# Set these to the desired values
ARTIFACT_ID=gomarkdoc
VERSION=0.4.1-1

GO_BUILD_FLAGS?=-mod=vendor -a -tags netgo $(LDFLAGS) -installsuffix cgo -o $(BINARY) ./cmd/gomarkdoc
MAKEFILES_VERSION=7.5.0

.DEFAULT_GOAL:=default

include build/make/build.mk
include build/make/self-update.mk
include build/make/dependencies-gomod.mk
include build/make/variables.mk
include build/make/clean.mk
include build/make/digital-signature.mk

default: compile-generic

PHONY: generate
generate:
	go generate .