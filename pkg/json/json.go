package jsonpkg

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type HTTPResponse struct {
	Success  bool        `json:"success"`
	Data     interface{} `json:"data,omitempty"`
	Error    *Error      `json:"error,omitempty"`
	Metadata *Metadata   `json:"metadata,omitempty"`
}

type Error struct {
	Message string `json:"message,omitempty"`
}

type Metadata struct {
	Page         int `json:"page"`
	PageSize     int `json:"page_size"`
	TotalPages   int `json:"total_pages"`
	TotalResults int `json:"total_results"`
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
	w.WriteHeader(statusCode)

	response := &HTTPResponse{
		Success: true,
		Data:    data,
	}

	// Serializa a resposta para JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		// Erro na serialização do JSON
		http.Error(w, "Internal Server Error: failed to serialize response", http.StatusInternalServerError)
		return
	}

	// Escreve o JSON no body da resposta
	_, err = w.Write(jsonResponse)
	if err != nil {
		// Loga o erro de escrita no body (não é possível alterar a resposta após WriteHeader)
		w.Write([]byte("Failed to write response: " + err.Error()))
	}
}

func ResponseErrorJson(w http.ResponseWriter, r *http.Request, statusCode int, err error) {
	fmt.Printf("ResponseErrorJson %v\n", err.Error())

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	response := &HTTPResponse{
		Success: true,
		Error: &Error{
			Message: err.Error(),
		},
	}

	// Serializa a resposta para JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		// Erro na serialização do JSON
		http.Error(w, "Internal Server Error: failed to serialize response", http.StatusInternalServerError)
		return
	}

	// Escreve o JSON no body da resposta
	if _, err = w.Write(jsonResponse); err != nil {
		// Loga o erro de escrita no body (não é possível alterar a resposta após WriteHeader)
		w.Write([]byte("Failed to write response: " + err.Error()))
	}
}
