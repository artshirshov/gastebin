package paste

import (
	"context"
	"encoding/json"
	model "github.com/artshirshov/gastebin/internal/model/paste"
	service "github.com/artshirshov/gastebin/internal/service/paste"
	hasher "github.com/artshirshov/gastebin/pkg/hash"
	"github.com/artshirshov/gastebin/pkg/logger"
	"github.com/artshirshov/gastebin/pkg/operation"
	"github.com/artshirshov/gastebin/pkg/rest"
	"go.uber.org/zap"
	"net/http"
)

type Handler interface {
	GetPaste(w http.ResponseWriter, r *http.Request) (rest.ResponseHolder[model.ResponseDto], error)
	CreatePaste(w http.ResponseWriter, r *http.Request) (rest.ResponseHolder[model.ResponseDto], error)
	UpdatePaste(w http.ResponseWriter, r *http.Request) (rest.ResponseHolder[model.ResponseDto], error)
	DeletePaste(w http.ResponseWriter, r *http.Request) (rest.ResponseHolder[operation.ResponseDto], error)
}

type handler struct {
	service service.Service
}

func NewHandler(service service.Service) Handler {
	return &handler{
		service: service,
	}
}

func (h *handler) GetPaste(w http.ResponseWriter, r *http.Request) (rest.ResponseHolder[model.ResponseDto], error) {
	ctx := r.Context()

	hash := r.PathValue("hash")

	dto, err := h.service.GetPasteByHash(ctx, hash)
	return rest.StatusOk(dto), err
}

func (h *handler) CreatePaste(w http.ResponseWriter, r *http.Request) (rest.ResponseHolder[model.ResponseDto], error) {
	ctx := r.Context()
	ctx = context.WithValue(ctx, "remoteAddr", hasher.GetUserIP(r))

	reqDto := &model.RequestDto{}

	err := json.NewDecoder(r.Body).Decode(&reqDto)
	if err != nil {
		logger.Log.With(zap.Error(err)).Error("error during read request dto")
		return rest.ResponseHolder[model.ResponseDto]{}, err
	}

	dto, err := h.service.CreatePaste(ctx, *reqDto)
	return rest.StatusCreated(dto), err
}

func (h *handler) UpdatePaste(w http.ResponseWriter, r *http.Request) (rest.ResponseHolder[model.ResponseDto], error) {
	ctx := r.Context()

	hash := r.PathValue("hash")

	reqDto := &model.RequestDto{}

	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&reqDto)
	if err != nil {
		logger.Log.With(zap.Error(err)).Error("error during read request dto")
		return rest.ResponseHolder[model.ResponseDto]{}, err
	}

	dto, err := h.service.UpdatePaste(ctx, hash, *reqDto)
	return rest.StatusOk(dto), err
}

func (h *handler) DeletePaste(w http.ResponseWriter, r *http.Request) (rest.ResponseHolder[operation.ResponseDto], error) {
	ctx := r.Context()

	hash := r.PathValue("hash")

	return rest.StatusOk(operation.ResponseDto{Message: "Deleted"}), h.service.DeletePaste(ctx, hash)
}
