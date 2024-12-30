package user

import (
	"context"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/merynayr/auth/internal/client/db"
	"github.com/merynayr/auth/internal/model"
	"github.com/merynayr/auth/internal/repository"
	"github.com/merynayr/auth/internal/repository/user/converter"
	modelRepo "github.com/merynayr/auth/internal/repository/user/model"
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

	query, args, err := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, emailColumn, passwordColumn, roleColumn, createdAtColumn, updatedAtColumn).
		Values(user.Name, user.Email, user.Password, user.Role, time.Now(), user.UpdatedAt).
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

func (r *repo) GetUser(ctx context.Context, userID int64) (*model.User, error) {
	op := "GetUser"
	log.Printf("[%s] request data | id: %v", op, userID)

	query, args, err := sq.Select(idColumn, nameColumn, emailColumn, roleColumn, createdAtColumn, updatedAtColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: userID}).
		ToSql()

	if err != nil {
		log.Printf("%s: failed to create builder: %v", op, err)
		return nil, err
	}

	q := db.Query{
		Name:     "user_repository.GetById",
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
