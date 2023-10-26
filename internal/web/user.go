package web

import (
	"fmt"
	"net/http"

	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/mach4101/geek_go_camp/webook/internal/domain"
	"github.com/mach4101/geek_go_camp/webook/internal/service"
)

// 定义和user相关的路由
type UserHandler struct {
	svc         *service.UserService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	const (
		emailRegexPattern = `^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`
		// 以字母开头，长度在6~18之间，只能包含字母、数字和下划线
		passwordRegexPattern = `^[a-zA-Z]\w{5,17}$`
	)

	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)

	return &UserHandler{
		svc:         svc,
		emailExp:    emailExp,
		passwordExp: passwordExp,
	}
}

// 路由注册
func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	// server.POST("/users/signup" u.SignUp)
	// server.POST("/users/login", u.Login)
	// server.POST("/users/edit", u.Edit)
	// server.GET("/users/profile", u.Profile)

	// 分组路由, 省略前缀
	ug := server.Group("/users")
	ug.GET("/profile", u.Profile)
	ug.POST("/edit", u.Edit)
	ug.POST("/login", u.Login)
	ug.POST("/signup", u.SignUp)
	ug.POST("/logout", u.Logout)
}

// 各handler的具体实现，类似controller
func (u *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		ConfirmPassword string `json:"ConfirmPassword"`
		Password        string `json:"password"`
	}

	var req SignUpReq

	// bind方法会根据Content-Type解析数据到req里头
	// 若解析错误，就会写回一个400错误码
	if err := ctx.Bind(&req); err != nil {
		return
	}

	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "两次输入的密码不一致")
		return
	}

	ok, err := u.emailExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	if !ok {
		ctx.String(http.StatusOK, "邮箱格式不对")
		return
	}

	ok, err = u.passwordExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	if !ok {
		ctx.String(http.StatusOK, "字母开头，长度在6~18之间，只能包含字母、数字和下划线")
		return
	}

	err = u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	if err == service.ErrDuplicateEmail {
		ctx.String(http.StatusOK, "邮箱冲突")
		return
	}

	if err != nil {
		ctx.String(http.StatusOK, "service error")
		return
	}

	ctx.String(http.StatusOK, "ok")
	fmt.Printf("%+v\n", req)

	// 数据库操作
}

func (u *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	user, err := u.svc.Login(ctx, req.Email, req.Password)

	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "用户名或者密码不对")
		return
	}

	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	// 登陆ok，设置和session相关的东西
	sess := sessions.Default(ctx)
	// 可以随便设置值
	sess.Set("userId", user.Id)
	sess.Options(sessions.Options{
		// Secure:   true,
		// HttpOnly: true,
		MaxAge: 60, // 设置一分钟过期
	})
	sess.Save()
	ctx.String(http.StatusOK, "登陆OK")
	return
}

func (u *UserHandler) Logout(ctx *gin.Context) {
	// 登陆ok，设置和session相关的东西
	sess := sessions.Default(ctx)
	// 可以随便设置值

	sess.Options(sessions.Options{
		MaxAge: -1,
	})
	sess.Save()
	ctx.String(http.StatusOK, "注销OK")
	return
}

func (u *UserHandler) Edit(ctx *gin.Context) {
}

func (u *UserHandler) Profile(ctx *gin.Context) {
}
