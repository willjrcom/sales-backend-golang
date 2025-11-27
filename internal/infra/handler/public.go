package handlerimpl

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	headerservice "github.com/willjrcom/sales-backend-go/internal/infra/service/header"
	jwtservice "github.com/willjrcom/sales-backend-go/internal/infra/service/jwt"
	companyusecases "github.com/willjrcom/sales-backend-go/internal/usecases/company"
	userusecases "github.com/willjrcom/sales-backend-go/internal/usecases/user"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type handlerPublicData struct {
	companyService *companyusecases.Service
	userService    *userusecases.Service
}

func NewHandlerPublicData(companyService *companyusecases.Service, userService *userusecases.Service) *handler.Handler {
	c := chi.NewRouter()
	h := &handlerPublicData{
		companyService: companyService,
		userService:    userService,
	}

	c.Get("/companies", h.handlerGetCompanies)
	c.Get("/users", h.handlerGetUsers)

	return handler.NewHandler("/public", c)
}

func (h *handlerPublicData) handlerGetCompanies(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if err := ensureIDToken(ctx, r); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusUnauthorized, err)
		return
	}

	companies, err := h.companyService.ListPublicCompanies(ctx)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, companies)
}

func (h *handlerPublicData) handlerGetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if err := ensureIDToken(ctx, r); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusUnauthorized, err)
		return
	}

	users, err := h.userService.ListPublicUsers(ctx)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, users)
}

func ensureIDToken(ctx context.Context, r *http.Request) error {
	idToken, err := headerservice.GetIDTokenFromHeader(r)
	if err != nil {
		return err
	}

	token, err := jwtservice.ValidateToken(ctx, idToken)
	if err != nil {
		return err
	}

	if token == nil || !token.Valid {
		return errors.New("invalid id-token")
	}

	return nil
}
