package models

import (
	"fmt"
	ti "github.com/gage-technologies/gigo-lib/db"
	"github.com/gage-technologies/gigo-lib/utils"
	"math"
	"reflect"
	"testing"
)

func TestCreateUser(t *testing.T) {

	badges := []int64{1, 2}

	user, err := CreateUser(69, "test", "testpass", "testemail@email.com",
		"phone", UserStatusBasic, "test", badges, []int64{1, 2, 3},
		"first", "last", 23, "", DefaultUserStart, "America/Chicago",
		AvatarSettings{}, 0)

	if err != nil {
		t.Errorf("failed to create user, err: %v", err)
		return
	}

	if user == nil {
		t.Errorf("failed to create user, err: user returned nil")
		return
	}

	if user.ID != 69 {
		t.Errorf("failed to create user, err: user returned incorrect id")
		return
	}

	if user.UserName != "test" {
		t.Errorf("failed to create user, err: user returned incorrect pfppath")
		return
	}

	isPass, err := utils.CheckPassword("testpass", user.Password)
	if err != nil {
		t.Errorf("failed to create user, password err: %v", err)
		return
	}

	if !isPass {
		t.Errorf("failed to create user, err: password did not match")
		return
	}

	if user.Email != "testemail@email.com" {
		t.Errorf("failed to create user, err: user returned incorrect email")
		return
	}

	if user.UserStatus != UserStatusBasic {
		t.Errorf("failed to create user, err: user returned incorrect user status")
		return
	}

	if user.Bio != "test" {
		t.Errorf("failed to create user, err: user returned incorrect user bio")
		return
	}

	if user.Badges[1] != badges[1] {
		t.Errorf("failed to create user, err: user returned incorrect user badges")
		return
	}

	if user.XP != 0 {
		t.Errorf("failed to create user, err: user returned incorrect user xp")
		return
	}

	if user.Level != Level1 {
		t.Errorf("failed to create user, err: user returned incorrect user level")
		return
	}

	if user.Tier != Tier1 {
		t.Errorf("failed to create user, err: user returned incorrect user tier")
		return
	}

	if user.Rank != NoobRank {
		t.Errorf("failed to create user, err: user returned incorrect user rank")
		return
	}

	if user.Coffee != 0 {
		t.Errorf("failed to create user, err: user returned incorrect user coffee")
		return
	}

	if user.Timezone != "America/Chicago" {
		t.Errorf("failed to create user, err: user returned incorrect user timezone")
		return
	}

	t.Logf("\nCreateUser Succeeded\n")
	return
}

func TestInitializeUserTableSQL(t *testing.T) {
	db, err := ti.CreateDatabase("gigo-dev-tidb", "4000", "mysql", "gigo-dev", "gigo-dev", "gigo_dev_test")
	if err != nil {
		t.Error("\nInitialize User Table\n    Error: ", err)
		return
	}

	defer db.DB.Exec("delete from users")
	defer db.DB.Exec("delete from user_active_times")
	defer db.DB.Exec("delete from user_badges")
	defer db.DB.Exec("delete from user_saved_posts")

	res, err := db.DB.Query("SELECT * FROM information_schema.tables WHERE table_schema = 'gigo_dev_test' AND table_name = 'users' LIMIT 1;")
	if err != nil {
		t.Error("\nInitialize failed\n    Error: ", err)
		return
	}

	if !res.Next() {
		t.Error("\nInitialize failed\n    Error: table was not created")
		return
	}

	res, err = db.DB.Query("SELECT * FROM information_schema.tables WHERE table_schema = 'gigo_dev_test' AND table_name = 'user_active_times' LIMIT 1;")
	if err != nil {
		t.Error("\nInitialize failed\n    Error: ", err)
		return
	}

	if !res.Next() {
		t.Error("\nInitialize failed\n    Error: table was not created")
		return
	}

	res, err = db.DB.Query("SELECT * FROM information_schema.tables WHERE table_schema = 'gigo_dev_test' AND table_name = 'user_badges' LIMIT 1;")
	if err != nil {
		t.Error("\nInitialize failed\n    Error: ", err)
		return
	}

	if !res.Next() {
		t.Error("\nInitialize failed\n    Error: table was not created")
		return
	}

	res, err = db.DB.Query("SELECT * FROM information_schema.tables WHERE table_schema = 'gigo_dev_test' AND table_name = 'user_saved_posts' LIMIT 1;")
	if err != nil {
		t.Error("\nInitialize failed\n    Error: ", err)
		return
	}

	if !res.Next() {
		t.Error("\nInitialize failed\n    Error: table was not created")
		return
	}

	t.Log("\nInitialize succeeded")

}

func TestUser_ToFromSQLNative(t *testing.T) {
	db, err := ti.CreateDatabase("gigo-dev-tidb", "4000", "mysql", "gigo-dev", "gigo-dev", "gigo_dev_test")
	if err != nil {
		t.Error("\nToSQlNative User Table\n    Error: ", err)
		return
	}

	defer db.DB.Exec("delete from users")
	defer db.DB.Exec("delete from user_active_times")
	defer db.DB.Exec("delete from user_badges")
	defer db.DB.Exec("delete from user_saved_posts")

	badges := []int64{1, 2}

	user, err := CreateUser(69, "test", "testpass", "testemail@email.com",
		"phone", UserStatusBasic, "test", badges, []int64{1, 2, 3},
		"first", "last", 23, "", DefaultUserStart, "America/Chicago",
		AvatarSettings{}, 0)
	if err != nil {
		t.Errorf("failed user to sql native, err: %v", err)
		return
	}

	stripeUser := "stripe user"
	stripSubscription := "strip subscription"

	user.StripeUser = &stripeUser
	user.StripeSubscription = &stripSubscription

	statements, err := user.ToSQLNative()
	if err != nil {
		t.Error("\nToSQLNative User Table\n    Error: ", err)
		return
	}

	if len(statements) < 1 {
		t.Errorf("failed user to sql native, err: no statements returned")
		return
	}

	if statements[0].Statement != "insert ignore into user_badges(user_id, badge_id) values(?, ?);" {
		fmt.Println(statements[0].Statement)
		t.Errorf("\nToSQlNative User Table\n    Error: %v", err)
		return
	}

	if !reflect.DeepEqual(statements[0].Values, []interface{}{int64(69), int64(1)}) {
		fmt.Println(statements[0].Values)
		t.Error("\nTo sql native failed\n    Error: incorrect values returned for user badges table")
		return
	}

	if statements[2].Statement != "insert ignore into user_saved_posts(user_id, post_id) values(?, ?);" {
		t.Errorf("\nToSQlNative User Table\n    Error: %v", err)
		return
	}

	if !reflect.DeepEqual(statements[2].Values, []interface{}{int64(69), int64(1)}) {
		t.Error("\nTo sql native failed\n    Error: incorrect values returned for user badges table")
		return
	}

	if statements[5].Statement != "insert ignore into users(_id, email, phone, user_status, user_name, password, bio, xp, level, tier, user_rank, coffee, first_name, last_name, gitea_id, external_auth, created_at, stripe_user, stripe_subscription, workspace_settings, encrypted_service_key, follower_count, start_user_info, highest_score, timezone, avatar_settings, broadcast_threshold, avatar_reward, stripe_account, exclusive_agreement, reset_token, has_broadcast, holiday_themes, tutorials) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);" {
		fmt.Println(statements[5].Statement)
		t.Errorf("\nToSQlNative User Table\n    Error: incorrect statement")
		return
	}

	tx, err := db.DB.Begin()
	if err != nil {
		t.Errorf("\nTo sql native failed\n    Error: failed to start tx: %v", err)
		return
	}

	for _, s := range statements {
		_, err = tx.Exec(s.Statement, s.Values...)
		if err != nil {
			_ = tx.Rollback()
			t.Errorf("\nTo sql native failed\n    Error: failed to execute insertion statement: %v\n    statement: %s\n    values: %v", err, s.Statement, s.Values)
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		t.Errorf(fmt.Sprintf("failed to commit transaction, err: %v", err))
		return
	}

	res, err := db.DB.Query("select * from users where _id = 69 limit 1")
	if err != nil {
		t.Errorf("\nToSQlNative User Table\n    Error: %v", err)
		return
	}

	if !res.Next() {
		t.Errorf("\nToSQlNative User Table\n    Error: no rows returned")
		return
	}

	loadedUser, err := UserFromSQLNative(db, res)
	if err != nil {
		t.Errorf("\nToSQlNative User Table\n    Error: %v", err)
		return
	}

	if user.ID != loadedUser.ID {
		t.Errorf("\nToSQlNative User Table\n    Error: incorrect values returned for user table")
		return
	}

	if user.Email != loadedUser.Email {
		t.Errorf("\nToSQlNative User Table\n    Error: incorrect values returned for user table")
		return
	}

	if user.UserStatus != UserStatusPremium {
		t.Errorf("\nToSQlNative User Table\n    Error: incorrect values returned for user table")
		return
	}

	if user.Password != loadedUser.Password {
		t.Errorf("\nToSQlNative User Table\n    Error: incorrect values returned for user table")
		return
	}

	if user.Bio != loadedUser.Bio {
		t.Errorf("\nToSQlNative User Table\n    Error: incorrect values returned for user table")
		return
	}

	if 2 < math.Abs(float64(user.CreatedAt.Unix()-loadedUser.CreatedAt.Unix())) {
		fmt.Println(user.CreatedAt)
		fmt.Println(loadedUser.CreatedAt)
		t.Errorf("\nToSQlNative User Table\n    Error: incorrect values returned for user table")
		return
	}

	if *user.StripeUser != *loadedUser.StripeUser {
		fmt.Println(*user.StripeUser)
		fmt.Println(*loadedUser.StripeUser)
		t.Errorf("\nToSQlNative User Table\n    Error: incorrect values returned for user table")
		return
	}

	if *user.StripeSubscription != *loadedUser.StripeSubscription {
		t.Errorf("\nToSQlNative User Table\n    Error: incorrect values returned for user table")
		return
	}

	if user.Level != loadedUser.Level {
		t.Errorf("\nToSQlNative User Table\n    Error: incorrect values returned for user table")
		return
	}

	if user.Tier != loadedUser.Tier {
		t.Errorf("\nToSQlNative User Table\n    Error: incorrect values returned for user table")
		return
	}

	if user.UserName != loadedUser.UserName {
		t.Errorf("\nToSQlNative User Table\n    Error: incorrect values returned for user table")
		return
	}

	if user.AuthRole != loadedUser.AuthRole {
		t.Errorf("\nToSQlNative User Table\n    Error: incorrect values returned for user table")
		return
	}

	if !reflect.DeepEqual(user.WorkspaceSettings, loadedUser.WorkspaceSettings) {
		t.Errorf("\nToSQlNative User Table\n    Error: incorrect values returned for user table")
		return
	}

	if user.Timezone != loadedUser.Timezone {
		t.Errorf("\nToSQlNative User Table\n    Error: incorrect values returned for user table")
		return
	}

	t.Logf("User To SQL Native Succeeded")
}

func TestInsertUser(t *testing.T) {
	db, err := ti.CreateDatabase("gigo-dev-tidb", "4000", "mysql", "gigo-dev", "gigo-dev", "gigo_dev_test")
	if err != nil {
		t.Error("\nToSQlNative User Table\n    Error: ", err)
		return
	}
	defer db.DB.Exec("delete from users where _id = 6942069")

	badges := []int64{1, 2}

	user, err := CreateUser(6942069, "test", "testpass", "testemail@email.com",
		"phone", UserStatusBasic, "test", badges, []int64{1, 2, 3},
		"first", "last", 23, "", DefaultUserStart, "America/Chicago",
		AvatarSettings{}, 0)
	if err != nil {
		t.Errorf("\nTestInsertUser\n    Error: %v", err)
		return
	}

	insertStatement, err := user.ToSQLNative()
	if err != nil {
		t.Errorf("\nTestInsertUser\n    Error: %v", err)
		return
	}

	for _, statement := range insertStatement {
		fmt.Println("Statement:", statement.Statement)
		fmt.Println("Values:", statement.Values)
	}

	for _, statement := range insertStatement {
		_, err = db.DB.Exec(statement.Statement, statement.Values...)
		if err != nil {
			t.Errorf("\nTestInsertUser\n    Error: %v", err)
		}
	}

	return
}
