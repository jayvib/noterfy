
define PROTO_HELP_INFO
# Use to generate the protocol buffer files
# under note package.
#
# Example:
#		make proto
endef
.PHONY: proto
proto:
	protoc \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		proto/note.proto