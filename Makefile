PROTO_SOURCES:=$(wildcard proto/*.proto)

.PHONY: proto
proto: $(PROTO_SOURCES)
	export PATH="$(shell go env GOPATH)/bin:$(PATH)" && \
	protoc \
	-I . \
	-I ./proto \
	--go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    --include_imports --include_source_info --descriptor_set_out=service_descriptor.pb \
	$(PROTO_SOURCES)
