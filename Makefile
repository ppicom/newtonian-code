.PHONY: proto clean deps

# Variables
PROTO_DIR = internal/infrastructure/api/grpc
PROTO_FILE = $(PROTO_DIR)/banking_v1.proto
PROTO_OUT_DIR = $(PROTO_DIR)/banking/v1

# Install protobuf generation dependencies
deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Clean generated files
clean:
	find . -name "*.pb.go" -delete

# Generate protocol buffer code
proto: clean deps
	protoc \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_FILE)

# Build the project
build: proto
	go build -o bin/server cmd/main.go

# Run the server
run: build
	./bin/server 