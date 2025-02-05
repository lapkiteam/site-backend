package main

import (
	"github.com/gin-gonic/gin"
	"github.com/lapkiteam/site-backend/app/auth"
	"github.com/lapkiteam/site-backend/app/user"
	"github.com/lapkiteam/site-backend/pkg/database"
	"gorm.io/gorm"
	"net/http"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&auth.SessionModel{})
	db.AutoMigrate(&user.UserModel{})
}

func main() {
	// Create DB connection and make migrate
	dbConnection := database.Init()
	Migrate(dbConnection)

	//Create Router and groups
	router := gin.Default()

	authGroup := router.Group("/")
	auth.AuthRegister(authGroup)
	authGroup.Use(auth.Auth())

	router.LoadHTMLFiles("site/auth/index.html")
	router.Static("/css", "site/auth/css")

	router.GET("/", func(ctx *gin.Context) {
		ctx.Request.Method = "GET"
		ctx.Redirect(http.StatusSeeOther, "/MissingBoar.Docs/")
		return
	})

	authorized := router.Group("/", auth.Auth())
	{
		authorized.StaticFS("/MissingBoar.Docs/", gin.Dir("site/docs/MissingBoar", false))
	}

	router.Run(":8080")
}
