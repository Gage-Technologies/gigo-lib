package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	ti "github.com/gage-technologies/gigo-lib/db"
	"github.com/kisielk/sqlstruct"
	"go.opentelemetry.io/otel"
	"time"
)

type DailyUsage struct {
	StartTime   time.Time  `json:"start_time" sql:"start_time"`
	EndTime     *time.Time `json:"end_time" sql:"end_time"`
	OpenSession int        `json:"open_session" sql:"open_session"`
}

type UserStats struct {
	ID                  int64         `json:"id" sql:"_id"`
	UserID              int64         `json:"user_id" sql:"user_id"`
	ChallengesCompleted int           `json:"challenges_completed" sql:"challenges_completed"`
	StreakActive        bool          `json:"streak_active" sql:"streak_active"`
	CurrentStreak       int           `json:"current_streak" sql:"current_streak"`
	LongestStreak       int           `json:"longest_streak" sql:"longest_streak"`
	TotalTimeSpent      time.Duration `json:"total_time_spent" sql:"total_time_spent"`
	AvgTime             time.Duration `json:"avg_time" sql:"avg_time"`

	DailyIntervals []*DailyUsage `json:"daily_intervals" sql:"daily_intervals"`

	DaysOnPlatform   int       `json:"days_on_platform" sql:"days_on_platform"`
	DaysOnFire       int       `json:"days_on_fire" sql:"days_on_fire"`
	StreakFreezes    int       `json:"streak_freezes" sql:"streak_freezes"`
	StreakFreezeUsed bool      `json:"streak_freeze_used" sql:"streak_freeze_used"`
	XpGained         int64     `json:"xp_gained" sql:"xp_gained"`
	Date             time.Time `json:"date" sql:"date"`
	Expiration       time.Time `json:"expiration" sql:"expiration"`
	Closed           bool      `json:"closed" sql:"closed"`
}

type UserStatsSQL struct {
	ID                  int64         `json:"id" sql:"_id"`
	UserID              int64         `json:"user_id" sql:"user_id"`
	ChallengesCompleted int           `json:"challenges_completed" sql:"challenges_completed"`
	StreakActive        bool          `json:"streak_active" sql:"streak_active"`
	CurrentStreak       int           `json:"current_streak" sql:"current_streak"`
	LongestStreak       int           `json:"longest_streak" sql:"longest_streak"`
	TotalTimeSpent      time.Duration `json:"total_time_spent" sql:"total_time_spent"`
	AvgTime             time.Duration `json:"avg_time" sql:"avg_time"`

	DaysOnPlatform   int       `json:"days_on_platform" sql:"days_on_platform"`
	DaysOnFire       int       `json:"days_on_fire" sql:"days_on_fire"`
	StreakFreezes    int       `json:"streak_freezes" sql:"streak_freezes"`
	StreakFreezeUsed bool      `json:"streak_freeze_used" sql:"streak_freeze_used"`
	XpGained         int64     `json:"xp_gained" sql:"xp_gained"`
	Date             time.Time `json:"date" sql:"date"`
	Expiration       time.Time `json:"expiration" sql:"expiration"`
	Closed           bool      `json:"closed" sql:"closed"`

	// Convenience fields for common queries
	OpenSession int    `json:"open_session" sql:"open_session"`
	Timezone    string `json:"timezone" sql:"timezone"`
}

type UserStatsFrontend struct {
	ID                  string        `json:"id" sql:"_id"`
	UserID              string        `json:"user_id" sql:"user_id"`
	ChallengesCompleted string        `json:"challenges_completed" sql:"challenges_completed"`
	StreakActive        bool          `json:"streak_active" sql:"streak_active"`
	CurrentStreak       string        `json:"current_streak" sql:"current_streak"`
	LongestStreak       string        `json:"longest_streak" sql:"longest_streak"`
	TotalTimeSpent      time.Duration `json:"total_time_spent" sql:"total_time_spent"`
	AvgTime             time.Duration `json:"avg_time" sql:"avg_time"`

	DaysOnPlatform   int       `json:"days_on_platform" sql:"days_on_platform"`
	DaysOnFire       int       `json:"days_on_fire" sql:"days_on_fire"`
	StreakFreezes    int       `json:"streak_freezes" sql:"streak_freezes"`
	StreakFreezeUsed bool      `json:"streak_freeze_used" sql:"streak_freeze_used"`
	XpGained         string    `json:"xp_gained" sql:"xp_gained"`
	Date             time.Time `json:"date" sql:"date"`
	Closed           bool      `json:"closed" sql:"closed"`
}

func CreateUserStats(id int64, userId int64, challengesCompleted int, streakActive bool, currentStreak int,
	longestStreak int, totalTimeSpent time.Duration, avgTime time.Duration, daysOnPlatform int, daysOnFire int, streakFreezes int,
	date time.Time, expiration time.Time, dailyUse []*DailyUsage) (*UserStats, error) {
	return &UserStats{
		ID:                  id,
		UserID:              userId,
		ChallengesCompleted: challengesCompleted,
		StreakActive:        streakActive,
		CurrentStreak:       currentStreak,
		LongestStreak:       longestStreak,
		TotalTimeSpent:      totalTimeSpent,
		AvgTime:             avgTime,
		DaysOnPlatform:      daysOnPlatform,
		DaysOnFire:          daysOnFire,
		StreakFreezes:       streakFreezes,
		Date:                date,
		DailyIntervals:      dailyUse,
		Expiration:          expiration,
	}, nil
}

func UserStatsFromSQLNative(db *ti.Database, rows *sql.Rows) (*UserStats, error) {
	// create new user stats object to load into
	userStatsSQL := new(UserStatsSQL)

	// scan row into user stats object
	err := sqlstruct.Scan(userStatsSQL, rows)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error scanning user stats in first scan: %v", err))
	}

	ctx, span := otel.Tracer("gigo-core").Start(context.Background(), "gigo-lib")
	callerName := "UserStatsFromSQLNative"
	res, err := db.QueryContext(ctx, &span, &callerName, "SELECT start_time, end_time, open_session from user_daily_usage where user_id = ? and date = ?", userStatsSQL.UserID, userStatsSQL.Date)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to query for daily usage, err: %v  query: %v", err, fmt.Sprintf("SELECT start_time, end_time, open_session from user_daily_usage where user_id = %v and date = %s;", userStatsSQL.UserID, userStatsSQL.Date)))
	}

	if res.Err() != nil {
		return nil, errors.New(fmt.Sprintf("failed to query for daily usage, err: %v", res.Err()))
	}

	dailyUses := make([]*DailyUsage, 0)

	for res.Next() {
		dailuse := new(DailyUsage)
		err := sqlstruct.Scan(dailuse, res)
		if err != nil {
			return nil, err
		}
		dailyUses = append(dailyUses, dailuse)
	}

	// create new user stats
	userStats := &UserStats{
		ID:                  userStatsSQL.ID,
		UserID:              userStatsSQL.UserID,
		ChallengesCompleted: userStatsSQL.ChallengesCompleted,
		StreakActive:        userStatsSQL.StreakActive,
		CurrentStreak:       userStatsSQL.CurrentStreak,
		LongestStreak:       userStatsSQL.LongestStreak,
		TotalTimeSpent:      userStatsSQL.TotalTimeSpent,
		AvgTime:             userStatsSQL.AvgTime,
		DailyIntervals:      dailyUses,
		DaysOnPlatform:      userStatsSQL.DaysOnPlatform,
		StreakFreezes:       userStatsSQL.StreakFreezes,
		StreakFreezeUsed:    userStatsSQL.StreakFreezeUsed,
		DaysOnFire:          userStatsSQL.DaysOnFire,
		XpGained:            userStatsSQL.XpGained,
		Date:                userStatsSQL.Date,
		Expiration:          userStatsSQL.Expiration,
		Closed:              userStatsSQL.Closed,
	}

	return userStats, nil
}

func (i *UserStats) ToFrontend() *UserStatsFrontend {
	return &UserStatsFrontend{
		ID:                  fmt.Sprintf("%d", i.ID),
		UserID:              fmt.Sprintf("%d", i.UserID),
		ChallengesCompleted: fmt.Sprintf("%d", i.ChallengesCompleted),
		StreakActive:        i.StreakActive,
		CurrentStreak:       fmt.Sprintf("%d", i.CurrentStreak),
		LongestStreak:       fmt.Sprintf("%d", i.LongestStreak),
		TotalTimeSpent:      i.TotalTimeSpent,
		AvgTime:             i.AvgTime,
		DaysOnPlatform:      i.DaysOnPlatform,
		StreakFreezes:       i.StreakFreezes,
		DaysOnFire:          i.DaysOnFire,
		XpGained:            fmt.Sprintf("%d", i.XpGained),
		Date:                i.Date,
		Closed:              i.Closed,
	}
}

func (i *UserStats) ToSQLNative() []*SQLInsertStatement {
	sqlStatements := make([]*SQLInsertStatement, 0)

	sqlStatements = append(sqlStatements, &SQLInsertStatement{
		Statement: "insert ignore into user_stats(_id, user_id, challenges_completed, streak_active, current_streak, longest_streak, total_time_spent, avg_time, days_on_platform, days_on_fire, streak_freezes, streak_freeze_used, xp_gained, date, expiration) values (?, ?, ?, ?, ?, ?, ?, ? ,? ,? ,?, ?, ?, ?, ?);",
		Values:    []interface{}{i.ID, i.UserID, i.ChallengesCompleted, i.StreakActive, i.CurrentStreak, i.LongestStreak, i.TotalTimeSpent, i.AvgTime, i.DaysOnPlatform, i.DaysOnFire, i.StreakFreezes, i.StreakFreezeUsed, i.XpGained, i.Date, i.Expiration},
	})

	for _, d := range i.DailyIntervals {
		sqlStatements = append(sqlStatements, &SQLInsertStatement{
			Statement: "insert ignore into user_daily_usage(user_id, start_time, end_time, open_session, date) values (?, ?, ?, ?, ?);",
			Values:    []interface{}{i.UserID, d.StartTime, d.EndTime, d.OpenSession, i.Date},
		})
	}

	// create insertion statement and return
	return sqlStatements
}
