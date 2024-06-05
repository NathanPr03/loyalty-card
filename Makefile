OAPI_PACKAGE := generated
OAPI_DIRECTORY := api/${OAPI_PACKAGE}
OAPI_SPEC_PATH := openapi.yml

.PHONY: codegen
codegen:
	oapi-codegen -generate types -package ${OAPI_PACKAGE} ${OAPI_SPEC_PATH} > ${OAPI_DIRECTORY}/models.gen.go