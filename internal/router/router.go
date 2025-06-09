package router

import (
	"github.com/gorilla/mux"
	"github.com/rohit154041/students-api/internal/http/handlers/student"
)

func Init() *mux.Router {
	router := mux.NewRouter()

	// CRUD routes for students
	router.HandleFunc("/api/students", student.CreateStudent).Methods("POST")
	router.HandleFunc("/api/students", student.GetStudents).Methods("GET")
	router.HandleFunc("/api/students/{id}", student.GetStudent).Methods("GET")
	router.HandleFunc("/api/students/{id}", student.UpdateStudent).Methods("PUT")
	router.HandleFunc("/api/students/{id}", student.DeleteStudent).Methods("DELETE")

	return router
}
