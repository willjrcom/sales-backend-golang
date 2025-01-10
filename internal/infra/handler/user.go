package handlerimpl

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
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

	route := "/user"

	c.With().Group(func(c chi.Router) {
		c.Post("/new", h.handlerNewUser)
		c.Patch("/update/password", h.handlerUpdateUserPassword)
		c.Post("/forget-password", h.handlerForgetUserPassword)
		c.Patch("/update/{id}", h.handlerUpdateUser)
		c.Post("/login", h.handlerLoginUser)
		c.Post("/access", h.handlerAccess)
		c.Post("/search", h.handlerSearchUser)
		c.Delete("/", h.handlerDeleteUser)
	})

	unprotectedRoutes := []string{
		fmt.Sprintf("%s/new", route),
		fmt.Sprintf("%s/login", route),
		fmt.Sprintf("%s/access", route),
		fmt.Sprintf("%s/update-password", route),
		fmt.Sprintf("%s/forget-password", route),
	}

	return handler.NewHandler(route, c, unprotectedRoutes...)
}

func (h *handlerUserImpl) handlerNewUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoUser := &userdto.UserCreateDTO{}
	if err := jsonpkg.ParseBody(r, dtoUser); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: err.Error()})
		return
	}

	id, err := h.s.CreateUser(ctx, dtoUser)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: id})
}

func (h *handlerUserImpl) handlerUpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoUser := &userdto.UserUpdatePasswordDTO{}
	if err := jsonpkg.ParseBody(r, dtoUser); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: err.Error()})
		return
	}

	if err := h.s.UpdateUserPassword(ctx, dtoUser); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerUserImpl) handlerForgetUserPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoUser := &userdto.UserForgetPasswordDTO{}
	if err := jsonpkg.ParseBody(r, dtoUser); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: err.Error()})
		return
	}

	if err := h.s.ForgetUserPassword(ctx, dtoUser); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerUserImpl) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoID := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	dtoUser := &userdto.UserUpdateDTO{}
	if err := jsonpkg.ParseBody(r, dtoUser); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: err.Error()})
		return
	}

	if err := h.s.UpdateUser(ctx, dtoID, dtoUser); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerUserImpl) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoUser := &userdto.UserLoginDTO{}
	if err := jsonpkg.ParseBody(r, dtoUser); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: err.Error()})
		return
	}

	token, err := h.s.LoginUser(ctx, dtoUser)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: token})
}

func (h *handlerUserImpl) handlerAccess(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	headerToken, _ := headerservice.GetAccessTokenHeader(r)
	accessToken, err := jwtservice.ValidateToken(ctx, headerToken)

	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusUnauthorized, jsonpkg.Error{Message: err.Error()})
		return
	}

	dtoSchema := &userdto.UserSchemaDTO{}
	if err := jsonpkg.ParseBody(r, dtoSchema); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: err.Error()})
		return
	}

	token, err := h.s.Access(ctx, dtoSchema, accessToken)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: token})
}

func (h *handlerUserImpl) handlerSearchUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoUser := &userdto.UserSearchDTO{}
	if err := jsonpkg.ParseBody(r, dtoUser); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: err.Error()})
		return
	}

	user, err := h.s.SearchUser(ctx, dtoUser)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: user})
}

func (h *handlerUserImpl) handlerDeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoUser := &userdto.UserDeleteDTO{}
	if err := jsonpkg.ParseBody(r, dtoUser); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: err.Error()})
		return
	}

	if err := h.s.DeleteUser(ctx, dtoUser); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}
