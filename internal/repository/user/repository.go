package user

import (
	"context"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
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
	db *pgxpool.Pool
}

// NewRepository возвращает новый объект репо слоя
func NewRepository(db *pgxpool.Pool) repository.UserRepository {
	return &repo{db: db}
}

func (r *repo) CreateUser(ctx context.Context, user *model.User) (int64, error) {
	op := "CreateUser"
	log.Printf("%s %v", op, user.ID)

	builderInsert := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, emailColumn, passwordColumn, roleColumn, createdAtColumn, updatedAtColumn).
		Values(user.Name, user.Email, user.Password, user.Role, time.Now(), user.UpdatedAt).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Printf("%s: failed to create builder: %v", op, err)
		return 0, err
	}

	var userID int64
	err = r.db.QueryRow(ctx, query, args...).Scan(&userID)
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

	builderGet := sq.Select(idColumn, nameColumn, emailColumn, roleColumn, createdAtColumn, updatedAtColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: userID})

	query, args, err := builderGet.ToSql()

	if err != nil {
		log.Printf("%s: failed to create builder: %v", op, err)
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		log.Printf("%s: failed to select user: %v", op, err)
		return nil, err
	}

	var user modelRepo.User

	for rows.Next() {
		err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			log.Printf("%s: failed to scan user: %v", op, err)
			return nil, err
		}
	}

	log.Printf("%s: selected user %d", op, userID)
	return converter.ToUserFromRepo(&user), nil
}
