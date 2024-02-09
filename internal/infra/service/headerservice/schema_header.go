package headerservice

import "net/http"

func GetSchemaHeader(r *http.Request) string {
	return r.Header.Get("schema")
}

func GetAccessTokenHeader(r *http.Request) string {
	return r.Header.Get("access-token")
}
