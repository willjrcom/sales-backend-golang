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

		// Verificar se a URL atual está na lista de URLs afetados
		shouldValidate := true
		unprotectUserDelete := r.Method == http.MethodDelete && strings.HasSuffix(r.URL.Path, "/user")

		if unprotectUserDelete {
			shouldValidate = false
		}

		if shouldValidate {
			for _, url := range c.UnprotectedRoutes {

				if strings.Contains(r.URL.Path, url) || unprotectUserDelete {
					shouldValidate = false
					break
				}
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

			ctx = context.WithValue(ctx, model.Schema("schema"), jwtservice.GetSchemaFromIDToken(validToken))
			ctx = context.WithValue(ctx, companyentity.UserValue("user_id"), jwtservice.GetUserIDFromToken(validToken))
		}

		// Chamando o próximo handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
