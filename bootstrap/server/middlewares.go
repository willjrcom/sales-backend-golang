package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	headerservice "github.com/willjrcom/sales-backend-go/internal/infra/service/header"
	jwtservice "github.com/willjrcom/sales-backend-go/internal/infra/service/jwt"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

func (c *ServerChi) middlewareAuthUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// liberação de preflight
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		// Verificar se a URL atual está na lista de URLs afetados
		shouldValidate := true
		unprotectUserDelete := r.Method == http.MethodDelete && strings.HasSuffix(r.URL.Path, "/user")

		if unprotectUserDelete {
			shouldValidate = false
		}

		// normalizar caminho
		path := strings.TrimSuffix(r.URL.Path, "/")

		// usar prefixo (ou exato) em vez de Contains
		for _, route := range c.UnprotectedRoutes {
			if strings.HasPrefix(path, route) {
				shouldValidate = false
				break
			}
		}

		if shouldValidate {
			// Executar a lógica desejada apenas para os endpoints selecionados
			fmt.Println("Antes de chamar o endpoint:", r.URL.Path)

			validToken, err := headerservice.GetAnyValidToken(ctx, r)

			if err != nil {
				jsonpkg.ResponseErrorJson(w, r, http.StatusUnauthorized, err)
				return
			}

			ctx = context.WithValue(ctx, model.Schema("schema"), jwtservice.GetSchemaFromAccessToken(validToken))
			ctx = context.WithValue(ctx, companyentity.UserValue("user_id"), jwtservice.GetUserIDFromToken(validToken))
		}

		// Chamando o próximo handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
