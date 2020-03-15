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
	{
		"CreateFamily",
		strings.ToUpper("Post"),
		"/family",
		handler.CreateFamily,
	},
	{
		"ShowFamily",
		strings.ToUpper("Get"),
		"/family",
		handler.ShowFamily,
	},
	{
		"UpdateFamily",
		strings.ToUpper("Put"),
		"/family",
		handler.UpdateFamily,
	},
	{
		"DeleteFamily",
		strings.ToUpper("Delete"),
		"/family",
		handler.DeleteFamily,
	},
	{
		"RegisterFamilyMember",
		strings.ToUpper("Post"),
		"/family/member",
		handler.RegisterFamilyMember,
	},
	{
		"DeleteFamilyMember",
		strings.ToUpper("Delete"),
		"/family/member/{member_id}",
		handler.DeleteFamilyMember,
	},
	{
		"ListFamilyMembers",
		strings.ToUpper("Get"),
		"/family/members",
		handler.ListFamilyMembers,
	}, /*
		{
			"CreateTask",
			strings.ToUpper("Post"),
			"/task",
			handler.CreateTask,
		},
		{
			"UpdateTask",
			strings.ToUpper("Put"),
			"/task",
			handler.UpdateTask,
		},
		{
			"DeleteTask",
			strings.ToUpper("Delete"),
			"/task/{task_id}",
			handler.DeleteTask,
		},
		{
			"ListTasks",
			strings.ToUpper("Get"),
			"/tasks",
			handler.ListTasks,
		},*/
}
