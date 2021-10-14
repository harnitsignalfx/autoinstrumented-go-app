# autoinstrumented-go-app
AutoInstrumented (sort of, with the otel libs) Golang App


The env variables to run this script are mentioned below

OTEL_RESOURCE_ATTRIBUTES=service.name=my-go-client-app,deployment.environment=istio OTEL_EXPORTER_JAEGER_ENDPOINT=https://ingest.<realm>.signalfx.com/v2/trace OTEL_PROPAGATORS=b3,b3multi SPLUNK_ACCESS_TOKEN=<token> go run istio-sender.go
