package handlers

import (
	"go-drive-project/entities"
	"go-drive-project/logger"
	"go-drive-project/models"
	"go-drive-project/utility"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func SignupHandler(c *gin.Context) {
	logger.Log.Println("Inside SignupHandler")
	var reqData = entities.SignUpReqestEntity{}
	jsonError := reqData.FromJSON(c.Request.Body)
	if jsonError != nil {
		logger.Log.Println("Error in parsing request body: ", jsonError)
		ThrowSignupResponse(c, false, http.StatusInternalServerError, "Failed to encode request body", "")
		return
	}
	if !strings.Contains(reqData.Email, "@") {
		logger.Log.Println("Invalid Email")
		ThrowSignupResponse(c, false, http.StatusBadRequest, "Invalid Email", "")
	}
	if reqData.Email == "" || reqData.Password == "" {
		logger.Log.Println("Email or Password is missing")
		ThrowSignupResponse(c, false, http.StatusBadRequest, "Email and Password are required", "")
		return
	}

	userName, index := utility.GetUserName(reqData.Email)
	if index <= 0 {
		logger.Log.Println("Invalid Email")
		ThrowSignupResponse(c, false, http.StatusBadRequest, "Invalid Email", "")
		return
	}

	reqData.UserName = userName

	isUserExists, token, err := models.SignupModel(reqData)
	if err != nil {
		logger.Log.Println("Error in checking user exists: ", err)
		ThrowSignupResponse(c, false, http.StatusInternalServerError, "Failed to check user exists", "")
		return
	}
	if isUserExists {
		logger.Log.Println("User already exists")
		ThrowSignupResponse(c, false, http.StatusConflict, "User already exists", "")
		return
	}
	ThrowSignupResponse(c, true, http.StatusOK, "User signed up successfully", token)
}

func ThrowSignupResponse(c *gin.Context, success bool, status int, message string, token string) {
	var response entities.SignUpResponseEntity
	response.Success = success
	response.Message = message
	response.Token = token
	c.JSON(status, response)
}

func LoginHandler(c *gin.Context) {
	logger.Log.Println("Inside LoginHandler")
	var reqData = entities.LoginRequestEntity{}
	jsonError := reqData.FromJSON(c.Request.Body)
	if jsonError != nil {
		logger.Log.Println("Error in parsing request body: ", jsonError)
		ThrowLoginResponse(c, false, http.StatusInternalServerError, "Failed to encode request body", "")
		return
	}
	if reqData.EmailOrUsername == "" || reqData.Password == "" {
		logger.Log.Println("Email or Password is missing")
		ThrowLoginResponse(c, false, http.StatusBadRequest, "Email and Password are required", "")
		return
	}
	isValidUser, token, err := models.LoginModel(reqData)
	if err != nil {
		logger.Log.Println("Error in checking user exists: ", err)
		ThrowLoginResponse(c, false, http.StatusInternalServerError, "Failed to check user exists", "")
		return
	}
	if !isValidUser {
		logger.Log.Println("Invalid Email or Password")
		ThrowLoginResponse(c, false, http.StatusUnauthorized, "Invalid Email or Password", "")
		return
	}
	ThrowLoginResponse(c, true, http.StatusOK, "logged in successfully", token)
}

func ThrowLoginResponse(c *gin.Context, success bool, status int, message string, token string) {
	var response entities.LoginResponseEntity
	response.Success = success
	response.Message = message
	response.Token = token
	c.JSON(status, response)
}
