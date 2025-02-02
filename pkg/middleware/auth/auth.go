package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/lapkiteam/docs-viewer/pkg/cookie/previewurl"
	"github.com/lapkiteam/docs-viewer/pkg/cookie/session"
	"net/http"
)

func Auth(tokens *[]string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookie, err := ctx.Cookie(session.CookieName)
		if err != nil {
			previewurl.SetCookie(ctx, ctx.Request.URL.Path, "/")
			ctx.Request.Method = "GET"
			ctx.Redirect(http.StatusSeeOther, "/auth")
			return
		}

		for _, token := range *tokens {
			if cookie == token {
				ctx.Next()
				return
			}
		}

		previewurl.SetCookie(ctx, ctx.Request.URL.Path, "/")
		ctx.Request.Method = "GET"
		ctx.Redirect(http.StatusSeeOther, "/auth")
		ctx.Next()
		return
	}
}
