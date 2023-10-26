package main

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/mach4101/geek_go_camp/webook/internal/repository"
	"github.com/mach4101/geek_go_camp/webook/internal/repository/dao"
	"github.com/mach4101/geek_go_camp/webook/internal/service"
	"github.com/mach4101/geek_go_camp/webook/internal/web"
	"github.com/mach4101/geek_go_camp/webook/internal/web/middleware"
)

func main() {
	db := initDB()
	server := initWebServer()

	u := initUser(db)
	u.RegisterRoutes(server)
	server.Run(":8080")
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"))
	// 只有在初始化中，panic
	if err != nil {
		panic(err)
	}

	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}

func initUser(db *gorm.DB) *web.UserHandler {
	ud := dao.NewUserDAO(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)
	return u
}

func initWebServer() *gin.Engine {
	server := gin.Default()

	server.Use(cors.New(cors.Config{
		// AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "domain.name")
		},

		MaxAge: 12 * time.Hour,
	}))

	// 第一个参数是sutentication key, 第二个是encryption key，最好是32位或者64位
	store := memstore.NewStore([]byte("nUCUFGagbcXzkDJ33spmZ6CyW8zNaFu3"), []byte("wm67pcvktHdVpiHbxqV5W7kfJssuQ0Ae"))
	server.Use(sessions.Sessions("mysession", store))

	// 增加登陆校验
	server.Use(middleware.NewLoginMiddlewareBuilder().IgnorePaths("/users/login").IgnorePaths("/users/signup").Build())

	return server
}
