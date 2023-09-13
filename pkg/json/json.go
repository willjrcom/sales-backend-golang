package jsonpkg

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

type HTTPResponse struct {
	Data interface{} `json:"data,omitempty"`
}

type Error struct {
	Message string `json:"message,omitempty"`
}

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

func ResponseJson(w http.ResponseWriter, r *http.Request, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Write status
	w.WriteHeader(statusCode)

	if statusCode == http.StatusNoContent || data == nil {
		noContent, _ := json.Marshal(struct{}{})
		w.Write([]byte(noContent))
		return
	}

	if jsonResponse, err := json.Marshal(data); err != nil {
		w.Write([]byte("Internal Server Error: " + err.Error()))
		return

	} else if _, err = w.Write(jsonResponse); err != nil {
		w.Write([]byte("Internal Server Error: " + err.Error()))
		return
	}
}
