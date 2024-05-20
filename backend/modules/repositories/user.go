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

func (r *UserRepo) AddFriend(req *entities.FriendReq) (*entities.FriendRes, error) {
	query := `
	INSERT INTO "friends"(
	"from_user_id",
	"to_user_id",
	"status"
	)
	VALUES ($1, $2, $3);
	`

	user := new(entities.FriendRes)

	rows, err := r.Db.Queryx(query, req.UserId, req.FriendId, req.Status)
	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.New("error, failed to add friend")
	}

	defer rows.Close()

	return user, nil
}

func (r *UserRepo) GetFriendsReq(user_id int) ([]entities.FriendInfoRes, error) {
	query := `
	SELECT
	"users"."id",
	"users"."username"
	FROM "users"
	JOIN "friends"
	ON "users"."id" = "friends"."from_user_id"
	WHERE "friends"."to_user_id" = $1 AND "friends"."status" = 0;
	`

	var friends []entities.FriendInfoRes

	err := r.Db.Select(&friends, query, user_id)
	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.New("error, failed to get friends")
	}

	return friends, nil
}

func (r *UserRepo) GetFriendReq(user_id int, friend_id int) (int, error) {
	query := `
	SELECT "status"
	FROM "friends"
	WHERE "to_user_id" = $1 AND "from_user_id" = $2;
	`

	var friend = new(entities.FriendRes)

	err := r.Db.Select(&friend, query, user_id)
	if err != nil {
		fmt.Println(err.Error())
		return 0, errors.New("error, failed to get friends")
	}

	return friend.Status, nil
}

func (r *UserRepo) GetFriends(user_id int) ([]entities.FriendInfoRes, error) {
	query := `
	SELECT
	"users"."id",
	"users"."username"
	FROM "users"
	JOIN "friends"
	ON "users"."id" = "friends"."to_user_id"
	WHERE "friends"."from_user_id" = $1 AND "friends"."status" = 1;
	`

	var friends []entities.FriendInfoRes

	err := r.Db.Select(&friends, query, user_id)
	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.New("error, failed to get friends")
	}

	return friends, nil
}
