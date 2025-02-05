package previewurl

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lapkiteam/site-backend/pkg/config"
)

const (
	CookieName = "previewUrl"
	maxAge     = 60 * 60 // 1 hour
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
