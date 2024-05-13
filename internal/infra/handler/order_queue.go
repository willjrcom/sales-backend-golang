package handlerimpl

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	orderprocessentity "github.com/willjrcom/sales-backend-go/internal/domain/order_process"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	orderqueuedto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order_queue"
	orderqueueusecases "github.com/willjrcom/sales-backend-go/internal/usecases/order_queue"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type handlerQueueImpl struct {
	s *orderqueueusecases.Service
}

func NewHandlerQueue(queueService *orderqueueusecases.Service) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerQueueImpl{
		s: queueService,
	}

	c.With().Group(func(c chi.Router) {
		c.Post("/start", h.handlerStartQueue)
		c.Post("/finish/{id}", h.handlerFinishQueue)
		c.Get("/{id}", h.handlerGetQueueByID)
		c.Get("/by-group-item/{id}", h.handlerGetQueuesByGroupItemId)
		c.Get("/all", h.handlerGetAllQueues)
	})

	return handler.NewHandler("/queue", c)
}

func (h *handlerQueueImpl) handlerStartQueue(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoQueue := &orderqueuedto.StartQueueInput{}
	if err := jsonpkg.ParseBody(r, dtoQueue); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: err.Error()})
		return
	}

	id, err := h.s.StartQueue(ctx, dtoQueue)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: id})
}

func (h *handlerQueueImpl) handlerFinishQueue(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	process := &orderprocessentity.OrderProcess{}
	if err := jsonpkg.ParseBody(r, process); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: err.Error()})
		return
	}

	if err := h.s.FinishQueue(ctx, process); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerQueueImpl) handlerGetQueueByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	queue, err := h.s.GetQueueById(ctx, dtoId)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: queue})
}

func (h *handlerQueueImpl) handlerGetQueuesByGroupItemId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	queues, err := h.s.GetQueuesByGroupItemId(ctx, dtoId)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: queues})
}

func (h *handlerQueueImpl) handlerGetAllQueues(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	queues, err := h.s.GetAllQueues(ctx)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: queues})
}
