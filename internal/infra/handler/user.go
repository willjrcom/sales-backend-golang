package handlerimpl

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	schemaentity "github.com/willjrcom/sales-backend-go/internal/domain/schema"
	userdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/user"
	headerservice "github.com/willjrcom/sales-backend-go/internal/infra/service/header"
	jwtservice "github.com/willjrcom/sales-backend-go/internal/infra/service/jwt"
	userusecases "github.com/willjrcom/sales-backend-go/internal/usecases/user"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type handlerUserImpl struct {
	s *userusecases.Service
}

func NewHandlerUser(userService *userusecases.Service) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerUserImpl{
		s: userService,
	}

	c.With().Group(func(c chi.Router) {
		c.Post("/new", h.handlerNewUser)
		c.Post("/update", h.handlerUpdateUser)
		c.Post("/login", h.handlerLoginUser)
		c.Post("/access", h.handlerAccess)
		c.Delete("/delete", h.handlerDeleteUser)
	})

	return handler.NewHandler("/user", c)
}

func (h *handlerUserImpl) handlerNewUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx = context.WithValue(ctx, schemaentity.Schema("schema"), schemaentity.DEFAULT_SCHEMA)

	user := &userdto.CreateUserInput{}
	jsonpkg.ParseBody(r, user)

	if err := h.s.CreateUser(ctx, user); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
	}
}

func (h *handlerUserImpl) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx = context.WithValue(ctx, schemaentity.Schema("schema"), schemaentity.DEFAULT_SCHEMA)

	user := &userdto.UpdatePasswordInput{}
	jsonpkg.ParseBody(r, user)

	if err := h.s.UpdateUser(ctx, user); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
	}
}

func (h *handlerUserImpl) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx = context.WithValue(ctx, schemaentity.Schema("schema"), schemaentity.DEFAULT_SCHEMA)

	user := &userdto.LoginUserInput{}
	jsonpkg.ParseBody(r, user)

	if token, err := h.s.LoginUser(ctx, user); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: token})
	}
}

func (h *handlerUserImpl) handlerAccess(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	headerToken, _ := headerservice.GetAccessTokenHeader(r)
	accessToken, err := jwtservice.ValidateToken(ctx, headerToken)

	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusUnauthorized, jsonpkg.Error{Message: err.Error()})
		return
	}

	schema := &userdto.AccessCompanyInput{}
	jsonpkg.ParseBody(r, schema)

	if token, err := h.s.Access(ctx, schema, accessToken); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: token})
	}
}

func (h *handlerUserImpl) handlerDeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx = context.WithValue(ctx, schemaentity.Schema("schema"), schemaentity.DEFAULT_SCHEMA)

	user := &userdto.DeleteUserInput{}
	jsonpkg.ParseBody(r, user)

	if err := h.s.DeleteUser(ctx, user); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
	}
}
