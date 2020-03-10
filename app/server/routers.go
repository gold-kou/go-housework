package server

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/gold-kou/go-housework/app/server/handler"
	"github.com/gold-kou/go-housework/app/server/middleware"
)

// Route represents a route of routing
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes represent routes of routing
type Routes []Route

// NewRouter returns a new Router specified by Routes
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = middleware.Logger(handler, route.Name)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

var routes = Routes{
	{
		"Login",
		strings.ToUpper("Post"),
		"/login",
		handler.Login,
	},
	{
		"CreateUser",
		strings.ToUpper("Post"),
		"/user",
		handler.CreateUser,
	},
	{
		"DeleteUser",
		strings.ToUpper("Delete"),
		"/user",
		handler.DeleteUser,
	},
}
