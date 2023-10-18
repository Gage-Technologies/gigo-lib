package models

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"time"
)

type UserSessionKey struct {
	ID         int64     `sql:"_id"`
	Key        string    `sql:"_key"`
	Expiration time.Time `sql:"expiration"`
}

func CreateUserSessionKey(id int64, key string, expiration time.Time) *UserSessionKey {
	return &UserSessionKey{
		ID:         id,
		Key:        key,
		Expiration: expiration,
	}
}

func (k *UserSessionKey) ToSQLNative() ([]SQLInsertStatement, error) {
	// decode key into bytes
	buf, err := base64.RawStdEncoding.DecodeString(k.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to base64 decode key: %v", err)
	}

	return []SQLInsertStatement{
		{
			Statement: "insert ignore into user_session_key(_id, _key, expiration) values (?, ?, ?)",
			Values:    []interface{}{k.ID, buf, k.Expiration},
		},
	}, nil
}

func UserSessionKeyFromSQLNative(rows *sql.Rows) (*UserSessionKey, error) {
	// create new instance to scan into
	var usk UserSessionKey

	// create buffer to hold key bytes
	var keyBytes []byte

	// scan from rows
	err := rows.Scan(&usk.ID, &keyBytes, &usk.Expiration)
	if err != nil {
		return nil, fmt.Errorf("failed to scan user sesssion key from cursor")
	}

	// base64 encode key
	usk.Key = base64.RawStdEncoding.EncodeToString(keyBytes)

	return &usk, nil
}
