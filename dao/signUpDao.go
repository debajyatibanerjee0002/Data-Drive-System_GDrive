package dao

import (
	"go-drive-project/entities"
	"go-drive-project/logger"
)

func (db *DbConn) CheckUserExistsDao(emailOrUsername string) (bool, error) {
	logger.Log.Println("Inside CheckUserExistsDao")

	var count int64
	err := db.DB.Model(&entities.User{}).
		Where("email = ? OR username = ? AND delete_flag = 0", emailOrUsername, emailOrUsername).
		Count(&count).Error

	if err != nil {
		logger.Log.Println("Error in checking user exists:", err)
		return false, err
	}

	return count > 0, nil

}

func (db *DbConn) InsertUserDao(reqData entities.SignUpReqestEntity) (int, error) {
	logger.Log.Println("Inside InsertUserDao")

	newUser := entities.User{
		Username:      reqData.UserName,
		Email:         reqData.Email,
		Password_hash: reqData.Password,
	}
	if err := db.DB.Create(&newUser).Error; err != nil {
		logger.Log.Println("Error in inserting user:", err)
		return 0, err
	}

	return int(newUser.ID), nil
}

func (db *DbConn) GetUserDetailsDao(emailOrUsername string) (entities.SignUpReqestEntity, error) {
	logger.Log.Println("Inside GetUserDetailsDao")
	var response entities.SignUpReqestEntity
	var user entities.User
	err := db.DB.Where("email = ? OR username = ? AND delete_flag = 0", emailOrUsername, emailOrUsername).
		First(&user).Error

	if err != nil {
		logger.Log.Println("Error in getting user details:", err)
		return response, err
	}

	response.Email = user.Email
	response.UserName = user.Username
	response.Password = user.Password_hash
	response.UserId = int(user.ID)

	return response, nil
}
