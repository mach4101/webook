package service

import (
	"context"
	"errors"

	"github.com/mach4101/geek_go_camp/webook/internal/domain"
	"github.com/mach4101/geek_go_camp/webook/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var ErrDuplicateEmail = repository.ErrUserDuplicateEmail
var ErrInvalidUserOrPassword = errors.New("账号或密码不对")

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hash)
	return svc.repo.Create(ctx, u)
}

func (svc *UserService) Login(ctx context.Context, Email, Password string) (domain.User, error) {
	// 查找用户
	u, err := svc.repo.FindByEmail(ctx, Email)

	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}

	if err != nil {
		return domain.User{}, err
	}

	// 比较密码
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(Password))

	if err != nil {

		// 接入日志之后需要记录日志
		return domain.User{}, ErrInvalidUserOrPassword
	}

	return u, nil
}
