package httpserver

import (
	"encoding/json"
	"go-service/internal/student"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func NewRouter(svc student.ReportGenerator, lg *zap.Logger) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	r.Get("/api/v1/students/{id}/report", student.ReportHandler(svc, lg))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "route not found", http.StatusNotFound)
	})

	r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("go-service up " + time.Now().Format(time.RFC3339)))
	})
	return r
}
