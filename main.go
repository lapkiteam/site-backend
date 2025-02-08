package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lapkiteam/site-backend/app/auth"
	"github.com/lapkiteam/site-backend/app/user"
	"github.com/lapkiteam/site-backend/pkg/config"
	"github.com/lapkiteam/site-backend/pkg/database"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"os/exec"
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
		authorized.StaticFS("/MissingBoar/", gin.Dir("site/games/MissingBoar", false))
	}

	router.POST("/MissingBoar.Docs/", func(ctx *gin.Context) {
		accessToken := ctx.Request.Header.Get("Access-Token")
		configUploadToken, _ := config.GetUploadToken()

		if accessToken != configUploadToken {
			ctx.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}

		//Recieve File
		file, _ := ctx.FormFile("file")

		ctx.SaveUploadedFile(file, file.Filename)

		//Remove old docs
		files, err := os.ReadDir("site/docs/MissingBoar/")
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range files {
			if file.IsDir() {
				os.RemoveAll("site/docs/MissingBoar/" + file.Name())
			}

			if file.Name() != ".gitignore" && file.IsDir() != true {
				os.Remove("site/docs/MissingBoar/" + file.Name())
			}
		}

		//Unpack Tar
		pwd, _ := os.Getwd()

		cmd := exec.Command("tar", "xf", pwd+"/"+file.Filename, "-C", pwd+"/site/docs/MissingBoar/")
		stdoutStderr, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", stdoutStderr)

		//Delete tar after unpack
		e := os.Remove(file.Filename)
		if e != nil {
			log.Fatal(e)
		}

		//Return response
		ctx.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	})

	router.POST("/MissingBoar/", func(ctx *gin.Context) {
		accessToken := ctx.Request.Header.Get("Access-Token")
		configUploadToken, _ := config.GetUploadToken()

		if accessToken != configUploadToken {
			ctx.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}

		file, _ := ctx.FormFile("file")

		ctx.SaveUploadedFile(file, file.Filename)

		files, err := os.ReadDir("site/games/MissingBoar/")
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range files {
			if file.IsDir() {
				os.RemoveAll("site/games/MissingBoar/" + file.Name())
			}

			if file.Name() != ".gitignore" && file.IsDir() != true {
				os.Remove("site/games/MissingBoar/" + file.Name())
			}
		}

		pwd, _ := os.Getwd()

		cmd := exec.Command("tar", "xf", pwd+"/"+file.Filename, "-C", pwd+"/site/games/MissingBoar/")
		stdoutStderr, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", stdoutStderr)

		e := os.Remove(file.Filename)
		if e != nil {
			log.Fatal(e)
		}

		ctx.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	})

	router.Run(":8080")
}
