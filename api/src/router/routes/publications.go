package routes

import "api/src/controllers"

var publicationsRoutes = []Route {
	{
		URI: "/publications",
		Method: "POST",
		Func: controllers.CreatePublication,
		RequireAuth: false,
	},
	{
		URI: "/publications",
		Method: "GET",
		Func: controllers.GetPublication,
		RequireAuth: false,
	},
	{
		URI: "/publications/{publicationId}",
		Method: "GET",
		Func: controllers.GetPublicationById,
		RequireAuth: false,
	},
	{
		URI: "/publications/{publicationId}",
		Method: "PUT",
		Func: controllers.UpdatePublication,
		RequireAuth: false,
	},
	{
		URI: "/publications/{publicationId}",
		Method: "DELETE",
		Func: controllers.DeletePublication,
		RequireAuth: false,
	},
	{
		URI: "/users/{userId}/publications",
		Method: "GET",
		Func: controllers.GetPublicationFromUser,
		RequireAuth: false,
	},
}