BINDIR=${CURDIR}/bin
PACKAGE=example/comments

bindir:
	mkdir -p ${BINDIR}

build: bindir
	echo "build comments"

run-migrations:
	echo "run migrations"

coverage:
	go test ./... -coverprofile cover.out.tmp
	cat cover.out.tmp | grep -v "_mock.go" | grep -v "\.pb\." > cover.out
	rm cover.out.tmp
	go tool cover -func cover.out

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


# Устанавливаем proto описания google/googleapis
vendor-proto/google/api:
	git clone -b master --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/googleapis/googleapis vendor-proto/googleapis && \
 	cd vendor-proto/googleapis && \
	git sparse-checkout set --no-cone google/api && \
	git checkout
	mkdir -p  vendor-proto/google
	mv vendor-proto/googleapis/google/api vendor-proto/google
	rm -rf vendor-proto/googleapis

vendor-proto/protoc-gen-openapiv2/options:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/grpc-ecosystem/grpc-gateway vendor-proto/grpc-ecosystem && \
 	cd vendor-proto/grpc-ecosystem && \
	git sparse-checkout set --no-cone protoc-gen-openapiv2/options && \
	git checkout
	mkdir -p vendor-proto/protoc-gen-openapiv2
	mv vendor-proto/grpc-ecosystem/protoc-gen-openapiv2/options vendor-proto/protoc-gen-openapiv2
	rm -rf vendor-proto/grpc-ecosystem


vendor-proto: vendor-rm  vendor-proto/google/protobuf vendor-proto/validate vendor-proto/google/api vendor-proto/protoc-gen-openapiv2/options

bin-deps:
	$(info Installing binary dependencies...)

	GOBIN=$(BINDIR) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1 && \
    GOBIN=$(BINDIR) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0 && \
	GOBIN=$(BINDIR) go install github.com/bufbuild/buf/cmd/buf@v1.21.0 && \
	GOBIN=$(BINDIR) go install github.com/envoyproxy/protoc-gen-validate@v1.0.4 && \
	GOBIN=$(BINDIR) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.19.1 && \
	GOBIN=$(BINDIR) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.19.1 && \
	GOBIN=$(BINDIR) go install github.com/go-swagger/go-swagger/cmd/swagger@v0.30.5




PROTO_PATH:="api/comments/v1"

protoc-generate:
	mkdir -p api/openapiv2
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
	--validate_out="lang=go,paths=source_relative:pkg/api/comments/v1" \
	--plugin=protoc-gen-grpc-gateway=$(BINDIR)/protoc-gen-grpc-gateway \
	--grpc-gateway_out pkg/${PROTO_PATH} \
	--grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true \
	--plugin=protoc-gen-openapiv2=$(BINDIR)/protoc-gen-openapiv2 \
	--openapiv2_out api/openapiv2 \
	--openapiv2_opt logtostderr=true \
	api/comments/v1/comments.proto
	go mod tidy

serve-swagger:
	bin/swagger serve api/openapiv2/comments.swagger.json

install-goose:
	$(info Installing goose binary into [$(BINDIR)]...)
	GOBIN=$(BINDIR) go install github.com/pressly/goose/v3/cmd/goose@v3.24.1

install-sqlc:
	$(info Installing sqlc binary into [$(BINDIR)]...)
	GOBIN=$(BINDIR) go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.28.0

up-migrations:
	$(BINDIR)/goose -dir migrations postgres "postgresql://comments-user-1:comments-password-1@192.168.56.1:5432/comments-db?sslmode=disable" up

compile-sql:
	$(BINDIR)/sqlc generate

