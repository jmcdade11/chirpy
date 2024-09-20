package database

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

type RefreshToken struct {
	ID         int       `json:"id"`
	Token      string    `json:"token"`
	Expiration time.Time `json:"expiration"`
}

var ErrRefreshAlreadyExists = errors.New("refresh token already exists")

func (db *DB) CreateRefreshToken(id int) (RefreshToken, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return RefreshToken{}, err
	}

	token, err := getNewToken()
	if err != nil {
		return RefreshToken{}, err
	}

	expiration := time.Now().UTC().Add(time.Hour * 24 * 60)

	refreshToken := RefreshToken{
		ID:         id,
		Token:      token,
		Expiration: expiration,
	}
	dbStructure.RefreshTokens[id] = refreshToken

	err = db.writeDB(dbStructure)
	if err != nil {
		fmt.Printf("Couldn't write refresh to db\n")
		return RefreshToken{}, err
	}

	return refreshToken, nil
}

func (db *DB) GetRefreshToken(token string) (RefreshToken, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return RefreshToken{}, err
	}
	for _, refreshToken := range dbStructure.RefreshTokens {
		if refreshToken.Token == token {
			return refreshToken, nil
		}
	}
	return RefreshToken{}, ErrNotExist
}

func (db *DB) DeleteRefreshToken(id int) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}
	_, ok := dbStructure.RefreshTokens[id]
	if !ok {
		return ErrNotExist
	}

	delete(dbStructure.RefreshTokens, id)
	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}

func getNewToken() (string, error) {
	c := 10
	b := make([]byte, c)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
