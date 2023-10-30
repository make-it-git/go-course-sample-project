all:
	mkdir -p internal/generated/schema
	oapi-codegen -package rider -generate chi-server,types,spec api/schema.yaml > internal/generated/schema/rider.gen.go