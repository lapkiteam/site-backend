package auth

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/lapkiteam/site-backend/app/user"
	"github.com/lapkiteam/site-backend/pkg/cookie/previewurl"
	"github.com/lapkiteam/site-backend/pkg/cookie/session"
	"github.com/lapkiteam/site-backend/pkg/database"
	"net/http"
	"strconv"
	"time"
)

func AuthRegister(router *gin.RouterGroup) {
	router.POST("/auth", authLogin)
	router.GET("/auth", authRetrieve)
}

func authLogin(ctx *gin.Context) {
	login := ctx.PostForm("login")
	password := ctx.PostForm("password")

	dbConnection := database.GetDB()
	var user user.UserModel
	dbConnection.First(&user, "login = ?", login)

	var userExists bool
	if user.Login != "" && password == user.Password {
		userExists = true
	} else {
		userExists = false
	}

	if userExists {
		currentTime := strconv.Itoa(int(time.Now().Unix()))

		token := base64.StdEncoding.EncodeToString([]byte(login + ":" + password + currentTime))

		dbConnection := database.GetDB()
		dbConnection.Create(&SessionModel{Token: token})

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

func authRetrieve(ctx *gin.Context) {
	cookie, err := ctx.Cookie(session.CookieName)
	if err != nil {
		ctx.HTML(http.StatusOK, "index.html", gin.H{})
		return
	}

	dbConnection := database.GetDB()
	var dbSession SessionModel
	dbConnection.First(&dbSession, "token = ?", cookie)

	if dbSession.Token != "" {
		ctx.Redirect(http.StatusSeeOther, "/MissingBoar.Docs/")
		return
	}

	ctx.HTML(http.StatusOK, "index.html", gin.H{})
	return
}
