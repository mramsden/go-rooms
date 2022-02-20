package users

import (
	"context"
	"database/sql"
)

type UsersRepository struct {
	db *sql.DB
}

func (r *UsersRepository) CreateUser(ctx context.Context, username, hashedPassword string) (*User, error) {
	u := User{}
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO users(username, password)
		VALUES ($1, $2)
		RETURNING id, username, password`, username, hashedPassword).Scan(&u.Id, &u.Username, &u.HashedPassword)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UsersRepository) GetUser(ctx context.Context, id int) (*User, error) {
	u := User{}
	err := r.db.QueryRowContext(ctx, `
		SELECT id, username, password FROM users 
		WHERE id = $1`, id).Scan(&u.Id, &u.Username, &u.HashedPassword)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UsersRepository) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	u := User{}
	err := r.db.QueryRowContext(ctx, `
		SELECT id, username, password FROM users
		WHERE username = $1`, username).Scan(&u.Id, &u.Username, &u.HashedPassword)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
