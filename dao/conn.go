package dao

import (
	"gorm.io/gorm"
)

// DbConn is used for initialized connection
type DbConn struct {
	DB *gorm.DB
}
