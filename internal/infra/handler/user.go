package handlerimpl

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	companydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/company"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
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
		c.Patch("/update/forget-password", h.handlerUpdateUserForgetPassword)
		c.Post("/forget-password", h.handlerForgetUserPassword)
		c.Patch("/update/{id}", h.handlerUpdateUser)
		c.Post("/login", h.handlerLoginUser)
		c.Post("/access", h.handlerAccess)
		c.Post("/search", h.handlerSearchUser)
		c.Delete("/", h.handlerDeleteUser)
		c.Get("/refresh-access-token", h.handlerRefreshAccessToken)
		c.Get("/companies", h.handlerGetCompanies)
		c.Post("/validate-password-reset-token", h.handlerValidatePasswordResetToken)
	})

	unprotectedRoutes := []string{
		fmt.Sprintf("%s/new", route),
		fmt.Sprintf("%s/login", route),
		fmt.Sprintf("%s/update/forget-password", route),
		fmt.Sprintf("%s/forget-password", route),
		fmt.Sprintf("%s/validate-password-reset-token", route),
	}

	return handler.NewHandler(route, c, unprotectedRoutes...)
}

func (h *handlerUserImpl) handlerNewUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoUser := &companydto.UserCreateDTO{}
	if err := jsonpkg.ParseBody(r, dtoUser); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	id, err := h.s.CreateUser(ctx, dtoUser)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, id)
}

func (h *handlerUserImpl) handlerUpdateUserForgetPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoUser := &companydto.UserUpdateForgetPasswordDTO{}
	if err := jsonpkg.ParseBody(r, dtoUser); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	email, err := jwtservice.ValidatePasswordResetToken(ctx, dtoUser.Token)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	if email != dtoUser.Email {
		err := fmt.Errorf("email do token invalido!")
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	if err := h.s.UpdateUserForgetPassword(ctx, dtoUser); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerUserImpl) handlerUpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoUser := &companydto.UserUpdatePasswordDTO{}
	if err := jsonpkg.ParseBody(r, dtoUser); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.s.UpdateUserPassword(ctx, dtoUser); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerUserImpl) handlerForgetUserPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoUser := &companydto.UserForgetPasswordDTO{}
	if err := jsonpkg.ParseBody(r, dtoUser); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.s.ForgetUserPassword(ctx, dtoUser); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerUserImpl) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoID := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	dtoUser := &companydto.UserUpdateDTO{}
	if err := jsonpkg.ParseBody(r, dtoUser); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.s.UpdateUser(ctx, dtoID, dtoUser); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerUserImpl) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoUser := &companydto.UserLoginDTO{}
	if err := jsonpkg.ParseBody(r, dtoUser); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	user, err := h.s.LoginUser(ctx, dtoUser)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, user)
}

func (h *handlerUserImpl) handlerAccess(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idToken, err := headerservice.GetIDTokenFromHeader(r)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusUnauthorized, err)
		return
	}

	validIdToken, err := jwtservice.ValidateToken(ctx, idToken)

	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusUnauthorized, err)
		return
	}

	dtoSchema := &companydto.UserSchemaDTO{}
	if err := jsonpkg.ParseBody(r, dtoSchema); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	accessToken, err := h.s.Access(ctx, dtoSchema, validIdToken)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, accessToken)
}

func (h *handlerUserImpl) handlerSearchUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoUser := &companydto.UserSearchDTO{}
	if err := jsonpkg.ParseBody(r, dtoUser); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	user, err := h.s.SearchUser(ctx, dtoUser)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, user)
}

func (h *handlerUserImpl) handlerDeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoUser := &companydto.UserDeleteDTO{}
	if err := jsonpkg.ParseBody(r, dtoUser); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.s.DeleteUser(ctx, dtoUser); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerUserImpl) handlerRefreshAccessToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// require an existing ID token to refresh
	accessToken, err := headerservice.GetAccessTokenFromHeader(r)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusUnauthorized, err)
		return
	}

	validAccessToken, err := jwtservice.ValidateToken(ctx, accessToken)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusUnauthorized, err)
		return
	}

	schema := jwtservice.GetSchemaFromAccessToken(validAccessToken)
	newAccessToken, err := jwtservice.CreateAccessToken(validAccessToken, schema)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, newAccessToken)
}

func (h *handlerUserImpl) handlerGetCompanies(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	companies, err := h.s.GetCompanies(ctx)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, companies)
}

func (h *handlerUserImpl) handlerValidatePasswordResetToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := &companydto.UserResetTokenRequestDTO{}
	if err := jsonpkg.ParseBody(r, &req); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	email, err := jwtservice.ValidatePasswordResetToken(ctx, req.Token)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, companydto.UserResetTokenResponseDTO{
		Valid: true,
		Email: email,
	})
}
