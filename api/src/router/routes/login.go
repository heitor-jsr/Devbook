package routes

import "api/src/controllers"

var loginRoute = Route{
	URI: "/login",
	Method: "POST",
	Func: controllers.Login,
	RequireAuth: false,
}