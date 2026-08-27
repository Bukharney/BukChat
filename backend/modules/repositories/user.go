package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/bukharney/bukchat/modules/entities"
	"github.com/bukharney/bukchat/pkg/apperrors"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type UserRepo struct {
	Db *sqlx.DB
}

func NewUsersRepo(db *sqlx.DB) entities.UsersRepository {
	return &UserRepo{Db: db}
}

func (r *UserRepo) Register(ctx context.Context, req *entities.UsersRegisterReq) (*entities.UsersRegisterRes, error) {
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

	rows, err := r.Db.QueryxContext(ctx, query, req.Username, req.Email, req.Password)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				if strings.Contains(pqErr.Constraint, "users_username_key") {
					return nil, apperrors.ErrUsernameExists
				}
				if strings.Contains(pqErr.Constraint, "users_email_key") {
					return nil, apperrors.ErrEmailExists
				}
			}
		}
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to register user: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.StructScan(user); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
	}

	return user, nil
}

func (r *UserRepo) ChangePassword(ctx context.Context, req *entities.UsersChangePasswordReq) (*entities.UsersChangedRes, error) {
	query := `
	UPDATE "users"
	SET "password" = $1
	WHERE "id" = $2;
	`

	res := new(entities.UsersChangedRes)

	rows, err := r.Db.QueryxContext(ctx, query, req.NewPassword, req.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to change password: %w", err)
	}
	defer rows.Close()

	res.Success = true
	for rows.Next() {
		if err := rows.StructScan(res); err != nil {
			return nil, fmt.Errorf("failed to scan password change result: %w", err)
		}
	}
	return res, nil
}

func (r *UserRepo) GetUserByUsername(ctx context.Context, username string) (*entities.UsersPassport, error) {
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
	if err := r.Db.GetContext(ctx, res, query, username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	return res, nil
}

func (r *UserRepo) DeleteAccount(ctx context.Context, user_id int) (*entities.UsersChangedRes, error) {
	query := `
	DELETE FROM "users"
	WHERE "id" = $1;
	`

	user := new(entities.UsersChangedRes)

	rows, err := r.Db.QueryxContext(ctx, query, user_id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete user: %w", err)
	}
	defer rows.Close()

	user.Success = true
	return user, nil
}

func (r *UserRepo) AddFriend(ctx context.Context, req *entities.FriendReq) (*entities.FriendRes, error) {
	query := `
	INSERT INTO "friends"(
	"from_user_id",
	"to_user_id",
	"room_id",
	"status"
	)
	VALUES ($1, $2, $3, $4)
	RETURNING "from_user_id", "to_user_id", "room_id", "status", "created_at";
	`

	user := new(entities.FriendRes)

	rows, err := r.Db.QueryxContext(ctx, query, req.UserId, req.FriendId, nil, req.Status)
	if err != nil {
		return nil, fmt.Errorf("failed to add friend: %w", err)
	}
	defer rows.Close()

	return user, nil
}

func (r *UserRepo) AcceptFriendReq(ctx context.Context, user_id int, friend_id int, room_id int) (*entities.FriendRes, error) {
	query := `
	UPDATE "friends"
	SET "status" = 1, "room_id" = $3
	WHERE "to_user_id" = $1 AND "from_user_id" = $2
	RETURNING "from_user_id", "to_user_id", "status", "created_at";
	`

	user := new(entities.FriendRes)

	rows, err := r.Db.QueryxContext(ctx, query, friend_id, user_id, room_id)
	if err != nil {
		return nil, fmt.Errorf("failed to accept friend request: %w", err)
	}
	defer rows.Close()

	return user, nil
}

func (r *UserRepo) RejectFriend(ctx context.Context, user_id int, friend_id int) (*entities.UsersChangedRes, error) {
	query := `
	DELETE FROM "friends"
	WHERE "to_user_id" = $1 AND "from_user_id" = $2 OR "to_user_id" = $2 AND "from_user_id" = $1;
	`

	user := new(entities.UsersChangedRes)

	rows, err := r.Db.QueryxContext(ctx, query, user_id, friend_id)
	if err != nil {
		return nil, fmt.Errorf("failed to reject friend request: %w", err)
	}
	defer rows.Close()

	user.Success = true
	return user, nil
}

func (r *UserRepo) GetFriendsReq(ctx context.Context, user_id int) ([]entities.FriendInfoRes, error) {
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

	err := r.Db.SelectContext(ctx, &friends, query, user_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []entities.FriendInfoRes{}, nil
		}
		return nil, fmt.Errorf("failed to get friend requests: %w", err)
	}

	return friends, nil
}

func (r *UserRepo) GetFriendReq(ctx context.Context, user_id int, friend_id int) (*entities.FriendRes, error) {
	query := `
	SELECT 
	"from_user_id", 
	"to_user_id", 
	"status", 
	"created_at"
	FROM "friends"
	WHERE ("to_user_id" = $1 AND "from_user_id" = $2) OR ("to_user_id" = $2 AND "from_user_id" = $1);
	`

	var friend = new(entities.FriendRes)

	err := r.Db.GetContext(ctx, friend, query, user_id, friend_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get friend request: %w", err)
	}

	return friend, nil
}

func (r *UserRepo) GetFriends(ctx context.Context, user_id int) ([]entities.FriendInfoRes, error) {
	query := `
	SELECT
	"users"."id",
	"users"."username",
	"friends"."status",
	"friends"."room_id"
  FROM "users"
  JOIN "friends"
	ON "users"."id" = "friends"."from_user_id" OR "users"."id" = "friends"."to_user_id"
  WHERE ("friends"."to_user_id" = $1 OR "friends"."from_user_id" = $1) AND "friends"."status" = 1;
	`

	var friends []entities.FriendInfoRes

	err := r.Db.SelectContext(ctx, &friends, query, user_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get friends: %w", err)
	}

	if len(friends) == 0 {
		return nil, apperrors.ErrNotFound
	}

	length := len(friends)
	for i := 0; i < length; i++ {
		if friends[i].Id == user_id {
			friends = append(friends[:i], friends[i+1:]...)
			length--
			i--
		}
	}

	return friends, nil
}
