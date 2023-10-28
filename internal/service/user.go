package service

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/mach4101/geek_go_camp/webook/internal/domain"
	"github.com/mach4101/geek_go_camp/webook/internal/repository"
)

var (
	ErrDuplicateEmail        = repository.ErrUserDuplicateEmail
	ErrInvalidUserOrPassword = errors.New("账号或密码不对")
)

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

	// 比较密码, 第一个是hash，第二个是密码
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(Password))

	if err != nil {
		// 接入日志之后需要记录日志
		return domain.User{}, ErrInvalidUserOrPassword
	}

	return u, nil
}

func (svc *UserService) Edit(ctx context.Context, oldUser domain.User, password string) error {
	// 第一部分和登录一样，主要是查找用户，比较密码，判断该用户是否存在
	u, err := svc.repo.FindByEmail(ctx, oldUser.Email)
	if err == repository.ErrUserNotFound {
		return ErrInvalidUserOrPassword
	}

	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(oldUser.Password))

	if err != nil {
		return ErrInvalidUserOrPassword
	}

	hashedPasswd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = svc.repo.UpdateByEmail(ctx, oldUser.Email, string(hashedPasswd))

	if err != nil {
		return errors.New("更新出问题")
	}
	return err
}
