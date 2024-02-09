package headerservice

import "net/http"

func GetAccessTokenHeader(r *http.Request) string {
	return r.Header.Get("access-token")
}

func GetIDTokenHeader(r *http.Request) string {
	return r.Header.Get("id-token")
}
