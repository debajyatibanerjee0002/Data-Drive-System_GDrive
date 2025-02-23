package dao

import (
	"go-drive-project/entities"
	"go-drive-project/logger"
	"mime/multipart"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

func (db *DbConn) CheckFolderExistsDao(folderId, userId uint) (bool, error) {
	logger.Log.Println("Inside CheckFolderExistsDao")

	var count int64
	err := db.DB.Model(&entities.Folder{}).
		Where("id = ? AND user_id = ? AND delete_flag = 0", folderId, userId).
		Count(&count).Error

	if err != nil {
		logger.Log.Println("Error in checking folder exists:", err)
		return false, err
	}

	return count > 0, nil
}

// func (db *DbConn) InsertFolderDao(fileHeader *multipart.FileHeader, filePath string, folderID, userID uint) (entities.File, error) {
// 	logger.Log.Println("Inside InsertFolderDao")

// 	file := entities.File{
// 		Name:     fileHeader.Filename,
// 		Path:     filePath,
// 		FolderID: &folderID,
// 		UserID:   userID,
// 	}
// 	if err := db.DB.Create(&file).Error; err != nil {
// 		logger.Log.Println("Error inserting file metadata:", err)
// 		return entities.File{}, err
// 	}

//		return file, nil
//	}
func (db *DbConn) InsertFolderDao(fileHeader *multipart.FileHeader, filePath, filePathIds string, folderID, userID uint) (entities.File, error) {
	logger.Log.Println("Inside InsertFolderDao")

	var existingFile entities.File

	// Check if file with the same path already exists
	logger.Log.Println(filePath)
	err := db.DB.Where("path = ?", filePath).First(&existingFile).Error

	if err == nil {
		// File exists -> Update the existing record
		existingFile.Name = fileHeader.Filename
		existingFile.FolderID = &folderID
		existingFile.UserID = userID

		if err := db.DB.Save(&existingFile).Error; err != nil {
			logger.Log.Println("Error updating file metadata:", err)
			return entities.File{}, err
		}
		logger.Log.Println("File updated successfully")
		return existingFile, nil
	} else if err != gorm.ErrRecordNotFound {
		// If there's an actual error (other than "not found"), return it
		logger.Log.Println("Error checking existing file:", err)
		return entities.File{}, err
	}

	// File does not exist -> Insert new record
	newFile := entities.File{
		Name:     fileHeader.Filename,
		Path:     filePath,
		FolderID: &folderID,
		UserID:   userID,
		Pathid:   filePathIds,
	}

	if err := db.DB.Create(&newFile).Error; err != nil {
		logger.Log.Println("Error inserting new file metadata:", err)
		return entities.File{}, err
	}

	logger.Log.Println("File inserted successfully")
	return newFile, nil
}

func (db *DbConn) GetFolderDetails(folderId uint) (entities.Folder, error) {
	logger.Log.Println("Inside GetFolderDetailsDao")

	var folder entities.Folder
	err := db.DB.Where("id = ? AND delete_flag = 0", folderId).
		First(&folder).Error

	if err != nil {
		logger.Log.Println("Error in getting folder details:", err)
		return folder, err
	}

	return folder, nil
}

func (db *DbConn) GetRootFolderNameDao(userID int64) (string, error) {
	logger.Log.Println("Inside GetRootFolderNaeDao")
	var folder entities.Folder
	err := db.DB.First(&folder).Where("user_id = ? AND parent_id = ? AND delete_flag = 0", userID, -1).Error
	if err != nil {
		return "", err // Return error if folder not found
	}
	folderName := folder.Name + "/"
	return folderName, nil
}

func (db *DbConn) GetFullPathDao(folderID uint) (string, string, error) {
	logger.Log.Println("Inside GetFullPathDao")
	logger.Log.Println("FolderID: ", folderID)
	var folder entities.Folder
	err := db.DB.First(&folder, folderID).Error
	if err != nil {
		return "", "", err // Return error if folder not found
	}

	// Store folder names in a slice (to reverse later)
	var pathParts []string
	var pathPartIds []string
	currentFolder := folder

	// Traverse up the hierarchy
	for {
		logger.Log.Println("Inside for loop")
		logger.Log.Println("Current Folder: ", currentFolder.ID)
		pathParts = append([]string{currentFolder.Name}, pathParts...) // Prepend folder name
		str := strconv.FormatUint(uint64(currentFolder.ID), 10)
		pathPartIds = append([]string{str}, pathPartIds...)
		if currentFolder.ParentID == 0 {
			break
		}

		var parentFolder entities.Folder
		logger.Log.Println("ParentID: ", currentFolder.ParentID)
		err := db.DB.First(&parentFolder, currentFolder.ParentID).Error
		if err != nil {
			return "", "", err
		}

		currentFolder = parentFolder
	}

	// Join path parts with "/"
	return "/" + strings.Join(pathParts, "/"), strings.Join(pathPartIds, "/"), nil
}

// func (db *DbConn) GetFoldersDao(userID uint) ([]entities.Folder, error) {
// 	logger.Log.Println("Inside GetFoldersDao")

// 	var folders []entities.Folder
// 	err := db.DB.Where("user_id = ? AND delete_flag = 0", userID).Find(&folders).Error
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Map folders by ID
// 	folderMap := make(map[uint]*entities.Folder)
// 	var roots []entities.Folder

// 	for i := range folders {
// 		folderMap[folders[i].ID] = &folders[i]
// 	}
// 	logger.Log.Println("FolderMap: ", folderMap)
// 	for i := range folders {
// 		folder := &folders[i]
// 		logger.Log.Println("Folder: ", folder)
// 		if folder.ParentID == folder.ID {
// 			roots = append(roots, *folder) // Root folder
// 		} else {
// 			parentFolder := folderMap[folder.ParentID]
// 			logger.Log.Println("Parent Folder: ", parentFolder.ID)
// 			logger.Log.Println("Folder: ", folder.ID)
// 			if parentFolder.ID != folder.ParentID {
// 			// 	logger.Log.Println("Inside if")
// 			parentFolder.Children = append(parentFolder.Children, *folder)
// 			logger.Log.Println("Parent Folder: ", parentFolder.Children)
// 			roots = append(roots, *folder)

// 			}
// 		}
// 	}

// 	return roots, nil

// }

// Function to build folder paths recursively
func buildFolderPaths(folderMap map[uint]*entities.Folder, folder *entities.Folder, path []uint, result *[][]uint) {
	// Add current folder ID to the path
	path = append(path, folder.ID)

	// Find children and recursively process
	isLeaf := true
	for _, f := range folderMap {
		if f.ParentID != 0 && f.ParentID == folder.ID {
			isLeaf = false
			buildFolderPaths(folderMap, f, path, result)
		}
	}

	// If it's a leaf node, store the path as []uint
	if isLeaf {
		pathCopy := make([]uint, len(path))
		copy(pathCopy, path) // Make a copy to avoid modifications in recursion
		*result = append(*result, pathCopy)
	}
}

// Function to generate folder paths
func generateFolderPaths(folders []entities.Folder) [][]uint {
	// Create a map of folders by ID
	folderMap := make(map[uint]*entities.Folder)
	for i := range folders {
		folderMap[folders[i].ID] = &folders[i]
	}

	// Find root folders and generate paths
	var result [][]uint
	for _, folder := range folders {
		if folder.ParentID == 0 {
			buildFolderPaths(folderMap, &folder, []uint{}, &result)
		}
	}

	return result
}

func (db *DbConn) GetFoldersDao(userID uint) ([][]entities.Folder, error) {
	var folders []entities.Folder
	var responseFolder [][]entities.Folder
	// Fetch all folders for the given user
	err := db.DB.Where("user_id = ? AND delete_flag = 0", userID).Find(&folders).Error
	if err != nil {
		return nil, err
	}

	paths := generateFolderPaths(folders)
	for _, p := range paths {
		logger.Log.Println(p)
		var temp []entities.Folder
		err := db.DB.Where("id IN ?", p).Find(&temp).Error
		if err != nil {
			return nil, err
		}
		for i, _ := range temp {
			temp[i].PathIds = p
		}
		responseFolder = append(responseFolder, temp)
	}
	logger.Log.Println("responseFolder: ", responseFolder)

	return responseFolder, nil
}

func (db *DbConn) GetFilesDao(userID uint) ([]entities.File, error) {
	logger.Log.Println("Inside GetFilesDao")

	var files []entities.File
	err := db.DB.Where("user_id = ? AND delete_flag = 0", userID).Find(&files).Error
	if err != nil {
		logger.Log.Println("Error fetching files:", err)
		return nil, err
	}

	return files, nil
}

func (db *DbConn) DeleteFileByUserIdDao(fileID int64, userID int64) error {
	logger.Log.Println("Inside DeleteFileByUserIdDao")

	// Find the file to ensure it exists
	var file entities.File
	if err := db.DB.Where("id = ? AND user_id = ?", fileID, userID).First(&file).Error; err != nil {
		logger.Log.Println("Error finding file:", err)
		return err
	}

	// Update delete_flag to 1
	if err := db.DB.Model(&file).Update("delete_flag", 1).Error; err != nil {
		logger.Log.Println("Error updating delete_flag:", err)
		return err
	}

	logger.Log.Println("File soft deleted successfully")
	return nil
}

func (db *DbConn) DeleteFolderByUserIdDao(folderID int64, userID int64) error {
	logger.Log.Println("Inside DeleteFolderByUserIdDao")

	// Find the file to ensure it exists
	var folder entities.Folder
	if err := db.DB.Where("id = ? AND user_id = ?", folderID, userID).First(&folder).Error; err != nil {
		logger.Log.Println("Error finding file:", err)
		return err
	}

	// Update delete_flag to 1
	if err := db.DB.Model(&folder).Update("delete_flag", 1).Error; err != nil {
		logger.Log.Println("Error updating delete_flag:", err)
		return err
	}

	logger.Log.Println("File soft deleted successfully")
	return nil
}

func (db *DbConn) CreateRootFolderByUserIdModel(userID int, folderName string) (entities.Folder, error) {
	logger.Log.Println("Inside CreateRootFolderByUserIdModel")

	newFolder := entities.Folder{
		Name:     folderName,
		UserID:   uint(userID),
		ParentID: 0,
	}
	if err := db.DB.Create(&newFolder).Error; err != nil {
		logger.Log.Println("Error in inserting user:", err)
		return entities.Folder{}, err
	}

	return newFolder, nil
}

func (db *DbConn) CreateFolderByUserIdDao(reqData entities.CreateFolderByUserIdRequestEntity) (entities.Folder, error) {
	logger.Log.Println("Inside CreateRootFolderByUserIdModel")

	newFolder := entities.Folder{
		Name:     reqData.FolderName,
		UserID:   uint(reqData.Userid),
		ParentID: reqData.FolderPathIds[len(reqData.FolderPathIds)-1],
	}

	logger.Log.Println("New Folder: ", newFolder)
	if err := db.DB.Create(&newFolder).Error; err != nil {
		logger.Log.Println("Error in inserting user:", err)
		return entities.Folder{}, err
	}

	return newFolder, nil
}

func (db *DbConn) GetFolderPathDao(reqData entities.CreateFolderByUserIdRequestEntity) (string, error) {
	logger.Log.Println("Inside CreateRootFolderByUserIdModel")
	var temp []entities.Folder
	err := db.DB.Where("id IN ? AND delete_flag = 0", reqData.FolderPathIds).Find(&temp).Error
	if err != nil {
		return "", err
	}

	var folderPath string
	for i, _ := range temp {
		// if i == len(temp)-1 {
		// 	folderPath = folderPath + temp[i].Name
		// } else {
		// 	folderPath = folderPath + temp[i].Name + "/"
		// }
		folderPath = folderPath + temp[i].Name + "/"
	}

	logger.Log.Println("FolderPath: ", folderPath)

	folderPath = folderPath + reqData.FolderName

	return folderPath, nil
}

func (db *DbConn) CheckIsThereAnyFolderPresentOrNotDao(userId int) (bool, error) {
	logger.Log.Println("Inside CheckIsThereAnyFolderPresentOrNotDao")
	var count int64
	err := db.DB.Model(&entities.Folder{}).Where("user_id = ? AND delete_flag = 0", userId).Count(&count).Error
	if err != nil {
		logger.Log.Println("Error counting folders:", err)
		return false, err
	}

	return count > 0, nil
}
