package student

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rohit154041/students-api/internal/models"
	"github.com/rohit154041/students-api/internal/storage"
	"github.com/rohit154041/students-api/internal/utils/response"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateStudent creates a new student record
func CreateStudent(w http.ResponseWriter, r *http.Request) {
	var student models.User

	// Decode JSON body
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.Error(err))
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

	// ✅ Check if a student with the same email already exists
	var existing models.User
	err := collection.FindOne(ctx, bson.M{"email": student.Email}).Decode(&existing)
	if err == nil {
		// A record was found → email already exists
		response.WriteJSON(w, http.StatusConflict, response.Fail("Student with this email already exists"))
		return
	}

	if err != mongo.ErrNoDocuments {
		// An error occurred other than "no document found"
		slog.Error("error checking for existing student", slog.String("error", err.Error()))
		response.WriteJSON(w, http.StatusInternalServerError, response.Error(err))
		return
	}

	// Proceed to insert new student
	result, err := collection.InsertOne(ctx, student)
	if err != nil {
		slog.Error("failed to insert student", slog.String("error", err.Error()))
		response.WriteJSON(w, http.StatusInternalServerError, response.Error(err))
		return
	}

	student.ID = result.InsertedID.(primitive.ObjectID)

	response.WriteJSON(w, http.StatusCreated, response.Success("Student created successfully"))
}


// GetStudents returns all students
func GetStudents(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := storage.GetCollection()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		slog.Error("failed to fetch students", slog.String("error", err.Error()))
		response.WriteJSON(w, http.StatusInternalServerError, response.Error(err))
		return
	}
	defer cursor.Close(ctx)

	var students []models.UserResponse
	if err := cursor.All(ctx, &students); err != nil {
		slog.Error("failed to decode students", slog.String("error", err.Error()))
		response.WriteJSON(w, http.StatusInternalServerError, response.Error(err))
		return
	}

	response.WriteJSON(w, http.StatusOK, response.Success("Students retrieved successfully",students))
}

// GetStudent returns a single student by ID
func GetStudent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.Fail("Invalid student ID"))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := storage.GetCollection()

	var student models.UserResponse
	err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&student)
	if err != nil {
		response.WriteJSON(w, http.StatusNotFound, response.Fail("Student not found"))
		return
	}

	response.WriteJSON(w, http.StatusOK, response.Success("Student retrieved successfully",student))
}

// UpdateStudent updates a student by ID
func UpdateStudent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.Fail("Invalid student ID"))
		return
	}

	var student models.UserResponse
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.Error(err))
		return
	}

	// Optional: validate input fields, at least one field must be provided
	if student.Name == "" && student.Email == "" {
		response.WriteJSON(w, http.StatusBadRequest, response.Fail("At least one field (Name or Email) must be provided to update"))
		return
	}

	updateFields := bson.M{}
	if student.Name != "" {
		updateFields["name"] = student.Name
	}
	if student.Email != "" {
		updateFields["email"] = student.Email
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := storage.GetCollection()

	result, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updateFields})
	if err != nil {
		slog.Error("failed to update student", slog.String("error", err.Error()))
		response.WriteJSON(w, http.StatusInternalServerError, response.Error(err))
		return
	}

	if result.MatchedCount == 0 {
		response.WriteJSON(w, http.StatusNotFound, response.Fail("Student not found"))
		return
	}

	// Return updated student data
	var updatedStudent models.UserResponse
	err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&updatedStudent)
	if err != nil {
		slog.Error("failed to fetch updated student", slog.String("error", err.Error()))
		response.WriteJSON(w, http.StatusInternalServerError, response.Error(err))
		return
	}

	response.WriteJSON(w, http.StatusOK, response.Success("Student updated successfully",updatedStudent))
}

// DeleteStudent deletes a student by ID
func DeleteStudent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.Fail("Invalid student ID"))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := storage.GetCollection()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		slog.Error("failed to delete student", slog.String("error", err.Error()))
		response.WriteJSON(w, http.StatusInternalServerError, response.Error(err))
		return
	}

	if result.DeletedCount == 0 {
		response.WriteJSON(w, http.StatusNotFound, response.Fail("Student not found"))
		return
	}

	response.WriteJSON(w, http.StatusOK, response.Success("Student deleted successfully"))
}
