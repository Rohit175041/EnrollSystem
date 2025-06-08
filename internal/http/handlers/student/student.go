package student

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"github.com/rohit154041/students-api/internal/types"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var student types.Student

		json.NewDecoder(r.Body).Decode(&student)
		slog.Info("creating student info")

		w.Write([]byte("welcome to students api"))
	}
}
