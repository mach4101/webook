package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"

	"github.com/mach4101/geek_go_camp/webook/internal/web"
)

// JWT登录校验
type LoginJWTMiddlewareBuilder struct {
	paths []string
}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{}
}

func (l *LoginJWTMiddlewareBuilder) IgnorePaths(path string) *LoginJWTMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 不需要登陆校验的：
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}

		// JWT校验
		tokenHeader := ctx.GetHeader("Authorization")

		if tokenHeader == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		segs := strings.SplitN(tokenHeader, " ", 2)

		if len(segs) != 2 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenStr := segs[1]

		claims := &web.UserClaims{}
		// 从tokenstr中把数据解析到claims中, 还原数据结构
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("nUCUFGagbcXzkDJ33spmZ6CyW8zNaFu3"), nil
		})
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if token == nil || !token.Valid || claims.Uid == 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 若在后续的操作中需要使用到token中携带的字段，俺么可以使用set
		ctx.Set("claims", claims)
	}
}
