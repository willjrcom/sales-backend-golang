package server

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	headerservice "github.com/willjrcom/sales-backend-go/internal/infra/service/header"
	jwtservice "github.com/willjrcom/sales-backend-go/internal/infra/service/jwt"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

// middlewareRecover captura panics em handlers e evita que a aplicação caia.
// Ele registra a pilha e retorna 500 com uma mensagem genérica.
func (c *ServerChi) middlewareRecover(next http.Handler) http.Handler {
	fmt.Println("middlewareRecover init")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				// Logar o panic e stack trace para diagnóstico
				fmt.Printf("panic recovered: %v\n%s", rec, debug.Stack())

				// Responder com erro genérico (não expor detalhes do panic ao cliente)
				jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, errors.New("internal server error"))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (c *ServerChi) middlewareAuthUser(next http.Handler) http.Handler {
	fmt.Println("middlewareAuthUser init")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Executar a lógica desejada apenas para os endpoints selecionados
		fmt.Println("Antes de chamar o endpoint:", r.URL.Path)

		// liberação de preflight
		if r.Method == http.MethodOptions {
			fmt.Println("middlewareAuthUser options")
			next.ServeHTTP(w, r)
			return
		}

		// Verificar se a URL atual está na lista de URLs afetados
		shouldValidate := true
		unprotectUserDelete := r.Method == http.MethodDelete && strings.HasSuffix(r.URL.Path, "/user")

		if unprotectUserDelete {
			fmt.Println("middlewareAuthUser unprotectUserDelete")
			shouldValidate = false
		}

		// normalizar caminho
		path := strings.TrimSuffix(r.URL.Path, "/")

		// usar prefixo (ou exato) em vez de Contains
		for _, route := range c.UnprotectedRoutes {
			if strings.HasPrefix(path, route) {
				shouldValidate = false
				fmt.Println("middlewareAuthUser unprotected route", route)
				break
			}
		}

		validToken, _ := headerservice.GetAccessTokenFromHeader(ctx, r)
		if shouldValidate && validToken != nil && validToken.Valid {
			schema := jwtservice.GetSchemaFromAccessToken(validToken)
			userID := jwtservice.GetUserIDFromToken(validToken)
			ctx = context.WithValue(ctx, model.Schema("schema"), schema)
			ctx = context.WithValue(ctx, companyentity.UserValue("user_id"), userID)

			// Chamando o próximo handler
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		appToken := r.Header.Get("Authorization")
		if shouldValidate && appToken != "" {
			// Remover o prefixo "Bearer " se existir
			appToken = strings.TrimPrefix(appToken, "Bearer ")

			// Tentar decodificar Base64
			schemaBytes, err := base64.StdEncoding.DecodeString(appToken)
			if err != nil {
				jsonpkg.ResponseErrorJson(w, r, http.StatusUnauthorized, err)
				return
			}

			schema := string(schemaBytes)
			ctx = context.WithValue(ctx, model.Schema("schema"), schema)
			// Chamando o próximo handler com o contexto atualizado
			next.ServeHTTP(w, r.WithContext(ctx))
			return

		}

		jsonpkg.ResponseErrorJson(w, r, http.StatusUnauthorized, errors.New("unauthorized"))
	})
}
