package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/lapkiteam/site-backend/pkg/cookie/previewurl"
	"github.com/lapkiteam/site-backend/pkg/cookie/session"
	"github.com/lapkiteam/site-backend/pkg/database"
	"net/http"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookie, err := ctx.Cookie(session.CookieName)
		if err != nil {
			previewurl.SetCookie(ctx, ctx.Request.URL.Path, "/")
			ctx.Request.Method = "GET"
			ctx.Redirect(http.StatusSeeOther, "/auth")
			return
		}

		dbConnection := database.GetDB()
		var dbSession SessionModel
		dbConnection.First(&dbSession, "token = ?", cookie)

		if dbSession.Token != "" {
			ctx.Next()
			return
		}

		previewurl.SetCookie(ctx, ctx.Request.URL.Path, "/")
		ctx.Request.Method = "GET"
		ctx.Redirect(http.StatusSeeOther, "/auth")
		ctx.Next()
		return
	}
}
