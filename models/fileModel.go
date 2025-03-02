package models

import (
	"errors"
	"go-drive-project/config"
	"go-drive-project/dao"
	"go-drive-project/entities"
	"go-drive-project/logger"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
)

func UploadFileModel(fileHeader *multipart.FileHeader, folderID, userID uint) (entities.File, error) {
	logger.Log.Println("Inside uploadFileModel")
	db, err := config.ConnectMySqlDbSingleton()
	if err != nil {
		log.Println("Database connection failed:", err)
		return entities.File{}, err
	}

	dataAccess := dao.DbConn{DB: db}

	isFolderExists, err := dataAccess.CheckFolderExistsDao(folderID, userID)
	if err != nil || !isFolderExists {
		log.Println("Folder not found or does not belong to user")
		return entities.File{}, errors.New("folder not found")
	}

	getFolderDetails, err := dataAccess.GetFolderDetails(folderID)
	if err != nil {
		log.Println("Error in getting folder details:", err)
		return entities.File{}, err
	}
	get_root_folder, err := dataAccess.GetRootFolderNameDao(int64(userID))
	if err != nil {
		logger.Log.Println("something went wrong while fetching root folder: ", err)
		return entities.File{}, err
	}
	logger.Log.Println("get_root_folder: ", get_root_folder)
	storage_path := get_root_folder
	var storage_path_ids string

	// checking root folder or not
	if getFolderDetails.ParentID == 0 {
		// storage_path = storage_path + getFolderDetails.Name
		str := strconv.FormatUint(uint64(getFolderDetails.ID), 10)
		storage_path_ids = str
		logger.Log.Println("storage_path_ids: ", storage_path_ids)
	} else {
		getFullPath, getFullPathIds, err := dataAccess.GetFullPathDao(getFolderDetails.ID)
		if err != nil {
			log.Println("Error in getting folder details:", err)
			return entities.File{}, err
		}
		storage_path = getFullPath
		storage_path_ids = getFullPathIds
		logger.Log.Println("storage_path_ids: ", storage_path_ids)
	}

	if err := os.MkdirAll(storage_path, os.ModePerm); err != nil {
		log.Println("Failed to create storage directory:", err)
		return entities.File{}, err
	}

	folderP := "storage/" + storage_path
	// Save file locally
	filePath := filepath.Join(folderP, fileHeader.Filename)
	file, err := fileHeader.Open()
	if err != nil {
		log.Println("Failed to open file:", err)
		return entities.File{}, err
	}
	defer file.Close()

	outFile, err := os.Create(filePath)
	if err != nil {
		log.Println("Failed to create file on server:", err)
		return entities.File{}, err
	}
	defer outFile.Close()

	if _, err := outFile.ReadFrom(file); err != nil {
		log.Println("Failed to save file:", err)
		return entities.File{}, err
	}

	// Insert file metadata into MySQL
	fileDetails, err := dataAccess.InsertFolderDao(fileHeader, filePath, storage_path_ids, folderID, userID)
	if err != nil {
		log.Println("Failed to insert file metadata:", err)
		return entities.File{}, err
	}

	// Return file metadata
	return entities.File{
		ID:       uint(fileDetails.ID),
		Name:     fileHeader.Filename,
		Path:     filePath,
		FolderID: &folderID,
		UserID:   userID,
	}, nil
}

func GetFoldersModel(userID uint) ([][]entities.Folder, error) {
	logger.Log.Println("Inside GetFoldersModel")
	db, err := config.ConnectMySqlDbSingleton()
	if err != nil {
		log.Println("Database connection failed:", err)
		return nil, err
	}

	dataAccess := dao.DbConn{DB: db}
	folders, err := dataAccess.GetFoldersDao(userID)
	if err != nil {
		log.Println("Error getting folders:", err)
		return nil, err
	}

	return folders, nil
}

func GetFilesModel(userID uint) ([]entities.File, error) {
	logger.Log.Println("Inside GetFilesModel")
	db, err := config.ConnectMySqlDbSingleton()
	if err != nil {
		log.Println("Database connection failed:", err)
		return nil, err
	}

	dataAccess := dao.DbConn{DB: db}
	files, err := dataAccess.GetFilesDao(userID)
	if err != nil {
		log.Println("Error getting folders:", err)
		return nil, err
	}

	return files, nil
}

func DeleteFileByUserIdModel(fileID int64, userID int64) error {
	logger.Log.Println("Inside DeleteFileByUserIdModel")
	db, err := config.ConnectMySqlDbSingleton()
	if err != nil {
		log.Println("Database connection failed:", err)
		return err
	}

	dataAccess := dao.DbConn{DB: db}
	err = dataAccess.DeleteFileByUserIdDao(fileID, userID)
	if err != nil {
		log.Println("Error getting folders:", err)
		return err
	}

	return nil
}

func DeleteFolderByUserIdModel(folderID int64, userID int64) error {
	logger.Log.Println("Inside DeleteFolderByUserIdModel")
	db, err := config.ConnectMySqlDbSingleton()
	if err != nil {
		log.Println("Database connection failed:", err)
		return err
	}

	dataAccess := dao.DbConn{DB: db}
	err = dataAccess.DeleteFolderByUserIdDao(folderID, userID)
	if err != nil {
		log.Println("Error getting folders:", err)
		return err
	}

	return nil
}

func CreateRootFolderByUserIdModel(userId int) (entities.Folder, error) {
	logger.Log.Println("Inside CreateFolderByUserIdModel")
	rootFolderName := "root_" + strconv.Itoa(userId)

	db, err := config.ConnectMySqlDbSingleton()
	if err != nil {
		log.Println("Database connection failed:", err)
		return entities.Folder{}, err
	}

	dataAccess := dao.DbConn{DB: db}

	checkIsThereAnyFolderPresentOrNot, err := dataAccess.CheckIsThereAnyFolderPresentOrNotDao(userId)
	if err != nil {
		logger.Log.Println("something went wrong when checking checkIsThereAnyFolderPresentOrNot: ", err)
		return entities.Folder{}, err
	}
	if checkIsThereAnyFolderPresentOrNot {
		logger.Log.Println("Folder is present")
		return entities.Folder{}, errors.New("parent folder is present")
	}

	inserRootFolder, err := dataAccess.CreateRootFolderByUserIdModel(userId, rootFolderName)
	if err != nil {
		logger.Log.Println("something went wrong when creating root folder: ", err)
		return entities.Folder{}, err
	}
	logger.Log.Println("created root folder Id: ", inserRootFolder)

	folderP := "storage/" + rootFolderName
	err = CreateFolder(folderP)
	if err != nil {
		logger.Log.Println("Error creating folder:", err)
		return entities.Folder{}, err
	}
	return inserRootFolder, nil

}

func CreateFolderByUserIdModel(reqData entities.CreateFolderByUserIdRequestEntity) (entities.Folder, error) {
	logger.Log.Println("Inside CreateFolderByUserIdModel")

	db, err := config.ConnectMySqlDbSingleton()
	if err != nil {
		log.Println("Database connection failed:", err)
		return entities.Folder{}, err
	}

	dataAccess := dao.DbConn{DB: db}

	inserFolder, err := dataAccess.CreateFolderByUserIdDao(reqData)
	if err != nil {
		logger.Log.Println("something went wrong when creating root folder: ", err)
		return entities.Folder{}, err
	}
	logger.Log.Println("created folder Id: ", inserFolder)

	folderPath, err := dataAccess.GetFolderPathDao(reqData)
	if err != nil {
		logger.Log.Println("failed to get folder path")
		return entities.Folder{}, err
	}

	folderP := "storage/" + folderPath
	err = CreateFolder(folderP)
	if err != nil {
		logger.Log.Println("Error creating folder:", err)
		return entities.Folder{}, err
	}
	return inserFolder, nil

}

func CreateFolder(folderPath string) error {
	err := os.MkdirAll(folderPath, os.ModePerm) // os.ModePerm = 0777 (default permission)
	if err != nil {
		return err
	}
	logger.Log.Println("Folder created successfully:", folderPath)
	return nil
}
