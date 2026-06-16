package middleware

import (
	"context"
	"net/http"

	"github.com/resto-fnb/backend/internal/domain"
	"github.com/resto-fnb/backend/pkg/response"
)

type contextKey string

const ChainKey contextKey = "chain"

func Tenant(resolver func(slug string) (*domain.Chain, error)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			slug := r.Header.Get("X-Tenant-Slug")
			if slug == "" {
				response.Error(w, http.StatusBadRequest, "X-Tenant-Slug header is required")
				return
			}

			chain, err := resolver(slug)
			if err != nil {
				response.Error(w, http.StatusNotFound, "tenant not found")
				return
			}

			ctx := context.WithValue(r.Context(), ChainKey, chain)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetChain(r *http.Request) *domain.Chain {
	chain, _ := r.Context().Value(ChainKey).(*domain.Chain)
	return chain
}
