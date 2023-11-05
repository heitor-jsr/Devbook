package routes

import "api/src/controllers"

var publicationsRoutes = []Route {
	{
		URI: "/publications",
		Method: "POST",
		Func: controllers.CreatePublication,
		RequireAuth: true,
	},
	{
		URI: "/publications",
		Method: "GET",
		Func: controllers.GetPublication,
		RequireAuth: true,
	},
	{
		URI: "/publications/{publicationId}",
		Method: "GET",
		Func: controllers.GetPublicationById,
		RequireAuth: true,
	},
	{
		URI: "/publications/{publicationId}",
		Method: "PUT",
		Func: controllers.UpdatePublication,
		RequireAuth: true,
	},
	{
		URI: "/publications/{publicationId}",
		Method: "DELETE",
		Func: controllers.DeletePublication,
		RequireAuth: true,
	},
	{
		URI: "/users/{userId}/publications",
		Method: "GET",
		Func: controllers.GetPublicationFromUser,
		RequireAuth: true,
	},
	{
		URI: "/publications/{publicationId}/like",
		Method: "POST",
		Func: controllers.LikePublication,
		RequireAuth: true,
	},
	{
		URI: "/publications/{publicationId}/deslike",
		Method: "POST",
		Func: controllers.DeslikePublication,
		RequireAuth: true,
	},
}