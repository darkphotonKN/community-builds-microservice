#!/bin/bash

# Script to generate protobuf files with googleapis support
# Usage: ./generate_proto.sh

set -e

echo "Generating protobuf files..."

# Make sure we're in the common directory
cd "$(dirname "$0")"

# Generate each proto file with googleapis support
protoc \
  --go_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_out=. \
  --go-grpc_opt=paths=source_relative \
  --proto_path=. \
  --proto_path=third_party/googleapis \
  api/proto/notification/notification.proto

echo "✅ notification.proto generated"

# Add other proto files here as needed
# protoc \
#   --go_out=. \
#   --go_opt=paths=source_relative \
#   --go-grpc_out=. \
#   --go-grpc_opt=paths=source_relative \
#   --proto_path=. \
#   --proto_path=third_party/googleapis \
#   api/proto/other/other.proto

echo "🎉 All protobuf files generated successfully!"
