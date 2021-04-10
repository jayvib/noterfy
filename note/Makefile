
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

define GEN_SWAGGER_DOC_REST
# Use to generate the Swagger documentation for rest API.
#
# Example:
#		make gen-swagger-doc-rest
endef
.PHONY: gen-swagger-doc-rest
gen-swagger-doc-rest:
	swag init -g ./api/v1/transport/rest/handler.go --output ./api/v1/transport/rest/docs