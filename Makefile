# ---- config ---------------------------------------------------------------

PROTO_DIRS  := alkibiades antisthenes aristippos aspasia diotima kritias kriton xenofon

# tools / images
PROTOC          ?= protoc
BUF             ?= buf
DOCKER          ?= docker
PROTO_DOC_IMAGE ?= localproto:latest

SPECTAQL_DIR    := ./sokrates/docs
SPECTAQL_CFG    := spectaql.yaml

ROOT ?= .
OWNER ?= odysseia-greek

# Make runs each recipe line in its own shell by default; enable bash + strict mode.
SHELL := /usr/bin/env bash
.SHELLFLAGS := -euo pipefail -c

.DEFAULT_GOAL := all

# ---- helper "functions" ---------------------------------------------------

# Run a bash snippet for each directory in a list:
# $(call for_each_dir,<dir list>,<bash body using $$dir>)
define for_each_dir
@for dir in $(1); do \
	echo "==> $$dir"; \
	$(2); \
done
endef

# Does $$dir contain a file? (used in loops)
define if_file
if [[ -f "$(1)" ]]; then $(2); fi
endef

# ---- targets --------------------------------------------------------------

.PHONY: all generate docs tidy tidy-go mods

all: generate docs

generate: generate-buf

# Run a bash snippet for each directory in a list:
define for_each_dir
@for dir in $(1); do \
	echo "==> $$dir"; \
	$(2) \
done
endef

.PHONY: generate-buf
generate-buf:
	$(call for_each_dir,$(PROTO_DIRS), \
		echo "Generating Protobuf (buf) in $$dir..."; \
		$(BUF) generate --template "$$dir/buf.gen.yaml" "$$dir"; \
	)
.PHONY: docs
docs: docs-grpc spectaql

.PHONY: tools
tools: $(PROTOC_GEN_DOC)

$(PROTOC_GEN_DOC):
	@mkdir -p $(TOOLS_DIR)
	@echo "Installing protoc-gen-doc..."
	@GOBIN=$(TOOLS_DIR) go install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@latest


.PHONY: docs-grpc
docs-grpc: tools
	@for dir in $(PROTO_DIRS); do \
		echo "Generating gRPC docs in $$dir..."; \
		PATH=$(TOOLS_DIR):$$PATH buf generate --template $$dir/buf.gen.docs.yaml $$dir; \
	done

.PHONY: spectaql
spectaql:
	@echo "==> Generating SpectaQL docs in $(SPECTAQL_DIR)..."
	@cd "$(SPECTAQL_DIR)" && spectaql -c "$(SPECTAQL_CFG)"

# --------------------------------------------------------------------------
# Go module tidying
# --------------------------------------------------------------------------
# Runs `go mod tidy` in *every* directory that has a go.mod (recursively).
# Uses `find` so it doesn't need you to maintain a list.

.PHONY: tidy tidy-go
tidy: tidy-go

tidy-go:
	@echo "==> Running 'go mod tidy' in all modules..."
	@mods=$$(find . -name go.mod -print0 | xargs -0 -n1 dirname | sort -u); \
	if [[ -z "$$mods" ]]; then \
		echo "No go.mod files found."; \
		exit 0; \
	fi; \
	for d in $$mods; do \
		echo "==> go mod tidy in $$d"; \
		( cd "$$d" && go mod tidy && go fmt ./... ); \
	done

# Convenience: list all module dirs found
.PHONY: mods
mods:
	@find . -name go.mod -print0 | xargs -0 -n1 dirname | sort -u


.PHONY: images-dev images-prod

images-dev:
	OWNER="$(OWNER)" ROOT="$(ROOT)" GROUP=dev \
		./bump-images.sh

images-prod:
	OWNER="$(OWNER)" ROOT="$(ROOT)" GROUP=prod \
		./bump-images.sh