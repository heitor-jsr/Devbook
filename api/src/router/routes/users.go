package routes

import (
	"api/src/controllers"
	"api/src/database"
	"database/sql"
	"log"
)

var db *sql.DB = connectDB()

func connectDB() *sql.DB {
    db, erro := database.Connect()
    if erro != nil {
        log.Fatal(erro)
    }
    return db
}

// representa todas as rotas de usuarios que teremos dentro da api.
// as funções que vão lidar com as rotas vão ser armazenadas no package controllers.
var usersRoutes = []Route{
	{
		URI: "/users",
		Method: "POST",
		Func: controllers.CreateUserFactory(db),
		RequireAuth: false,
	},
	{
		URI: "/users",
		Method: "GET",
		Func: controllers.GetUsers,
		RequireAuth: true,
	},	
	{
		URI: "/users/{userId}",
		Method: "GET",
		Func: controllers.GetUSerById,
		RequireAuth: true,
	},	
	{
		URI: "/users/{userId}",
		Method: "PUT",
		Func: controllers.UpdateUser,
		RequireAuth: true,
	},	
	{
		URI: "/users/{userId}",
		Method: "DELETE",
		Func: controllers.DeleteUser,
		RequireAuth: true,
	},
	{
		URI: "/users/{userId}/follow",
		Method: "POST",
		Func: controllers.FollowUser,
		RequireAuth: true,
	},
	{
		URI: "/users/{userId}/unfollow",
		Method: "POST",
		Func: controllers.UnfollowUser,
		RequireAuth: true,
	},
	{
		URI: "/users/{userId}/followers",
		Method: "GET",
		Func: controllers.GetFollowers,
		RequireAuth: true,
	},
	{
		URI: "/users/{userId}/followings",
		Method: "GET",
		Func: controllers.GetFollowings,
		RequireAuth: true,
	},
	{
		URI: "/users/{userId}/change-password",
		Method: "POST",
		Func: controllers.ChangePassword,
		RequireAuth: true,
	},
}