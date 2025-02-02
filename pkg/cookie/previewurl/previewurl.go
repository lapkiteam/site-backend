package previewurl

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lapkiteam/docs-viewer/pkg/config/env"
)

const (
	CookieName = "previewUrl"
	maxAge     = 60 * 60 // 1 hour
)

func SetCookie(ctx *gin.Context, value string, path string) {
	if path == "" {
		path = "/"
	}

	appUrl, err := env.GetAppUrl()
	if err != nil {
		fmt.Println(err)
	}

	ctx.SetCookie(CookieName, value, maxAge, path, appUrl, false, true)
}
