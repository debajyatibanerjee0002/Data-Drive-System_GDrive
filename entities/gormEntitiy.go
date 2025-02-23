package entities

import "time"

// User represents the users table in the database
type User struct {
	ID            uint      `gorm:"primaryKey;autoIncrement"`
	Username      string    `gorm:"size:255;unique;not null"`
	Email         string    `gorm:"size:255;unique;not null"`
	Password_hash string    `gorm:"type:text;not null"`
	ActiveFlag    int       `gorm:"default:1"`
	DeleteFlag    int       `gorm:"default:0"`
	CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

type Folder struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	Name       string    `gorm:"size:255;not null"`
	ParentID   uint      `gorm:"default:0"`
	UserID     uint      `gorm:"not null"`
	ActiveFlag int       `gorm:"default:1"`
	DeleteFlag int       `gorm:"default:0"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	PathIds    []uint    `gorm:"-" json:"path_ids,omitempty"`
	Children   []Folder  `gorm:"-" json:"children,omitempty"`
}

type File struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	Name       string    `gorm:"size:255;not null"`
	Path       string    `gorm:"not null"`
	FolderID   *uint     `gorm:"default:null"` // Nullable foreign key
	UserID     uint      `gorm:"not null"`
	Pathid     string    `gorm:"not null"`
	ActiveFlag int       `gorm:"default:1"`
	DeleteFlag int       `gorm:"default:0"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
