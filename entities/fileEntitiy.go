package entities

import (
	"encoding/json"
	"io"
)

type FolderResponseEntity struct {
	Sucess  bool       `json:"success"`
	Message string     `json:"message"`
	Folders [][]Folder `json:"folders"`
}

type FilesResponseEntity struct {
	Sucess  bool   `json:"success"`
	Message string `json:"message"`
	Files   []File `json:"files"`
}

type DeleteFileByUserIdResponseEntity struct {
	Sucess  bool   `json:"success"`
	Message string `json:"message"`
}

type CreateFolderByUserIdResponseEntity struct {
	Sucess   bool     `json:"success"`
	Message  string   `json:"message"`
	Response []Folder `json:"response"`
}

type CreateFolderByUserIdRequestEntity struct {
	FolderPathIds []uint `json:"filePathIds"`
	FolderName    string `json:"folderName"`
	Userid        int    `json:"userid"`
}

func (w *CreateFolderByUserIdRequestEntity) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(w)
}
