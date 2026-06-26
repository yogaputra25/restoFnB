package report

import (
	"net/http"

	"github.com/resto-fnb/backend/internal/domain"
	"github.com/resto-fnb/backend/internal/middleware"
	"github.com/resto-fnb/backend/pkg/response"
)

type ReportService interface {
	GetSalesReport(chainID, period, from, to string) (*domain.SalesReportResponse, error)
	GetSalesReportCSV(chainID, period, from, to string) ([]byte, string, error)
}

type Handler struct {
	svc ReportService
}

func NewHandler(svc ReportService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Sales(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	user := middleware.GetUser(r)
	if user == nil {
		response.Error(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	period := r.URL.Query().Get("period")
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	format := r.URL.Query().Get("format")

	if from == "" || to == "" {
		response.Error(w, http.StatusBadRequest, "from and to are required")
		return
	}

	if format == "csv" {
		data, filename, err := h.svc.GetSalesReportCSV(user.ChainID, period, from, to)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.Header().Set("Content-Type", "text/csv; charset=utf-8")
		w.Header().Set("Content-Disposition", "attachment; filename="+filename)
		w.Write(data)
		return
	}

	report, err := h.svc.GetSalesReport(user.ChainID, period, from, to)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, report)
}
