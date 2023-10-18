package models

import (
	"fmt"
	"testing"
	"time"

	ti "github.com/gage-technologies/gigo-lib/db"
)

func TestCreateUserStats(t *testing.T) {
	sts, err := CreateUserStats(69420, 6969420420, 0, false, 0,
		0, time.Minute, time.Minute, 0, 0, 0, time.Now(),
		time.Now().Add(time.Hour*24), nil)
	if err != nil {
		t.Error("\nCreate User Stats Table Failed")
		return
	}

	if sts == nil {
		t.Error("\nCreate User Stats Table Failed\n    Error: creation returned nil")
		return
	}

	if sts.ID != 69420 {
		t.Error("\nCreate User Stats Table Failed\n    Error: wrong id")
		return
	}

	if sts.UserID != 6969420420 {
		t.Error("\nCreate User Stats Table Failed\n    Error: wrong user id")
		return
	}

	if sts.ChallengesCompleted != 0 {
		t.Error("\nCreate User Stats Table Failed\n    Error: wrong challenges completed")
		return
	}

	if sts.CurrentStreak != 0 {
		t.Error("\nCreate User Stats Table Failed\n    Error: wrong current streak")
		return
	}

	if sts.LongestStreak != 0 {
		t.Error("\nCreate User Stats Table Failed\n    Error: wrong longest streak")
		return
	}

	if sts.AvgTime != time.Minute {
		t.Error("\nCreate User Stats Table Failed\n    Error: wrong length")
		return
	}

	t.Log("\nCreate User Stats Table Succeeded")
}

func TestUserStats_ToSQLNative(t *testing.T) {
	rec, err := CreateUserStats(69420, 6969420420, 0, false, 0,
		0, time.Minute, time.Minute, 0, 0, 0, time.Now(),
		time.Now().Add(time.Hour*24), nil)
	if err != nil {
		t.Error("\nCreate User Stats Table Failed")
		return
	}

	statements := rec.ToSQLNative()

	for _, statement := range statements {
		if statement.Statement != "insert ignore into user_stats(_id, user_id, challenges_completed, streak_active, current_streak, longest_streak, total_time_spent, avg_time, days_on_platform, days_on_fire, streak_freezes, streak_freeze_used, xp_gained, date, expiration) values (?, ?, ?, ?, ?, ?, ?, ? ,? ,? ,?, ?, ?, ?, ?);" {
			t.Errorf("\nUser Stats to sql native failed\n    Error: incorrect statement returned")
			return
		}

		if len(statement.Values) != 15 {
			fmt.Println("number of values returned: ", len(statement.Values))
			t.Errorf("\nuser stats to sql native failed\n    Error: incorrect values returned %v", statement.Values)
			return
		}
	}

	t.Log("\nuser stats to sql native succeeded")
}

func TestUserStatsFromSQLNative(t *testing.T) {
	db, err := ti.CreateDatabase("gigo-dev-tidb", "4000", "mysql", "gigo-dev", "gigo-dev", "gigo_dev_test")
	if err != nil {
		t.Error("\nInitialize user stats table sql failed\n    Error: ", err)
		return
	}

	defer db.DB.Exec("DROP TABLE user_stats")
	defer db.DB.Exec("DROP TABLE user_daily_usage")

	post, err := CreateUserStats(69420, 6969420420, 0, false, 0,
		0, time.Minute, time.Minute, 0, 0, 0, time.Now(),
		time.Now().Add(time.Hour*24), nil)
	if err != nil {
		t.Error("\nCreate User Stats Table Failed")
		return
	}

	statements := post.ToSQLNative()

	for _, statement := range statements {
		stmt, err := db.DB.Prepare(statement.Statement)
		if err != nil {
			t.Errorf("\nuser stats from sql native failed\n    Error: %v    statement: %v", err, statement.Statement)
			return
		}

		_, err = stmt.Exec(statement.Values...)
		if err != nil {
			t.Error("\nuser stats from sql native failed\n    Error: ", err)
			return
		}

		rows, err := db.DB.Query("select * from user_stats")
		if err != nil {
			t.Error("\nuser stats from sql native failed\n    Error: ", err)
			return
		}

		if !rows.Next() {
			t.Error("\nuser stats from sql native failed\n    Error: no rows found")
			return
		}

		sts, err := UserStatsFromSQLNative(db, rows)
		if err != nil {
			fmt.Println(rows)
			fmt.Println(statement.Statement)
			fmt.Println(statement.Values)
			t.Errorf("\nuser stats from sql native failed\n    Error: %v", err)
			return
		}

		if sts == nil {
			t.Error("\nuser stats from sql native failed\n    Error: creation returned nil")
			return
		}

		if sts.ID != 69420 {
			t.Error("\nuser stats from sql native failed\n    Error: wrong id")
			return
		}

		if sts.UserID != 6969420420 {
			t.Error("\nuser stats from sql native failed\n    Error: wrong user id")
			return
		}

		if sts.ChallengesCompleted != 0 {
			t.Error("\nuser stats from sql native failed\n    Error: wrong challenges completed")
			return
		}

		if sts.CurrentStreak != 0 {
			t.Error("\nuser stats from sql native failed\n    Error: wrong current streak")
			return
		}

		if sts.LongestStreak != 0 {
			t.Error("\nuser stats from sql native failed\n    Error: wrong longest streak")
			return
		}
	}

}
