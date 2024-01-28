package schemaservice

import "net/http"

func GetSchemaHeader(r *http.Request) string {
	return r.Header.Get("schema")
}
