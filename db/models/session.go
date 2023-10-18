package models

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	ti "github.com/gage-technologies/gigo-lib/db"
	"github.com/gage-technologies/gigo-lib/session"
	"github.com/go-redis/redis/v8"
	"go.opentelemetry.io/otel"
	"time"
)

// UserSession
//
//	UserSessions are used to manage ephemeral data for
//	a user's session. This includes the ephemeral service
//	key. UserSessions are store in redis for the duration
//	of a session or until their expiration. Once a user
//	session has expired the user for that session should
//	be required to login before continuing interactions
//	on the system.
type UserSession struct {
	ID                  int64           `json:"_id"`
	UserID              int64           `json:"user_id"`
	Started             time.Time       `json:"started"`
	Expiration          time.Time       `json:"expiration"`
	EncryptedServiceKey string          `json:"encrypted_service_key"`
	SessionKey          *UserSessionKey `json:"session_key,omitempty"`
}

func CreateUserSession(id int64, userId int64, serviceKey string, expiration time.Time) (*UserSession, error) {
	// create a new password to encrypt the service key
	pass, err := session.GenerateServicePassword()
	if err != nil {
		return nil, fmt.Errorf("failed to generate password for service key: %v", err)
	}

	// base64 decode password to get bytes for encrypting
	passBytes, err := base64.RawStdEncoding.DecodeString(pass)
	if err != nil {
		return nil, fmt.Errorf("failed to decode password for service key: %v", err)
	}

	// encrypt the service key
	encryptedServiceKey, err := session.EncryptServicePassword(serviceKey, passBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt service key: %v", err)
	}

	// return new user session and key
	return &UserSession{
		ID:                  id,
		UserID:              userId,
		Started:             time.Now(),
		Expiration:          expiration,
		EncryptedServiceKey: encryptedServiceKey,
		SessionKey: &UserSessionKey{
			ID:         id,
			Expiration: expiration,
			Key:        pass,
		},
	}, nil
}

// LoadUserSession
//
//	Loads existing user session from redis and its key from sql. If there
//	is no session for the passed user an error of "no session" is returned
func LoadUserSession(db *ti.Database, rdb redis.UniversalClient, userId int64) (*UserSession, error) {
	// retrieve session from redis
	sessionBytes, err := rdb.Get(context.Background(), fmt.Sprintf("gigo-user-sess-%d", userId)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("no session")
		}
		return nil, fmt.Errorf("failed to retrieve session from redis: %v", err)
	}

	// unmarshal session
	var session UserSession
	err = json.Unmarshal(sessionBytes, &session)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal session from redis: %v", err)
	}

	// god forbid someone leaks a session key to the redis storage
	// clean it up here and scream bloody murder
	if session.SessionKey != nil {
		_ = rdb.Del(context.Background(), fmt.Sprintf("gigo-user-sess-%d", userId))
		ctx, span := otel.Tracer("gigo-core").Start(context.Background(), "gigo-lib")
		callerName := "LoadUserSession"
		_, _ = db.ExecContext(ctx, &span, &callerName, "delete from user_session_key where _id = ?", session.ID)
		return nil, fmt.Errorf("session key has been leaked to redis for %d", userId)
	}

	ctx, span := otel.Tracer("gigo-core").Start(context.Background(), "gigo-lib")
	callerName := "LoadUserSession"
	// query for session key
	res, err := db.QueryContext(ctx, &span, &callerName, "select * from user_session_key where _id = ? limit 1", session.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to query user_session_key: %v", err)
	}

	// ensure closure of cursor
	defer res.Close()

	// attempt to load key into first position of cursor
	if !res.Next() {
		return nil, fmt.Errorf("no session key found")
	}

	// load key from cursor
	session.SessionKey, err = UserSessionKeyFromSQLNative(res)
	if err != nil {
		return nil, fmt.Errorf("failed to load session key: %v", err)
	}

	return &session, nil
}

// Store
//
//	Stores the session in redis with the session expiration and stores the
//	session key in sql with the session expiration
func (s *UserSession) Store(db *ti.Database, rdb redis.UniversalClient) error {
	ctx, span := otel.Tracer("gigo-core").Start(context.Background(), "gigo-lib")
	callerName := "StoreUserSession"
	// open tx for insertion
	tx, err := db.BeginTx(ctx, &span, &callerName, nil)
	if err != nil {
		return fmt.Errorf("failed to open transaction: %v", err)
	}
	defer tx.Rollback()

	// format key for insertion
	statements, err := s.SessionKey.ToSQLNative()
	if err != nil {
		return fmt.Errorf("failed to format key for insertion: %v", err)
	}

	// iterate insertion statements performing sql insertion
	for _, statement := range statements {
		ctx, span = otel.Tracer("gigo-core").Start(context.Background(), "gigo-lib")
		_, err = tx.ExecContext(ctx, &callerName, statement.Statement, statement.Values...)
		if err != nil {
			return fmt.Errorf("failed to insert user session key: %v", err)
		}
	}

	// save key
	sessKey := s.SessionKey

	// set key to nil before we save to redis
	s.SessionKey = nil

	// json marshall user session for redis insertion
	bytes, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("failed to marshal user session: %v", err)
	}

	// set key
	s.SessionKey = sessKey

	// insert user session into redis
	res := rdb.Set(context.TODO(), fmt.Sprintf("gigo-user-sess-%d", s.UserID), bytes, time.Until(s.Expiration))
	if res.Err() != nil {
		return fmt.Errorf("failed to insert user session into redis: %v", res.Err())
	}

	// commit tx
	err = tx.Commit(&callerName)
	if err != nil {
		// clean up user session
		_ = rdb.Del(context.TODO(), fmt.Sprintf("gigo-user-sess-%d", s.UserID))
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

// GetServiceKey
//
//	Decrypts the EncryptedServiceKey using the SessionKey
//	and returns the plain-text service key
func (s *UserSession) GetServiceKey() (string, error) {
	// ensure that the session key is loaded
	if s.SessionKey == nil {
		return "", fmt.Errorf("session key not loaded")
	}

	// base64 decode session key
	keyBytes, err := base64.RawStdEncoding.DecodeString(s.SessionKey.Key)
	if err != nil {
		return "", fmt.Errorf("failed to decode session key: %v", err)
	}

	// decrypt the service key
	serviceKey, err := session.DecryptServicePassword(s.EncryptedServiceKey, keyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt service key: %v", err)
	}

	return serviceKey, nil
}
