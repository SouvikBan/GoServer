package main

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB
var err error

type (
	//Users describes a User type
	Users struct {
		ID       uint
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Admin    bool   `json:"admin"`
	}
	//Genre descrbes each genre
	Genre struct {
		ID      uint
		Gname   string `json:"gname"`
		Quizzes []Quiz
	}

	//Quiz describes each quiz type in the app
	Quiz struct {
		ID        uint
		GenreID   uint
		Quizname  string `json:"quizname"`
		Questions []Question
	}

	//Question describes each question type
	Question struct {
		ID       uint
		QuizID   uint
		Qtype    string `json:"qtype"`
		Question string `json:"question"`
		Options  string `json:"options"`
		Answer   string `json:"answer"`
	}
	//Scores describes each user's score type
	Scores struct {
		ID        uint
		UserID    uint
		QuizID    uint
		Attempts  uint64 `json:"attempts"`
		BestScore uint64 `json:"bestscore"`
	}
)

func init() {

	db, err = gorm.Open("sqlite3", "./gorm.db")
	if err != nil {
		panic("failed to connect database")
	}

	db.CreateTable(&Users{}, &Genre{}, &Quiz{}, &Question{}, &Scores{})
	db.AutoMigrate(&Users{}, &Genre{}, &Quiz{}, &Question{}, &Scores{})

}

func main() {
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	router.Use(cors.New(config))

	users := router.Group("/api/users")
	{
		users.POST("/", CreateUser)
		users.GET("/", GetAllUsers)
		users.GET("/:id", GetUser)
		users.DELETE("/:id", DeleteUser)
	}
	genres := router.Group("/api/genres")
	{
		genres.POST("/", CreateGenre)
		genres.GET("/", GetAllGenres)
		genres.DELETE("/:id", DeleteGenre)
	}

	quizzes := router.Group("/api/quizzes")
	{
		quizzes.POST("/", CreateQuiz)
		quizzes.DELETE("/:id", DeleteQuiz)
	}

	questions := router.Group("/api/questions")
	{
		questions.POST("/", CreateQuestion)
		questions.GET("/:id", GetAllQuestions)
		questions.PUT("/:id", EditQuestion)
		questions.DELETE("/:id", DeleteQuestion)
	}
	scores := router.Group("/api/scores")
	{
		scores.POST("/", CreateScores)
		scores.GET("/:id", GetScores)
	}

	router.Run()
}

//CreateUser adds an user to the database
func CreateUser(c *gin.Context) {
	var user Users
	c.BindJSON(&user)

	db.Create(&user)

	c.JSON(http.StatusOK, user)

}

//GetAllUsers gets all the user from the database
func GetAllUsers(c *gin.Context) {
	var user []Users
	if err := db.Find(&user).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(http.StatusOK, user)
	}
}

//GetUser gets an user of a specific id from the database
func GetUser(c *gin.Context) {
	id := c.Params.ByName("id")
	var user Users
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(http.StatusOK, user)
	}
}

//DeleteUser deletes an user with a specific id from the database
func DeleteUser(c *gin.Context) {
	id := c.Params.ByName("id")
	var user Users
	retVal := db.Where("id=?", id).Delete(&user)
	fmt.Println(retVal)

	c.JSON(http.StatusOK, user)

}

//CreateGenre creates a genre from the database
func CreateGenre(c *gin.Context) {
	var genre Genre
	c.BindJSON(&genre)

	db.Create(&genre)

	c.JSON(http.StatusOK, genre)

}

//GetAllGenres gets each genre and its detail
func GetAllGenres(c *gin.Context) {
	var genre []Genre
	if err := db.Find(&genre).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		for i, g := range genre {
			var quizzes []Quiz
			db.Model(&g).Related(&quizzes)
			genre[i].Quizzes = append(genre[i].Quizzes, quizzes...)
		}
		c.JSON(http.StatusOK, genre)
	}
}

//DeleteGenre deletes a genre
func DeleteGenre(c *gin.Context) {
	id := c.Params.ByName("id")
	var genre Genre
	retVal := db.Where("id=?", id).Delete(&genre)
	fmt.Println(retVal)
	c.JSON(http.StatusOK, genre)
}

//CreateQuiz creates a quiz for a genre
func CreateQuiz(c *gin.Context) {
	var quiz Quiz
	c.BindJSON(&quiz)

	db.Create(&quiz)

	c.JSON(http.StatusOK, quiz)
}

//DeleteQuiz deletes quiz of a genre
func DeleteQuiz(c *gin.Context) {
	id := c.Params.ByName("id")
	var quiz Quiz
	retVal := db.Where("id=?", id).Delete(&quiz)
	fmt.Println(retVal)
	c.JSON(http.StatusOK, quiz)
}

//CreateQuestion creates of a particular quiz
func CreateQuestion(c *gin.Context) {
	var question Question
	c.BindJSON(&question)

	db.Create(&question)

	c.JSON(http.StatusOK, question)
}

//GetAllQuestions gets a question of a particular quiz
func GetAllQuestions(c *gin.Context) {
	id := c.Params.ByName("id")
	var quiz Quiz
	if err := db.Where("id = ?", id).First(&quiz).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		var questions []Question
		db.Model(&quiz).Related(&questions)
		quiz.Questions = append(quiz.Questions, questions...)
		c.JSON(http.StatusOK, quiz.Questions)
	}
}

//EditQuestion edits a question of a particular quiz
func EditQuestion(c *gin.Context) {

}

//DeleteQuestion deletes a question of a particular quiz
func DeleteQuestion(c *gin.Context) {
	id := c.Params.ByName("id")
	var question Question
	retVal := db.Where("id=?", id).Delete(&question)
	fmt.Println(retVal)
	c.JSON(http.StatusOK, question)
}

//GetScores gives you the score for each quiz of a user
func GetScores(c *gin.Context) {
	id := c.Params.ByName("id")
	var scores []Scores
	db.Where("UserID=?", id).Find(&scores)
	c.JSON(http.StatusOK, scores)
}

//CreateScores keeps track of each score for every user
func CreateScores(c *gin.Context) {
	var scores Scores
	c.BindJSON(&scores)

	db.Create(&scores)

	c.JSON(http.StatusOK, "Done")
}
