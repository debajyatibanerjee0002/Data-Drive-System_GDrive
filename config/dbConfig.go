package config

import (
	"fmt"
	"log"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	dbconn *gorm.DB
	once   sync.Once
)

func ConnectMySqlDbValue() (string, string, string, string, string) {
	DBDRIVER := viper.GetString("DBDRIVER")
	DBUSER := viper.GetString("DBUSER")
	DBPASWORD := viper.GetString("DBPASWORD")
	DBURL := viper.GetString("DBURL")
	DBNAME := viper.GetString("DBNAME")

	return DBDRIVER, DBUSER, DBPASWORD, DBURL, DBNAME

}

// func ConnectMySqlDbSingleton() (*sql.DB, error) {
// 	DBDRIVER, DBUSER, DBPASWORD, DBURL, DBNAME := ConnectMySqlDbValue()

// 	logger.Log.Println("DBDRIVER, DBUSER, DBPASWORD, DBURL, DBNAME: ", DBDRIVER, DBUSER, DBPASWORD, DBURL, DBNAME)

// 	dbDriver := DBDRIVER    // Database Driver Name
// 	dbUser := DBUSER        // Database Username
// 	dbPassword := DBPASWORD // Database  Password
// 	dbUrl := DBURL          // Database ip/host with port
// 	dbName := DBNAME        // Database Name
// 	if dbconn == nil {
// 		d, err := sql.Open(dbDriver, dbUser+":"+dbPassword+"@"+dbUrl+"/"+dbName)
// 		if err != nil {
// 			// panic(err.Error())
// 			return nil, err
// 		}
// 		dbconn = d
// 	}
// 	return dbconn, nil
// }

func ConnectMySqlDbSingleton() (*gorm.DB, error) {
	var err error

	once.Do(func() { // Ensures the connection is initialized only once
		_, DBUSER, DBPASSWORD, DBHOST, DBNAME := ConnectMySqlDbValue()

		// GORM connection string format: "user:password@tcp(host:port)/dbname"
		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			DBUSER, DBPASSWORD, DBHOST, DBNAME)

		dbconn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal("Failed to connect to MySQL:", err)
		}

		log.Println("Database connected successfully!")
	})

	return dbconn, err
}
