
generate-basic-model:
	oapi-codegen -generate types -o ./pkg/model/generated/v1/openapi_types.gen.go -package v1 ./openapi/basic.yaml

generate-basic-server:
	oapi-codegen -generate 'server' -o ./internal/server/openapi_server.gen.go -package server -import-mapping=./basic.yaml:my-app/model/generated/v1  ./openapi/basic.yaml


generate-model:
	oapi-codegen -generate 'types,skip-prune' -o ./pkg/model/generated/v1/openapi_types.gen.go -package v1 ./openapi/types.yaml

generate-server:
	oapi-codegen -generate 'chi-server,types' -o ./internal/oapi-codegen/openapi_server.gen.go -package server ./openapi/basic.yaml

generate-server-config:
	oapi-codegen --config openapi_config.yaml ./openapi/openapi.spec.yaml

generate-types-config:
	oapi-codegen --config openapi_types_config.yaml ./openapi/types.yaml
	#oapi-codegen --config openapi_types_config.yaml ./openapi/types.yaml

generate-oapi-model:
	./bin/openapi-generator-cli generate -i ./openapi/basic.yaml -g go-server -o ./internal/server/openaapigenerator --generate-alias-as-model  --global-property models

generate-oapi-server:
	./bin/openapi-generator-cli generate -i ./openapi/basic.yaml -g go-server -o ./internal/server/openaapigenerator --generate-alias-as-model  --global-property apis,supportingFiles  --package-name test

