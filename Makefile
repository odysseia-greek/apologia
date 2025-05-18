PROTO_DIRS := alkibiades antisthenes aristippos kritias kriton xenofon

.PHONY: all
all: generate docs

.PHONY: generate
generate:
	@for dir in $(PROTO_DIRS); do \
		echo "Generating Protobuf files in $$dir..."; \
		(cd $$dir && \
		 protoc --go_out=. --go_opt=paths=source_relative \
		        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
		        proto/$$dir.proto); \
	done

.PHONY: docs
docs:
	@for dir in $(PROTO_DIRS); do \
		echo "Generating docs in $$dir..."; \
		docker run --rm \
			-v $$PWD/$$dir/docs:/out \
			-v $$PWD/$$dir/proto:/protos \
			localproto:latest --doc_opt=html,docs.html; \
		docker run --rm \
			-v $$PWD/$$dir/docs:/out \
			-v $$PWD/$$dir/proto:/protos \
			localproto:latest --doc_opt=markdown,docs.md; \
	done
	cd ./sokrates/docs && spectaql -c  spectaql.yaml

