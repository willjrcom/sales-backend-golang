package headerservice

import (
	"net/http"
	"strconv"
)

func GetPageAndPerPage(r *http.Request, defaultPage int, defaultPerPage int) (int, int) {
	page := defaultPage
	if p := r.URL.Query().Get("page"); p != "" {
		if iv, err := strconv.Atoi(p); err == nil && iv > 0 {
			page = iv
		}
	}

	perPage := defaultPerPage
	if pp := r.URL.Query().Get("per_page"); pp != "" {
		if ipv, err := strconv.Atoi(pp); err == nil && ipv > 0 {
			perPage = ipv
		}
	}

	return page, perPage
}
