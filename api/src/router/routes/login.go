package routes

var loginRoute = Route{
	URI: "/login",
	Method: "POST",
	Func: controllers.Login,
	RequireAuth: false,
}