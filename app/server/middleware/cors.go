package middleware

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// CORS provides Cross-Origin Resource Sharing middleware.
func CORS(router *mux.Router) http.Handler {
	options := []handlers.CORSOption{
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		handlers.AllowedOrigins([]string{"*"}),
	}
	return handlers.CORS(options...)(router)
}
