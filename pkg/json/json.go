package jsonpkg

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

func ParseBody(r *http.Request, output interface{}) error {
	// Get the contentType for comparisons
	ct := r.Header.Get("Content-Type")
	if r.Method == http.MethodGet || r.Method == http.MethodHead {
		return nil
	}

	if strings.Contains(ct, "application/json") {
		body, err := io.ReadAll(r.Body)
		r.Body.Close()
		r.Body = io.NopCloser(bytes.NewBuffer(body))
		if err != nil {
			return err
		}

		err = json.Unmarshal(body, &output)
		if err != nil {
			return err
		}

		return nil
	}

	return errors.New("header content-type is different from application/json")
}
