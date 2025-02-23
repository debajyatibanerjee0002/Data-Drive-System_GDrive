package handlers

import (
	"go-drive-project/entities"
	"go-drive-project/logger"
	"go-drive-project/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func UploadFileHandler(c *gin.Context) {
	log.Println("Inside UploadFileHandler")

	userIDStr := c.PostForm("user_id")
	folderIDStr := c.PostForm("folder_id")

	// Validate required fields
	if userIDStr == "" || folderIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id and folder_id are required"})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		log.Println("Invalid user_id:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}

	folderID, err := strconv.ParseUint(folderIDStr, 10, 32)
	if err != nil {
		log.Println("Invalid folder_id:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid folder_id"})
		return
	}

	// Get uploaded file
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		log.Println("File upload failed:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}
	defer file.Close()

	// Call file upload model function
	uploadedFile, err := models.UploadFileModel(fileHeader, uint(folderID), uint(userID))
	if err != nil {
		log.Println("Error uploading file:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Success response
	c.JSON(http.StatusOK, gin.H{"success": true, "file": uploadedFile})
}

func GetFoldersHandler(c *gin.Context) {
	logger.Log.Println("Inside GetFoldersHandler")

	userIDStr := c.Query("userid")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userid is required"})
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 32)
	if err != nil {
		log.Println("Invalid user_id:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}

	folders, err := models.GetFoldersModel(uint(userID))
	if err != nil {
		log.Println("Error getting folders:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(folders) == 0 {
		ThrowGetFoldersResponse(c, true, http.StatusNoContent, "Folders retrieved successfully", folders)
		return
	}
	ThrowGetFoldersResponse(c, true, http.StatusOK, "Folders retrieved successfully", folders)

}

func ThrowGetFoldersResponse(c *gin.Context, success bool, status int, message string, response [][]entities.Folder) {
	var res entities.FolderResponseEntity
	res.Sucess = success
	res.Message = message
	res.Folders = response
	c.JSON(status, res)
}

func GetFilesByUserId(c *gin.Context) {
	logger.Log.Println("Inside GetFilesByUserId")

	userIDStr := c.Query("userid")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userid is required"})
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 32)
	if err != nil {
		log.Println("Invalid user_id:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}

	files, err := models.GetFilesModel(uint(userID))
	if err != nil {
		log.Println("Error getting folders:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ThrowGetFilesResponse(c, true, http.StatusOK, "Folders retrieved successfully", files)

}

func ThrowGetFilesResponse(c *gin.Context, success bool, status int, message string, response []entities.File) {
	var res entities.FilesResponseEntity
	res.Sucess = success
	res.Message = message
	res.Files = response
	c.JSON(status, res)
}

func DeleteFileByUserId(c *gin.Context) {
	logger.Log.Println("Inside DeleteFileByUserId")

	userIDStr := c.Query("userid")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userid is required"})
		return
	}
	fileIdStr := c.Query("fileid")
	if fileIdStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userid is required"})
		return
	}
	userID, err := strconv.ParseInt(userIDStr, 10, 32)
	if err != nil {
		log.Println("Invalid user_id:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}
	fileID, err := strconv.ParseInt(fileIdStr, 10, 32)
	if err != nil {
		log.Println("Invalid formId:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid formId"})
		return
	}

	err = models.DeleteFileByUserIdModel(fileID, userID)
	if err != nil {
		log.Println("Error getting folders:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ThrowDeleteFileByUserIdResponse(c, true, http.StatusOK, "Successfully deleted file")

}

func ThrowDeleteFileByUserIdResponse(c *gin.Context, success bool, status int, message string) {
	var res entities.DeleteFileByUserIdResponseEntity
	res.Sucess = success
	res.Message = message
	c.JSON(status, res)
}

func CreateFolderByUserId(c *gin.Context) {
	logger.Log.Println("Inside CreateFolderByUserId")

	var reqData = entities.CreateFolderByUserIdRequestEntity{}
	jsonError := reqData.FromJSON(c.Request.Body)
	if jsonError != nil {
		logger.Log.Println("Error in parsing request body: ", jsonError)
		ThrowCreateFolderByUserIdResponse(c, false, http.StatusInternalServerError, "Failed to encode request body", entities.Folder{})
		return
	}

	// Need to create root folder
	if len(reqData.FolderPathIds) == 0 {
		response, err := models.CreateRootFolderByUserIdModel(reqData.Userid)
		if err != nil {
			logger.Log.Println("something went wrong while creating root folder: ", err)
			ThrowCreateFolderByUserIdResponse(c, false, http.StatusInternalServerError, "something went wrong while creating root folder or root folder is present", entities.Folder{})
			return
		}
		ThrowCreateFolderByUserIdResponse(c, true, http.StatusOK, "Successfully created folder", response)
	} else {
		logger.Log.Println("Creating...")
		response, err := models.CreateFolderByUserIdModel(reqData)
		if err != nil {
			logger.Log.Println("something went wrong while creating folder: ", err)
			ThrowCreateFolderByUserIdResponse(c, false, http.StatusInternalServerError, "something went wrong while creating folder", entities.Folder{})
			return
		}
		ThrowCreateFolderByUserIdResponse(c, true, http.StatusOK, "Successfully created folder", response)
	}
}

func ThrowCreateFolderByUserIdResponse(c *gin.Context, success bool, status int, message string, response entities.Folder) {
	var res entities.CreateFolderByUserIdResponseEntity
	res.Sucess = success
	res.Message = message
	res.Response = append(res.Response, response)
	c.JSON(status, res)
}

func DeleteFolderByUserId(c *gin.Context) {
	logger.Log.Println("Inside DeleteFolderByUserId")

	userIDStr := c.Query("userid")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userid is required"})
		return
	}
	folderIdStr := c.Query("folderid")
	if folderIdStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userid is required"})
		return
	}
	userID, err := strconv.ParseInt(userIDStr, 10, 32)
	if err != nil {
		log.Println("Invalid user_id:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}
	folderID, err := strconv.ParseInt(folderIdStr, 10, 32)
	if err != nil {
		log.Println("Invalid formId:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid formId"})
		return
	}

	err = models.DeleteFolderByUserIdModel(folderID, userID)
	if err != nil {
		log.Println("Error getting folders:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ThrowDeleteFileByUserIdResponse(c, true, http.StatusOK, "Successfully deleted folder")
}
