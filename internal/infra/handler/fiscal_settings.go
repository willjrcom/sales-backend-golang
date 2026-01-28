package handlerimpl

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	fiscalsettingsdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/fiscal_settings"
	fiscalsettingsusecases "github.com/willjrcom/sales-backend-go/internal/usecases/fiscal_settings"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type handlerFiscalSettingsImpl struct {
	s *fiscalsettingsusecases.Service
}

func NewFiscalSettingsHandler(service *fiscalsettingsusecases.Service) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerFiscalSettingsImpl{
		s: service,
	}

	c.With().Group(func(c chi.Router) {
		c.Get("/", h.handlerGetFiscalSettings)
		c.Patch("/", h.handlerUpdateFiscalSettings)
	})

	return handler.NewHandler("/company/fiscal-settings", c)
}

func (h *handlerFiscalSettingsImpl) handlerGetFiscalSettings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dto, err := h.s.GetFiscalSettings(ctx)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, dto)
}

func (h *handlerFiscalSettingsImpl) handlerUpdateFiscalSettings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dto := &fiscalsettingsdto.FiscalSettingsUpdateDTO{}
	if err := jsonpkg.ParseBody(r, dto); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.s.UpdateFiscalSettings(ctx, dto); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}
