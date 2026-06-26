package report

import (
	"bytes"
	"encoding/csv"
	"fmt"

	"github.com/resto-fnb/backend/internal/domain"
)

type ReportRepo interface {
	GetSalesReport(chainID, period, from, to string) (*domain.SalesReportResponse, error)
}

type Usecase struct {
	repo ReportRepo
}

func NewUsecase(repo ReportRepo) *Usecase {
	return &Usecase{repo: repo}
}

func (uc *Usecase) GetSalesReport(chainID, period, from, to string) (*domain.SalesReportResponse, error) {
	if period == "" {
		period = "daily"
	}
	return uc.repo.GetSalesReport(chainID, period, from, to)
}

func (uc *Usecase) GetSalesReportCSV(chainID, period, from, to string) ([]byte, string, error) {
	report, err := uc.GetSalesReport(chainID, period, from, to)
	if err != nil {
		return nil, "", err
	}

	var buf bytes.Buffer
	buf.WriteString("\xef\xbb\xbf")

	w := csv.NewWriter(&buf)
	w.Comma = ';'

	w.Write([]string{"Periode", "Kategori", "Menu Item", "Terjual", "Pendapatan"})

	dateLabel := fmt.Sprintf("%s - %s", report.From, report.To)
	for _, cat := range report.Categories {
		for _, item := range cat.Items {
			w.Write([]string{
				dateLabel,
				cat.CategoryName,
				item.Name,
				fmt.Sprintf("%d", item.Quantity),
				fmt.Sprintf("%.0f", item.Revenue),
			})
		}
	}
	w.Flush()

	if err := w.Error(); err != nil {
		return nil, "", err
	}

	filename := fmt.Sprintf("sales-report-%s-%s.csv", report.From, report.To)
	return buf.Bytes(), filename, nil
}
