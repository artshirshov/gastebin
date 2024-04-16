package router

import (
	hh "github.com/artshirshov/gastebin/internal/handler/healthcheck"
	ph "github.com/artshirshov/gastebin/internal/handler/paste"
	"github.com/artshirshov/gastebin/pkg/logger"
	"github.com/artshirshov/gastebin/pkg/rest"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	// CorsMaxAge Максимальное значение, не игнорируемое ни одним из основных браузеров.
	CorsMaxAge = 300
)

func New(
	pasteHandler ph.Handler,
	healthHandler hh.Handler,

) *chi.Mux {
	r := chi.NewRouter()

	// Настройка CORS
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"},
		// AllowedOrigins: []string{"https://*", "http://*"},
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           CorsMaxAge,
	}))

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", rest.JsonWrapperHandler(healthHandler.CheckHealth))

		r.Route("/pastes", func(r chi.Router) {
			r.Post("/", rest.JsonWrapperHandler(pasteHandler.CreatePaste))
			r.Route("/{hash}", func(r chi.Router) {
				r.Get("/", rest.JsonWrapperHandler(pasteHandler.GetPaste))
				r.Put("/", rest.JsonWrapperHandler(pasteHandler.UpdatePaste))
				r.Delete("/", rest.JsonWrapperHandler(pasteHandler.DeletePaste))
			})
		})
	})

	workDir, _ := os.Getwd()

	// Конфигурация Swagger
	openAPIDir := http.Dir(filepath.Join(workDir, "api/openapi"))
	FileServer(r, "/api/v1/swagger", openAPIDir)
	return r
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		logger.Log.Debug("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
