package repositories

import (
	"errors"
	"fmt"

	"github.com/bukharney/giga-chat/modules/entities"
	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	Db *sqlx.DB
}

func NewUsersRepo(db *sqlx.DB) entities.UsersRepository {
	return &UserRepo{Db: db}
}

func (r *UserRepo) Register(req *entities.UsersRegisterReq) (*entities.UsersRegisterRes, error) {
	query := `
	INSERT INTO "users"(
		"username",
		"email",
		"password"
	)
	VALUES ($1, $2, $3)
	RETURNING "id", "username";
	`
	user := new(entities.UsersRegisterRes)

	rows, err := r.Db.Queryx(query, req.Username, req.Email, req.Password)
	if err != nil {
		e := err.Error()
		if e == "sql: no rows in result set" {
			return nil, errors.New("error, user not found")
		} else if e == "pq: duplicate key value violates unique constraint \"users_username_key\"" {
			return nil, errors.New("error, username already exists")
		} else if e == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
			return nil, errors.New("error, email already exists")
		} else {
			return nil, errors.New("error, failed to query")
		}
	}

	defer rows.Close()

	for rows.Next() {
		if err := rows.StructScan(user); err != nil {
			fmt.Println(err.Error())
			return nil, errors.New("error, failed to scan")
		}
	}

	return user, nil
}

func (r *UserRepo) ChangePassword(req *entities.UsersChangePasswordReq) (*entities.UsersChangedRes, error) {
	query := `
	UPDATE "users"
	SET "password" = $1
	WHERE "id" = $2;
	`

	res := new(entities.UsersChangedRes)

	rows, err := r.Db.Queryx(query, req.NewPassword, req.Id)
	if err != nil {
		e := err.Error()
		if e == "sql: no rows in result set" {
			return nil, errors.New("error, user not found")
		} else {
			return nil, errors.New("error, failed to query")
		}
	} else {
		res.Success = true
	}

	for rows.Next() {
		if err := rows.StructScan(res); err != nil {
			fmt.Println(err.Error())
			return nil, errors.New("error, failed to scan")
		}
	}
	return res, nil
}

func (r *UserRepo) GetUserByUsername(username string) (*entities.UsersPassport, error) {
	query := `
	SELECT
	"id",
	"username",
	"password",
	"email"
	FROM "users"
	WHERE "username" = $1;
	`
	res := new(entities.UsersPassport)
	if err := r.Db.Get(res, query, username); err != nil {
		fmt.Println(err.Error())
		return nil, errors.New("error, user not found")
	}
	return res, nil
}

func (r *UserRepo) DeleteAccount(user_id int) (*entities.UsersChangedRes, error) {
	query := `
	DELETE FROM "users"
	WHERE "id" = $1;
	`

	user := new(entities.UsersChangedRes)

	rows, err := r.Db.Queryx(query, user_id)
	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.New("error, failed to delete user")
	}

	defer rows.Close()

	user.Success = true

	return user, nil
}
