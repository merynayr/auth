package user

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/merynayr/auth/internal/client/db"
	"github.com/merynayr/auth/internal/model"
	"github.com/merynayr/auth/internal/repository"
	"github.com/merynayr/auth/internal/repository/user/converter"
	modelRepo "github.com/merynayr/auth/internal/repository/user/model"
	"github.com/merynayr/auth/internal/utils/hash"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	tableName = "users"

	idColumn        = "id"
	nameColumn      = "username"
	emailColumn     = "email"
	passwordColumn  = "password"
	roleColumn      = "role"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

// Структура репо с клиентом базы данных (интерфейсом)
type repo struct {
	db db.Client
}

// NewRepository возвращает новый объект репо слоя
func NewRepository(db db.Client) repository.UserRepository {
	return &repo{db: db}
}

func (r *repo) CreateUser(ctx context.Context, user *model.User) (int64, error) {
	op := "CreateUser"
	log.Printf("%s %v", op, user.ID)

	exist, err := r.IsNameExist(ctx, user.Name)
	if err != nil {
		return 0, err
	}
	if exist {
		return 0, fmt.Errorf("user with name %s already exists", user.Name)
	}

	exist, err = r.IsEmailExist(ctx, user.Email)
	if err != nil {
		return 0, err
	}
	if exist {
		return 0, fmt.Errorf("user with email %s already exists", user.Email)
	}

	passHash, err := hash.EncryptPassword(user.Password)
	if err != nil {
		return 0, err
	}

	query, args, err := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, emailColumn, passwordColumn, roleColumn, createdAtColumn, updatedAtColumn).
		Values(user.Name, user.Email, passHash, user.Role, time.Now(), user.UpdatedAt).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		log.Printf("%s: failed to create builder: %v", op, err)
		return 0, err
	}

	q := db.Query{
		Name:     "user_repository.CreateUser",
		QueryRaw: query,
	}

	var userID int64
	err = r.db.DB().ScanOneContext(ctx, &userID, q, args...)
	if err != nil {
		log.Printf("%s: failed to insert user: %v", op, err)
		return 0, err
	}

	log.Printf("%s: inserted user with id: %d", op, userID)
	return userID, nil
}

func (r *repo) GetUserByID(ctx context.Context, userID int64) (*model.User, error) {
	op := "GetUser"
	log.Printf("[%s] request data | id: %v", op, userID)

	exist, err := r.IsExistByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("user with id %d doesn't exist", userID)
	}

	query, args, err := sq.Select(idColumn, nameColumn, emailColumn, roleColumn, createdAtColumn, updatedAtColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: userID}).
		ToSql()

	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "user_repository.GetUser",
		QueryRaw: query,
	}

	var user modelRepo.User
	err = r.db.DB().ScanOneContext(ctx, &user, q, args...)
	if err != nil {
		log.Printf("%s: failed to select user: %v", op, err)
		return nil, err
	}

	log.Printf("%s: selected user %d", op, userID)
	return converter.ToUserFromRepo(&user), nil
}

// GetUserByName получает из БД информацию пользователя
func (r *repo) GetUserByName(ctx context.Context, name string) (*model.UserInfo, error) {
	exist, err := r.IsNameExist(ctx, name)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, fmt.Errorf("user with name %s doesn't exist", name)
	}

	query, args, err := sq.Select(nameColumn, passwordColumn, roleColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{nameColumn: name}).
		Limit(1).ToSql()

	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "user_repository.getUserByName",
		QueryRaw: query,
	}

	var user modelRepo.UserInfo
	err = r.db.DB().ScanOneContext(ctx, &user, q, args...)
	if err != nil {
		return nil, err
	}

	return converter.ToUserInfoFromRepo(&user), nil
}

// UpdateUser обновляет данные пользователя по id
func (r *repo) UpdateUser(ctx context.Context, user *model.UserUpdate) (*emptypb.Empty, error) {
	op := "UpdateUser"
	log.Printf("[%s] request data | id: %v", op, user.ID)

	exist, err := r.IsExistByID(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("user with id %d doesn't exist", user.ID)
	}

	builderUpdate := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Set(updatedAtColumn, time.Now())

	if user.Name != nil {
		builderUpdate = builderUpdate.Set(nameColumn, *user.Name)
	}

	if user.Email != nil {
		builderUpdate = builderUpdate.Set(emailColumn, *user.Email)
	}

	builderUpdate = builderUpdate.Where(sq.Eq{idColumn: user.ID})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		log.Printf("%s: failed to create builder: %v", op, err)
		return nil, err
	}

	q := db.Query{
		Name:     "user_repository.UpdateUser",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		log.Printf("%s: failed to update user: %v", op, err)
		return nil, err
	}

	log.Printf("%s: updated user %d", op, user.ID)
	return &emptypb.Empty{}, nil
}

// DeleteUser удаляет пользователя по id
func (r *repo) DeleteUser(ctx context.Context, userID int64) (*emptypb.Empty, error) {
	exist, err := r.IsExistByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("user with id %d doesn't exist", userID)
	}

	query, args, err := sq.Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: userID}).
		ToSql()

	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "user_repository.DeleteUser",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// IsExistById проверяет, существует ли в БД пользователь с указанным ID
func (r *repo) IsExistByID(ctx context.Context, userID int64) (bool, error) {
	op := "IsExistByID"
	log.Printf("[%s] request data | id: %v", op, userID)

	query, args, err := sq.Select("1").
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: userID}).
		ToSql()

	if err != nil {
		return false, err
	}

	q := db.Query{
		Name:     "user_repository.IsExistByID",
		QueryRaw: query,
	}

	var user int
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&user)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// IsEmailExist проверяет, существует ли в БД указанный email
func (r *repo) IsEmailExist(ctx context.Context, email string) (bool, error) {
	query, args, err := sq.Select("1").
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{emailColumn: email}).
		Limit(1).ToSql()

	if err != nil {
		return false, err
	}

	q := db.Query{
		Name:     "user_repository.IsEmailExist",
		QueryRaw: query,
	}

	var one int

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&one)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Println("No rows", one)
			return false, nil
		}
		log.Println(err)

		return false, err
	}

	return true, nil
}

// IsNameExist проверяет, существует ли в БД указанный name
func (r *repo) IsNameExist(ctx context.Context, name string) (bool, error) {
	query, args, err := sq.Select("1").
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{nameColumn: name}).
		Limit(1).ToSql()

	if err != nil {
		return false, err
	}

	q := db.Query{
		Name:     "user_repository.IsNameExist",
		QueryRaw: query,
	}

	var one int

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&one)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
