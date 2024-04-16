package healthcheck

import (
	"github.com/artshirshov/gastebin/pkg/operation"
	"github.com/artshirshov/gastebin/pkg/rest"
	"net/http"
)

type Handler interface {
	CheckHealth(w http.ResponseWriter, r *http.Request) (rest.ResponseHolder[operation.ResponseDto], error)
}

type handler struct{}

func NewHandler() Handler {
	return &handler{}
}

func (h *handler) CheckHealth(_ http.ResponseWriter, _ *http.Request) (
	rest.ResponseHolder[operation.ResponseDto],
	error,
) {
	return rest.StatusOk(operation.ResponseDto{Message: "Healthy"}), nil
}
