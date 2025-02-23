package models

import (
	"go-drive-project/config"
	"go-drive-project/dao"
	"go-drive-project/entities"
	"go-drive-project/logger"
	"go-drive-project/utility"
)

func SignupModel(reqData entities.SignUpReqestEntity) (bool, string, error) {
	logger.Log.Println("Inside SignupModel")
	logger.Log.Println("Request Data: ", reqData)

	db, err := config.ConnectMySqlDbSingleton()

	if err != nil {
		logger.Log.Println("database connection failure", err)
		return false, "", err
	}
	dataAccess := dao.DbConn{DB: db}
	isUserExists, err := dataAccess.CheckUserExistsDao(reqData.Email)
	if err != nil {
		logger.Log.Println("Error in checking user exists: ", err)
		return false, "", err
	}
	if isUserExists {
		logger.Log.Println("User already exists")
		return true, "", nil
	}
	encryptedPassword, err := utility.HashPassword(reqData.Password)
	if err != nil {
		logger.Log.Println("Error in encrypting password: ", err)
		return false, "", err
	}
	reqData.Password = encryptedPassword
	newUserId, err := dataAccess.InsertUserDao(reqData)
	if err != nil {
		logger.Log.Println("Error in inserting user: ", err)
		return false, "", err
	}
	token, err := utility.GenerateJWT(reqData.UserName, reqData.Email, newUserId)
	if err != nil {
		logger.Log.Println("Error in generating token: ", err)
		return false, "", err
	}
	return false, token, nil
}

func LoginModel(reqData entities.LoginRequestEntity) (bool, string, error) {
	logger.Log.Println("Inside LoginModel")
	logger.Log.Println("Request Data: ", reqData)

	db, err := config.ConnectMySqlDbSingleton()

	if err != nil {
		logger.Log.Println("database connection failure", err)
		return false, "", err
	}
	dataAccess := dao.DbConn{DB: db}
	isUserExists, err := dataAccess.CheckUserExistsDao(reqData.EmailOrUsername)
	if err != nil {
		logger.Log.Println("Error in checking user exists: ", err)
		return false, "", err
	}
	if !isUserExists {
		logger.Log.Println("User doesn't exists")
		return false, "", nil
	}
	getUserDetails, err := dataAccess.GetUserDetailsDao(reqData.EmailOrUsername)
	if err != nil {
		logger.Log.Println("Error in getting user details: ", err)
		return false, "", err
	}
	logger.Log.Println("User details: ", getUserDetails)
	isPasswordMatched := utility.CheckPasswordHash(reqData.Password, getUserDetails.Password)
	if !isPasswordMatched {
		logger.Log.Println("Password doesn't match")
		return false, "", nil
	}
	token, err := utility.GenerateJWT(getUserDetails.UserName, getUserDetails.Email, getUserDetails.UserId)
	if err != nil {
		logger.Log.Println("Error in generating token: ", err)
		return false, "", err
	}
	return true, token, nil
}
