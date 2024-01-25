package handlerimpl

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	schemadto "github.com/willjrcom/sales-backend-go/internal/infra/dto/schema"
	schemausecases "github.com/willjrcom/sales-backend-go/internal/usecases/schema"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type handlerSchemaImpl struct {
	s *schemausecases.Service
}

func NewHandlerSchema(schemaService *schemausecases.Service) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerSchemaImpl{
		s: schemaService,
	}

	c.With().Group(func(c chi.Router) {
		c.Post("/new", h.handlerNewSchema)
	})

	return handler.NewHandler("/schema", c)
}

func (h *handlerSchemaImpl) handlerNewSchema(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	schema := &schemadto.SchemaInput{}

	jsonpkg.ParseBody(r, schema)
	h.s.NewSchema(ctx, schema)
}
