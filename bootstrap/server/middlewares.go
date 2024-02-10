package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	schemaentity "github.com/willjrcom/sales-backend-go/internal/domain/schema"
	userentity "github.com/willjrcom/sales-backend-go/internal/domain/user"
	headerservice "github.com/willjrcom/sales-backend-go/internal/infra/service/header"
	jwtservice "github.com/willjrcom/sales-backend-go/internal/infra/service/jwt"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Verificar se a URL atual está na lista de URLs afetados
		shouldLog := false
		for _, url := range []string{"/category-product"} {
			if strings.Contains(r.URL.Path, url) {
				shouldLog = true
				break
			}
		}

		if shouldLog {
			// Executar a lógica desejada apenas para os endpoints selecionados
			fmt.Println("Antes de chamar o endpoint:", r.URL.Path)

			tokenString, err := headerservice.GetIDTokenHeader(r)

			if err != nil {
				jsonpkg.ResponseJson(w, r, http.StatusUnauthorized, jsonpkg.Error{Message: err.Error()})
				return
			}

			token, error := jwtservice.ValidateToken(ctx, tokenString)

			if error != nil {
				jsonpkg.ResponseJson(w, r, http.StatusUnauthorized, jsonpkg.Error{Message: error.Error()})
				return
			}

			ctx = context.WithValue(ctx, schemaentity.Schema("schema"), jwtservice.GetSchemaFromToken(token))
			ctx = context.WithValue(ctx, userentity.UserValue("user"), jwtservice.GetUserFromToken(token))
		}

		// Chamando o próximo handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
