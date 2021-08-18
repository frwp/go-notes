package main

import (
	"fmt"
	"log"
	"os"

	"github.com/RianWardanaPutra/notes-v1/controller"
	"github.com/RianWardanaPutra/notes-v1/model"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

func ConnectDb() (*gorm.DB, error) {
	er := godotenv.Load()
	if er != nil {
		log.Fatal("Error loading .env file")
	}

	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPasswd := os.Getenv("DB_PASSWD")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Jakarta", dbHost, dbUser, dbPasswd, dbName)
	db, err := gorm.Open("postgres", dsn)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	db.AutoMigrate(model.User{}, model.Note{})
	db.Model(&model.Note{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "CASCADE")

	return db, nil
}

func main() {
	fmt.Println("Hello")

	// gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	db, err := ConnectDb()

	if err != nil {
		log.Fatal(err)
		return
	}

	c := controller.NewController(db)

	v1 := r.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.GET("/", c.ListUsers)
			users.POST("/", c.AddUser)
			users.GET("/:id", c.FindUserById)
			users.PUT("/:id", c.UpdateUserById)
			users.DELETE("/:id", c.DeleteUserById)
			users.GET("/:id/notes", c.FindNotesByUserId)
			users.POST("/:id/new", c.AddNote)
			users.PUT("/:id/notes/:note_id", c.UpdateNoteByIdAndUID)
			users.DELETE("/:id/notes/:note_id", c.DeleteNote)
		}
		notes := v1.Group("/notes")
		{
			notes.GET("/", c.ListNotes)
			notes.GET("/:note_id", c.FindNoteById)
		}
	}

	r.GET("/", func(c *gin.Context) {
		fmt.Println("landing hit")
		c.JSON(200, gin.H{
			"message": "Welcome",
		})
	})

	r.Run("127.0.0.1:8080")
}
