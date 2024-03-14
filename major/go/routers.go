/*
 * MySocial API
 *
 * Social network. There are users on the social network who can register, create, update and delete posts, get their own wall with posts, view walls with posts from other users, like and view statistics on posts in the form of the number of likes and views, as well as leave comments on posts and view them.
 *
 * API version: 1.0
 * Contact: neazhazha@edu.hse.ru
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},

	Route{
		"UserAuth",
		strings.ToUpper("Get"),
		"/user/auth",
		UserAuth,
	},

	Route{
		"UserReg",
		strings.ToUpper("Post"),
		"/user/registration",
		UserReg,
	},

	Route{
		"UserUpdate",
		strings.ToUpper("Put"),
		"/user/update",
		UserUpdate,
	},
}
