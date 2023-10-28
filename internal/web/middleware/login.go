package middleware

import (
	"encoding/gob"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type LoginMiddlewareBuilder struct {
	paths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) IgnorePaths(path string) *LoginMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	gob.Register(time.Now())

	return func(ctx *gin.Context) {
		// 不需要登陆校验的：
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}

		// 取出session
		sess := sessions.Default(ctx)
		// 是否已经登录
		id := sess.Get("userId")

		if id == nil {
			// 没有登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 若已经登录，啥刷新超时时间

		updateTime := sess.Get("update_time")
		sess.Set("userId", id)
		sess.Options(sessions.Options{
			MaxAge: 60,
		})

		now := time.Now()

		if updateTime == nil {
			// 还没有刷新
			sess.Set("update_time", now)
			if err := sess.Save(); err != nil {
				panic(err)
			}
		}

		// 之前存过update_time
		updateTimeVal, _ := updateTime.(time.Time)

		// 被改过
		// 如果距离上次刷新时间超过了十秒，那么刷新登录状态

		if now.Sub(updateTimeVal) > time.Second*10 {
			sess.Set("update_time", now)
			sess.Save()
			return
		}
	}
}
