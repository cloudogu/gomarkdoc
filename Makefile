# Set these to the desired values
ARTIFACT_ID=gomarkdoc
VERSION=0.4.1-1

GOTAG?=1.19.7
MAKEFILES_VERSION=7.5.0

GO_BUILD_FLAGS?=-mod=vendor -a -tags netgo $(LDFLAGS) -installsuffix cgo -o $(BINARY) ./cmd/gomarkdoc
.DEFAULT_GOAL:=default

include build/make/variables.mk
include build/make/build.mk
include build/make/self-update.mk
include build/make/dependencies-gomod.mk
include build/make/clean.mk


default: compile
