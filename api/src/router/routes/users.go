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

var uc = controllers.NewUserController(db)

// representa todas as rotas de usuarios que teremos dentro da api.
// as funções que vão lidar com as rotas vão ser armazenadas no package controllers.
var usersRoutes = []Route{
	{
		URI: "/users",
		Method: "POST",
		Func: uc.CreateUser,
		RequireAuth: false,
	},
	{
		URI: "/users",
		Method: "GET",
		Func: uc.GetUsers,
		RequireAuth: true,
	},	
	{
		URI: "/users/{userId}",
		Method: "GET",
		Func: uc.GetUSerById,
		RequireAuth: true,
	},	
	{
		URI: "/users/{userId}",
		Method: "PUT",
		Func: uc.UpdateUser,
		RequireAuth: true,
	},	
	{
		URI: "/users/{userId}",
		Method: "DELETE",
		Func: uc.DeleteUser,
		RequireAuth: true,
	},
	{
		URI: "/users/{userId}/follow",
		Method: "POST",
		Func: uc.FollowUser,
		RequireAuth: true,
	},
	{
		URI: "/users/{userId}/unfollow",
		Method: "POST",
		Func: uc.UnfollowUser,
		RequireAuth: true,
	},
	{
		URI: "/users/{userId}/followers",
		Method: "GET",
		Func: uc.GetFollowers,
		RequireAuth: true,
	},
	{
		URI: "/users/{userId}/followings",
		Method: "GET",
		Func: uc.GetFollowings,
		RequireAuth: true,
	},
	{
		URI: "/users/{userId}/change-password",
		Method: "POST",
		Func: uc.ChangePassword,
		RequireAuth: true,
	},
}