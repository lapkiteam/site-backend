package main

import (
	"encoding/base64"
	authconfig "github.com/lapkiteam/docs-viewer/pkg/config/auth"
	"github.com/lapkiteam/docs-viewer/pkg/cookie/previewurl"
	"github.com/lapkiteam/docs-viewer/pkg/cookie/session"
	"github.com/lapkiteam/docs-viewer/pkg/middleware/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

var tokens []string

func main() {
	router := gin.Default()
	router.Use()

	router.LoadHTMLFiles("auth/index.html")
	router.Static("/css", "auth/css")

	router.GET("/auth", func(ctx *gin.Context) {
		cookie, err := ctx.Cookie(session.CookieName)
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
		ctx.Redirect(http.StatusSeeOther, "/MissingBoar.Docs/")
		return
	})

	authorized := router.Group("/", auth.Auth(&tokens))
	{
		authorized.StaticFS("/MissingBoar.Docs/", gin.Dir("site/docs/MissingBoar", false))
	}

	router.Run(":8080")
}

func postAuthEndpoint(ctx *gin.Context) {
	login := ctx.PostForm("login")
	password := ctx.PostForm("password")

	userExists := authconfig.CheckLogin(login, password)

	if userExists {
		token := base64.StdEncoding.EncodeToString([]byte(login + ":" + password))
		tokens = append(tokens, token)
		session.SetCookie(ctx, token, "/")

		ctx.Request.Method = "GET"
		redirectPath, err := ctx.Cookie(previewurl.CookieName)
		if err != nil {
			ctx.Redirect(http.StatusSeeOther, "/MissingBoar.Docs/")
		}

		ctx.Redirect(http.StatusSeeOther, redirectPath)
		return
	}

	ctx.Request.Method = "GET"
	ctx.Redirect(http.StatusSeeOther, "/auth")
	return
}
