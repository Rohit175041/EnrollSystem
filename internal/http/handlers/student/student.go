package student

import (
	"context"
	"encoding/json"
	"github.com/rohit154041/students-api/internal/models"
	"github.com/rohit154041/students-api/internal/storage"
	"github.com/rohit154041/students-api/internal/utils/response"
	"log/slog"
	"net/http"
	"time"
)

func CreateStudent(w http.ResponseWriter, r *http.Request) {
	var student models.User

	// Decode request body into struct
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.Error(err))
		return
	}

	// Basic validation
	if student.Name == "" || student.Email == "" {
		response.WriteJSON(w, http.StatusBadRequest, response.Fail("Name and Email are required"))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := storage.GetCollection()

	_, err := collection.InsertOne(ctx, student)
	if err != nil {
		slog.Error("failed to insert student", slog.String("error", err.Error()))
		response.WriteJSON(w, http.StatusInternalServerError, response.Error(err))
		return
	}
	// Success response
	response.WriteJSON(w, http.StatusOK, response.Success("Student created successfully"))
}
