package repository

import (
	"context"

	"github.com/mach4101/geek_go_camp/webook/internal/domain"
	"github.com/mach4101/geek_go_camp/webook/internal/repository/cache"
	"github.com/mach4101/geek_go_camp/webook/internal/repository/dao"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFfound
)

type UserRepository struct {
	dao   *dao.UserDAO
	cache *cache.UserCache
}

func NewUserRepository(dao *dao.UserDAO, c *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   dao,
		cache: c,
	}
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})

	// 操作缓存
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       u.Id,
		Email:    u.Password,
		Password: u.Password,
	}, nil
}

func (r *UserRepository) UpdateByEmail(ctx context.Context, email, password string) error {
	err := r.dao.UpdateByEmail(ctx, email, password)
	return err
}

func (r *UserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	// 先从cache找，没有找到再查询数据库
	u, err := r.cache.Get(ctx, id)

	// 有数据
	if err == nil {
		return u, nil
	}

	// 数据库找
	ue, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	u = domain.User{
		Id:       ue.Id,
		Email:    ue.Email,
		Password: ue.Password,
	}

	// 添加到缓存
	go func() {
		err = r.cache.Set(ctx, u)
		if err != nil {
			// 监控
		}
	}()

	return u, err
}
