BINDIR=${CURDIR}/bin
PACKAGE=example/external

bindir:
	mkdir -p ${BINDIR}

build: bindir
	echo "build external"

PROTOC = PATH="$$PATH:$(BINDIR)" protoc

vendor-rm:
	rm -rf vendor-proto

vendor-proto/google/protobuf:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/protocolbuffers/protobuf vendor-proto/protobuf &&\
	cd vendor-proto/protobuf &&\
	git sparse-checkout set --no-cone src/google/protobuf &&\
	git checkout
	mkdir -p vendor-proto/google
	mv vendor-proto/protobuf/src/google/protobuf vendor-proto/google
	rm -rf vendor-proto/protobuf


vendor-proto/validate:
	git clone -b main --single-branch --depth=2 --filter=tree:0 \
		https://github.com/bufbuild/protoc-gen-validate vendor-proto/tmp && \
		cd vendor-proto/tmp && \
		git sparse-checkout set --no-cone validate &&\
		git checkout
		mkdir -p vendor-proto/validate
		mv vendor-proto/tmp/validate vendor-proto/
		rm -rf vendor-proto/tmp


vendor-proto: vendor-rm  vendor-proto/google/protobuf vendor-proto/validate

bin-deps:
	$(info Installing binary dependencies...)

	GOBIN=$(BINDIR) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1 && \
    GOBIN=$(BINDIR) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0 && \
	GOBIN=$(BINDIR) go install github.com/bufbuild/buf/cmd/buf@v1.21.0 && \
	GOBIN=$(BINDIR) go install github.com/envoyproxy/protoc-gen-validate@v1.0.4 && \

PROTO_PATH:="api/v1"

protoc-generate:
	protoc \
	-I ${PROTO_PATH} \
	-I vendor-proto \
	--plugin=protoc-gen-go=$(BINDIR)/protoc-gen-go \
	--go_out pkg/${PROTO_PATH} \
	--go_opt paths=source_relative \
	--plugin=protoc-gen-go-grpc=$(BINDIR)/protoc-gen-go-grpc \
	--go-grpc_out pkg/${PROTO_PATH} \
	--go-grpc_opt paths=source_relative \
	--plugin=protoc-gen-validate=$(BINDIR)/protoc-gen-validate \
	--validate_out="lang=go,paths=source_relative:pkg/api/v1" \
	api/v1/external.proto
	go mod tidy

