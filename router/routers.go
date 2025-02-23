package router

import (
	"go-drive-project/handlers"

	"github.com/gin-gonic/gin"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc gin.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"404 Not Found",
		"GET,POST,PUT,PATCH,DELETE,OPTIONS",
		"/",
		ThrowBlankResponse,
	},
	Route{
		"createUser",
		"POST",
		"/signup",
		handlers.SignupHandler,
	},
	Route{
		"loginUser",
		"POST",
		"/login",
		handlers.LoginHandler,
	},
	Route{
		"uploadFile",
		"POST",
		"/uploadFile",
		handlers.UploadFileHandler, // this uploadfile handle both create and update same file
	},
	Route{
		"getFolders",
		"GET",
		"/getFolders",
		handlers.GetFoldersHandler,
	},
	Route{
		"getFilesByUserId",
		"GET",
		"/getFilesByUserId",
		handlers.GetFilesByUserId,
	},
	Route{
		"deleteFileByUserId",
		"DELETE",
		"/deleteFileByUserId",
		handlers.DeleteFileByUserId,
	},
	Route{
		"createFolderByUserId",
		"POST",
		"/createFolderByUserId",
		handlers.CreateFolderByUserId,
	},
	Route{
		"deleteFolderByUserId",
		"DELETE",
		"/deleteFolderByUserId",
		handlers.DeleteFolderByUserId,
	},
}
