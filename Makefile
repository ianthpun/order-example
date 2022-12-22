.PHONY: install-proto-tools
install-proto-tools:
	GOFLAGS="" go install github.com/golang/protobuf/protoc-gen-go@v1.5.2
	GOFLAGS="" go install github.com/bufbuild/buf/cmd/buf@v1.4.0
	GOFLAGS="" go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0

.PHONY: generate-proto
generate-proto: install-proto-tools
	cd ./api/proto \
	 && buf generate
