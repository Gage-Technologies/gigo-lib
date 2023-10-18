package models

import (
	"context"
	ti "github.com/gage-technologies/gigo-lib/db"
	session3 "github.com/gage-technologies/gigo-lib/session"
	"github.com/go-redis/redis/v8"
	"testing"
	"time"
)

func TestStoreLoadUserSession(t *testing.T) {
	db, err := ti.CreateDatabase("gigo-dev-tidb", "4000", "mysql", "gigo-dev", "gigo-dev", "gigo_dev_test")
	if err != nil {
		t.Fatalf("\n%s failed\n    Error: %v", t.Name(), err)
	}

	rdb := redis.NewClient(&redis.Options{Addr: "gigo-dev-redis:6379", Password: "gigo-dev", DB: 7})

	defer rdb.Del(context.TODO(), "gigo-user-sess-420")
	defer db.DB.Exec("drop table user_session_key")

	serviceKey, err := session3.GenerateServicePassword()
	if err != nil {
		t.Fatalf("\n%s failed\n    Error: %v", t.Name(), err)
	}

	session, err := CreateUserSession(69, 420, serviceKey, time.Now().Add(time.Hour))
	if err != nil {
		t.Fatalf("\n%s failed\n    Error: %v", t.Name(), err)
	}

	err = session.Store(db, rdb)
	if err != nil {
		t.Fatalf("\n%s failed\n    Error: %v", t.Name(), err)
	}

	session2, err := LoadUserSession(db, rdb, 420)
	if err != nil {
		t.Fatalf("\n%s failed\n    Error: %v", t.Name(), err)
	}

	if session.ID != session2.ID {
		t.Fatalf("\n%s failed\n    %+v != %+v", t.Name(), session, session2)
	}

	if session.UserID != session2.UserID {
		t.Fatalf("\n%s failed\n    %+v != %+v", t.Name(), session, session2)
	}

	if session.EncryptedServiceKey != session2.EncryptedServiceKey {
		t.Fatalf("\n%s failed\n    %+v != %+v", t.Name(), session, session2)
	}

	if session.SessionKey.ID != session2.SessionKey.ID {
		t.Fatalf("\n%s failed\n    %+v != %+v", t.Name(), session, session2)
	}

	if session.SessionKey.Key != session2.SessionKey.Key {
		t.Fatalf("\n%s failed\n    %+v != %+v", t.Name(), session, session2)
	}

	t.Logf("\n%s succeeded", t.Name())
}

func TestUserSession_GetServiceKey(t *testing.T) {
	serviceKey, err := session3.GenerateServicePassword()
	if err != nil {
		t.Fatalf("\n%s failed\n    Error: %v", t.Name(), err)
	}

	session, err := CreateUserSession(69, 420, serviceKey, time.Now().Add(time.Hour))
	if err != nil {
		t.Fatalf("\n%s failed\n    Error: %v", t.Name(), err)
	}

	key, err := session.GetServiceKey()
	if err != nil {
		t.Fatalf("\n%s failed\n    Error: %v", t.Name(), err)
	}

	if serviceKey != key {
		t.Fatalf("\n%s failed\n    %+v!= %+v", t.Name(), serviceKey, key)
	}

	t.Logf("\n%s succeeded", t.Name())
}
