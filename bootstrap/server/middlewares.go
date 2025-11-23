package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	jwtpkg "github.com/dgrijalva/jwt-go"
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

		if shouldValidate {
			// Executar a lógica desejada apenas para os endpoints selecionados
			fmt.Println("Antes de chamar o endpoint:", r.URL.Path)

			// Criar contexto com timeout de 5 segundos para validação do token
			ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			// Canal para receber resultado da validação
			type tokenResult struct {
				token *jwtpkg.Token
				err   error
			}
			resultChan := make(chan tokenResult, 1)

			// Executar validação em goroutine
			go func() {
				fmt.Println("Validando token para:", r.URL.Path)
				token, err := headerservice.GetAnyValidToken(ctxWithTimeout, r)
				resultChan <- tokenResult{token: token, err: err}
			}()

			// Aguardar resultado ou timeout
			select {
			case result := <-resultChan:
				if result.err != nil {
					fmt.Println("Erro ao validar token:", result.err.Error())
					jsonpkg.ResponseErrorJson(w, r, http.StatusUnauthorized, result.err)
					return
				}
				fmt.Println("Token validado com sucesso para:", r.URL.Path)
				ctx = context.WithValue(ctx, model.Schema("schema"), jwtservice.GetSchemaFromAccessToken(result.token))
				ctx = context.WithValue(ctx, companyentity.UserValue("user_id"), jwtservice.GetUserIDFromToken(result.token))

			case <-ctxWithTimeout.Done():
				fmt.Println("TIMEOUT ao validar token para:", r.URL.Path)
				jsonpkg.ResponseErrorJson(w, r, http.StatusRequestTimeout, errors.New("timeout ao validar token"))
				return
			}
		}

		body, err := json.Marshal(r.Body)
		if err != nil {
			jsonpkg.ResponseErrorJson(w, r, http.StatusProcessing, err)
			return
		}

		fmt.Printf("body: %v\n", string(body))

		// Chamando o próximo handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
