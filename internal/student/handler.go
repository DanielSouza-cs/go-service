package student

import (
	"context"
	"errors"
	"go-service/internal/response"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type ReportGenerator interface {
	GenerateReport(ctx context.Context, id int64) ([]byte, error)
}

func ReportHandler(svc ReportGenerator, lg *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			response.WriteError(w, "Invalid student ID format", http.StatusBadRequest)
			return
		}

		pdf, err := svc.GenerateReport(r.Context(), id)
		if err != nil {
			if errors.Is(err, ErrStudentNotFound) {
				response.WriteError(w, "Student not found", http.StatusNotFound)
				return
			}

			lg.Error("failed to generate report",
				zap.Int64("studentID", id),
				zap.Error(err),
			)
			response.WriteError(w, "An internal error occurred", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", `attachment; filename="student_`+idStr+`_report.pdf"`)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(pdf)
	}
}
