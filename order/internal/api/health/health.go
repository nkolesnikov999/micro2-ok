package health

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/nkolesnikov999/micro2-OK/platform/pkg/logger"
)

func Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		resp := map[string]any{
			"status":  "SERVING",
			"service": "order-api",
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			logger.Error(r.Context(), "failed to write health response", zap.Error(err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	})
}
