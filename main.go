package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	tokens               []string
	sessionCookieName    = "sessionId"
	previewUrlCookieName = "previewUrl"
	serverName = https://www.lapkiteam.fun
	//serverName = "localhost"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookie, err := ctx.Cookie(sessionCookieName)
		if err != nil {
			ctx.Request.Method = "GET"
			ctx.Redirect(http.StatusSeeOther, "/auth")
			return
		}

		for _, token := range tokens {
			if cookie == token {
				ctx.Next()
				return
			}
		}

		ctx.Request.Method = "GET"
		ctx.Redirect(http.StatusSeeOther, "/auth")
		ctx.Next()
		return
	}
}

func main() {
	router := gin.Default()
	router.Use()

	router.LoadHTMLFiles("auth/index.html")
	router.Static("/css", "auth/css")

	router.GET("/auth", func(ctx *gin.Context) {
		cookie, err := ctx.Cookie(sessionCookieName)
		if err != nil {
			ctx.HTML(http.StatusOK, "index.html", gin.H{})
			return
		}

		for _, token := range tokens {
			if cookie == token {
				ctx.Redirect(http.StatusSeeOther, "/MissingBoar.Docs/")
				return
			}
		}

		ctx.HTML(http.StatusOK, "index.html", gin.H{})
		return

	})
	router.POST("/auth", postAuthEndpoint)
	router.GET("/", func(ctx *gin.Context) {
		ctx.Request.Method = "GET"
		ctx.Redirect(http.StatusSeeOther, "/auth")
		return
	})

	authorized := router.Group("/", Auth())
	{
		authorized.StaticFS("/MissingBoar.Docs/", gin.Dir("site", false))
	}

	router.Run(":8080")
}

func getAuthDataFromEnv() map[string]string {
	envData := map[string]string{}

	file, err := os.Open(".env")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		readData := strings.Split(scanner.Text(), "=")

		envData[readData[0]] = readData[1]
	}

	return envData
}

func postAuthEndpoint(ctx *gin.Context) {
	login := ctx.PostForm("login")
	password := ctx.PostForm("password")

	dataFromEnv := getAuthDataFromEnv()

	for envLogin, envPassword := range dataFromEnv {
		if (login == envLogin) && (password == envPassword) {
			token := base64.StdEncoding.EncodeToString([]byte(login + ":" + password))
			tokens = append(tokens, token)
			ctx.SetCookie(sessionCookieName, token, 60*60, "/", serverName, false, true)
			ctx.Request.Method = "GET"
			ctx.Redirect(http.StatusSeeOther, "/MissingBoar.Docs/")
			return
		}
	}

	ctx.Request.Method = "GET"
	ctx.Redirect(http.StatusSeeOther, "/auth")
	return
}
