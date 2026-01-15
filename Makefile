# ---- config ---------------------------------------------------------------

PROTO_DIRS  := alkibiades antisthenes aristippos aspasia diotima kritias kriton xenofon

# tools / images
PROTOC          ?= protoc
BUF             ?= buf
DOCKER          ?= docker
PROTO_DOC_IMAGE ?= localproto:latest

SPECTAQL_DIR    := ./sokrates/docs
SPECTAQL_CFG    := spectaql.yaml

# Make runs each recipe line in its own shell by default; enable bash + strict mode.
SHELL := bash
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

SHELL := bash
.SHELLFLAGS := -euo pipefail -c

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
docs: docs-proto spectaql

.PHONY: docs-proto
docs-proto:
	$(call for_each_dir,$(PROTO_DIRS), \
		echo "Generating docs in $$dir..."; \
		$(DOCKER) run --rm \
			-v "$$PWD/$$dir/docs:/out" \
			-v "$$PWD/$$dir/proto:/protos" \
			$(PROTO_DOC_IMAGE) --doc_opt=html,docs.html; \
		$(DOCKER) run --rm \
			-v "$$PWD/$$dir/docs:/out" \
			-v "$$PWD/$$dir/proto:/protos" \
			$(PROTO_DOC_IMAGE) --doc_opt=markdown,docs.md; \
	)

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
		( cd "$$d" && go mod tidy ); \
	done

# Convenience: list all module dirs found
.PHONY: mods
mods:
	@find . -name go.mod -print0 | xargs -0 -n1 dirname | sort -u