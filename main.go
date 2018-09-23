package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB
var err error

type (
	//Users describes a User type
	Users struct {
		ID        uint
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
	}
)

func init() {
	db, err := gorm.Open("sqlite3", "./gorm.db")
	if err != nil {
		panic("failed to connect database")
	}
	db.CreateTable(&Users{})

	defer db.Close()

	db.AutoMigrate(&Users{})
}

func main() {
	router := gin.Default()

	users := router.Group("/api/users")
	{
		users.POST("/", CreateUser)
		users.GET("/", GetAllUsers)
	}
	router.Run()
}

func CreateUser(c *gin.Context) {
	var user Users
	c.BindJSON(&user)

	db.Create(&user)

	c.JSON(http.StatusOK, user)

}

func GetAllUsers(c *gin.Context) {
	var user []Users
	if err := db.Find(&user).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(http.StatusOK, user)
	}
}
