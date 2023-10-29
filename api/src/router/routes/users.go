package routes

import (
	"api/src/controllers"
)

// representa todas as rotas de usuarios que teremos dentro da api.
// as funções que vão lidar com as rotas vão ser armazenadas no package controllers.
var usersRoutes = []Route{
	{
		URI: "/users",
		Method: "POST",
		Func: controllers.CreateUser,
		RequireAuth: false,
	},
	{
		URI: "/users",
		Method: "GET",
		Func: controllers.GetUsers,
		RequireAuth: false,
	},	
	{
		URI: "/users/{userId}",
		Method: "GET",
		Func: controllers.GetUSerById,
		RequireAuth: false,
	},	
	{
		URI: "/users/{userId}",
		Method: "PUT",
		Func: controllers.UpdateUser,
		RequireAuth: false,
	},	
	{
		URI: "/users/{userId}",
		Method: "DELETE",
		Func: controllers.DeleteUser,
		RequireAuth: false,
	},
}