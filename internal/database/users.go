package database

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type User struct {
	ID             string    `db:"ID"`
	FullName       string    `db:"FullName"`
	Email          string    `db:"Email"`
	HashedPassword string    `db:"HashedPassword"`
	ProfilePicture string    `db:"ProfilePicture"`
	CreatedAt      time.Time `db:"CreatedAt"`
	UpdatedAt      time.Time `db:"UpdatedAt"`
}

func (db *DB) InsertUser(email, hashedPassword string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
		INSERT INTO Users (Created, Email, HashedPassword)
		VALUES (?, ?, ?)`

	result, err := db.ExecContext(ctx, query, time.Now(), email, hashedPassword)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), err
}

func (db *DB) GetUser(id int) (*User, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var user User

	query := `SELECT * FROM Users WHERE ID = ?`

	err := db.GetContext(ctx, &user, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, false, nil
	}

	return &user, true, err
}

func (db *DB) GetUserByEmail(email string) (*User, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var user User

	query := `SELECT * FROM Users WHERE Email = ?`

	err := db.GetContext(ctx, &user, query, email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, false, nil
	}

	return &user, true, err
}

func (db *DB) UpdateUserHashedPassword(id int, hashedPassword string) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `UPDATE Users SET HashedPassword = ? WHERE ID = ?`

	_, err := db.ExecContext(ctx, query, hashedPassword, id)
	return err
}
