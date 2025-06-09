package router

import (
	"net/http"
	"github.com/rohit154041/students-api/internal/http/handlers/student"
)

func Init() *http.ServeMux {
	mux := http.NewServeMux()

	// Separate GET and POST routes
	mux.HandleFunc("POST /api/students", student.CreateStudent)
	// mux.HandleFunc("POST /api/students", student.PostHandler())

	return mux
}
