package session

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lapkiteam/site-backend/pkg/config"
)

const (
	CookieName = "sessionId"
	maxAge     = 60 * 60 * 24 * 30 // 1 month
)

func SetCookie(ctx *gin.Context, value string, path string) {
	if path == "" {
		path = "/"
	}

	appUrl, err := config.GetAppUrl()
	if err != nil {
		fmt.Println(err)
	}

	ctx.SetCookie(CookieName, value, maxAge, path, appUrl, false, true)
}
