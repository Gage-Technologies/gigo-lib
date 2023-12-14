package models

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	ti "github.com/gage-technologies/gigo-lib/db"
	"github.com/gage-technologies/gigo-lib/session"
	"github.com/gage-technologies/gigo-lib/storage"
	"github.com/gage-technologies/gigo-lib/utils"
	"github.com/gage-technologies/gotp"
	"github.com/kisielk/sqlstruct"
	"go.opentelemetry.io/otel"
)

type UserStatus int

const (
	UserStatusBasic UserStatus = iota
	UserStatusPremium
)

func (s UserStatus) String() string {
	switch s {
	case UserStatusBasic:
		return "Basic"
	case UserStatusPremium:
		return "Premium"
	}
	return "Unknown"
}

type UserTutorial struct {
	All           bool `json:"all" sql:"all"`
	Home          bool `json:"home" sql:"home"`
	Challenge     bool `json:"challenge" sql:"challenge"`
	Workspace     bool `json:"workspace" sql:"workspace"`
	Nemesis       bool `json:"nemesis" sql:"nemesis"`
	Stats         bool `json:"stats" sql:"stats"`
	CreateProject bool `json:"create_project" sql:"create_project"`
	Launchpad     bool `json:"launchpad" sql:"launchpad"`
	Vscode        bool `json:"vscode" sql:"vscode"`
}

var DefaultUserTutorial = UserTutorial{}

type User struct {
	ID                  int64              `json:"_id" sql:"_id"`
	UserName            string             `json:"user_name" sql:"user_name"`
	Password            string             `json:"password"`
	Email               string             `json:"email" sql:"email"`
	Phone               string             `json:"phone" sql:"phone"`
	UserStatus          UserStatus         `json:"user_status" sql:"user_status"`
	Bio                 string             `json:"bio" sql:"bio"`
	Badges              []int64            `json:"badges" sql:"badges"`
	XP                  uint64             `json:"xp" sql:"xp"`
	Level               LevelType          `json:"level" sql:"level"`
	Tier                TierType           `json:"tier" sql:"tier"`
	Rank                RankType           `json:"user_rank" sql:"user_rank"`
	Coffee              uint64             `json:"coffee" sql:"coffee"`
	SavedPosts          []int64            `json:"saved_posts,omitempty" sql:"saved_posts"`
	FirstName           string             `json:"first_name" sql:"first_name"`
	LastName            string             `json:"last_name" sql:"last_name"`
	CreatedAt           time.Time          `json:"created_at" sql:"created_at"`
	WorkspaceSettings   *WorkspaceSettings `json:"workspace_settings" sql:"workspace_settings"`
	EncryptedServiceKey string             `json:"encrypted_service_key" sql:"encrypted_service_key"`
	StartUserInfo       *UserStart         `json:"start_user_info" sql:"start_user_info"`
	HighestScore        uint64             `json:"highest_score" sql:"highest_score"`
	Timezone            string             `json:"timezone" sql:"timezone"`
	AvatarSettings      *AvatarSettings    `json:"avatar_settings" sql:"avatar_settings"`
	BroadcastThreshold  uint64             `json:"broadcast_threshold" sql:"broadcast_threshold"`
	AvatarReward        *int64             `json:"avatar_reward" sql:"avatar_reward"`
	ExclusiveAgreement  bool               `json:"exclusive_agreement" sql:"exclusive_agreement"`
	ResetToken          *string            `json:"reset_token" sql:"reset_token"`
	HasBroadcast        bool               `json:"has_broadcast" sql:"has_broadcast"`
	HolidayThemes       bool               `json:"holiday_themes" sql:"holiday_themes"`
	Tutorials           *UserTutorial      `json:"tutorials" sql:"tutorials"`

	// Gitea
	GiteaID     int64 `json:"gitea_id" sql:"gitea_id"`
	IsEphemeral bool  `json:"is_ephemeral" sql:"is_ephemeral"`

	// Auth
	Otp          *string            `json:"otp,omitempty" sql:"otp"`
	OtpValidated *bool              `json:"otp_validated,omitempty" sql:"otp_validated"`
	AuthRole     AuthenticationRole `json:"auth_role" sql:"auth_role"`
	ExternalAuth string             `json:"external_auth" sql:"external_auth"`

	// Stripe
	StripeUser         *string `json:"stripe_user,omitempty" sql:"stripe_user,omitempty"`
	StripeAccount      *string `json:"stripe_account,omitempty" sql:"stripe_account,omitempty"`
	StripeSubscription *string `json:"stripe_subscription" sql:"stripe_subscription"`
	FollowerCount      uint64  `json:"follower_count" sql:"follower_count"`

	ReferredBy *int64 `json:"referred_by" sql:"referred_by"`
}

type UserSQL struct {
	ID                  int64      `json:"_id" sql:"_id"`
	UserName            string     `json:"user_name" sql:"user_name"`
	Password            string     `json:"password"`
	Email               string     `json:"email" sql:"email"`
	Phone               string     `json:"phone" sql:"phone"`
	UserStatus          UserStatus `json:"user_status" sql:"user_status"`
	Bio                 string     `json:"bio" sql:"bio"`
	XP                  uint64     `json:"xp" sql:"xp"`
	Level               LevelType  `json:"level" sql:"level"`
	Tier                TierType   `json:"tier" sql:"tier"`
	Rank                RankType   `json:"user_rank" sql:"user_rank"`
	Coffee              uint64     `json:"coffee" sql:"coffee"`
	FirstName           string     `json:"first_name" sql:"first_name"`
	LastName            string     `json:"last_name" sql:"last_name"`
	GiteaID             int64      `json:"gitea_id" sql:"gitea_id"`
	CreatedAt           time.Time  `json:"created_at" sql:"created_at"`
	WorkspaceSettings   []byte     `json:"workspace_settings" sql:"workspace_settings"`
	EncryptedServiceKey []byte     `json:"encrypted_service_key" sql:"encrypted_service_key"`
	StartUserInfo       []byte     `json:"start_user_info" sql:"start_user_info"`
	HighestScore        uint64     `json:"highest_score" sql:"highest_score"`
	Timezone            string     `json:"timezone" sql:"timezone"`
	AvatarSettings      []byte     `json:"avatar_settings" sql:"avatar_settings"`
	BroadcastThreshold  uint64     `json:"broadcast_threshold" sql:"broadcast_threshold"`
	AvatarReward        *int64     `json:"avatar_reward" sql:"avatar_reward"`
	ExclusiveAgreement  bool       `json:"exclusive_agreement" sql:"exclusive_agreement"`
	ResetToken          *string    `json:"reset_token" sql:"reset_token"`
	HasBroadcast        bool       `json:"has_broadcast" sql:"has_broadcast"`
	HolidayThemes       bool       `json:"holiday_themes" sql:"holiday_themes"`
	Tutorials           []byte     `json:"tutorials" sql:"tutorials"`

	IsEphemeral bool `json:"is_ephemeral" sql:"is_ephemeral"`

	// Auth
	Otp          *string            `json:"otp,omitempty" sql:"otp"`
	OtpValidated *bool              `json:"otp_validated,omitempty" sql:"otp_validated"`
	AuthRole     AuthenticationRole `json:"auth_role" sql:"auth_role"`
	ExternalAuth string             `json:"external_auth" sql:"external_auth"`

	StripeUser         *string `json:"stripe_user" sql:"stripe_user"`
	StripeAccount      *string `json:"stripe_account" sql:"stripe_account"`
	StripeSubscription *string `json:"stripe_subscription" sql:"stripe_subscription"`
	FollowerCount      uint64  `json:"follower_count" sql:"follower_count"`

	ReferredBy *int64 `json:"referred_by" sql:"referred_by"`
}

type UserSearch struct {
	ID       int64  `json:"_id" sql:"_id"`
	UserName string `json:"user_name" sql:"user_name"`
}

type UserFrontend struct {
	ID                 string     `json:"_id" sql:"_id"`
	PFPPath            string     `json:"pfp_path" sql:"pfp_path"`
	UserName           string     `json:"user_name" sql:"user_name"`
	Email              string     `json:"email" sql:"email"`
	Phone              string     `json:"phone" sql:"phone"`
	UserStatus         UserStatus `json:"user_status" sql:"user_status"`
	UserStatusString   string     `json:"user_status_string" sql:"user_status_string"`
	Bio                string     `json:"bio" sql:"bio"`
	XP                 uint64     `json:"xp" sql:"xp"`
	Level              LevelType  `json:"level" sql:"level"`
	Tier               TierType   `json:"tier" sql:"tier"`
	Rank               RankType   `json:"user_rank" sql:"user_rank"`
	Coffee             uint64     `json:"coffee" sql:"coffee"`
	SavedPosts         []string   `json:"saved_posts,omitempty" sql:"saved_posts"`
	FirstName          string     `json:"first_name" sql:"first_name"`
	LastName           string     `json:"last_name" sql:"last_name"`
	CreatedAt          time.Time  `json:"created_at" sql:"created_at"`
	FollowerCount      uint64     `json:"follower_count" sql:"follower_count"`
	HighestScore       uint64     `json:"highest_score" sql:"highest_score"`
	Timezone           string     `json:"timezone" sql:"timezone"`
	BroadcastThreshold uint64     `json:"broadcast_threshold" sql:"broadcast_threshold"`
	AvatarReward       *string    `json:"avatar_reward" sql:"avatar_reward"`
	ExclusiveAgreement bool       `json:"exclusive_agreement" sql:"exclusive_agreement"`
	ResetToken         *string    `json:"reset_token" sql:"reset_token"`
	HasBroadcast       bool       `json:"has_broadcast" sql:"has_broadcast"`
	HolidayThemes      bool       `json:"holiday_themes" sql:"holiday_themes"`
}

func CreateUser(id int64, userName string, password string, email string, phone string,
	userStatus UserStatus, bio string, badges []int64, savedPosts []int64, firstName string,
	lasName string, giteaID int64, externalAuth string, starInfo UserStart, timezone string, avatar AvatarSettings,
	broadcastThreshold uint64, referredBy *int64) (*User, error) {

	// hash password
	hashedPass, err := utils.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash user password: %v", err)
	}

	// create internal service password
	serviceKey, err := session.GenerateServicePassword()
	if err != nil {
		return nil, fmt.Errorf("failed to generate internal service password: %v", err)
	}

	// encrypt internal service password using plain-text user password
	encryptedServiceKey, err := session.EncryptServicePassword(serviceKey, []byte(password))
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt internal service password: %v", err)
	}

	return &User{
		ID:                  id,
		UserName:            userName,
		Password:            hashedPass,
		Email:               email,
		Phone:               phone,
		UserStatus:          userStatus,
		Bio:                 bio,
		Badges:              badges,
		XP:                  0,
		Level:               Level1,
		Tier:                Tier1,
		Rank:                NoobRank,
		Coffee:              0,
		SavedPosts:          savedPosts,
		FirstName:           firstName,
		LastName:            lasName,
		GiteaID:             giteaID,
		CreatedAt:           time.Now(),
		ExternalAuth:        externalAuth,
		WorkspaceSettings:   &DefaultWorkspaceSettings,
		EncryptedServiceKey: encryptedServiceKey,
		FollowerCount:       0,
		StartUserInfo:       &starInfo,
		HighestScore:        0,
		Timezone:            timezone,
		AvatarSettings:      &avatar,
		BroadcastThreshold:  broadcastThreshold,
		ExclusiveAgreement:  false,
		ResetToken:          nil,
		HasBroadcast:        false,
		HolidayThemes:       true,
		Tutorials:           &DefaultUserTutorial,
		ReferredBy:          referredBy,
	}, nil
}

func (i *User) EditUser(userName *string, password *string, email *string, phone *string, userStatus *UserStatus,
	bio *string, badges []int64, savedPosts []int64, firstName *string, lasName *string, giteaID *int64,
	externalAuth *string, starInfo *UserStart, timezone *string, avatar *AvatarSettings, broadcastThreshold *uint64) (*User, *SQLInsertStatement, error) {

	params := make([]string, 0)
	values := make([]interface{}, 0)

	if userName != nil {
		i.UserName = *userName
		params = append(params, "user_name = ?")
		values = append(values, *userName)
	}
	if password != nil {

		// hash password
		hashedPass, err := utils.HashPassword(*password)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to hash user password: %v", err)
		}

		i.Password = hashedPass

		// create internal service password
		serviceKey, err := session.GenerateServicePassword()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to generate internal service password: %v", err)
		}

		// encrypt internal service password using plain-text user password
		encryptedServiceKey, err := session.EncryptServicePassword(serviceKey, []byte(*password))
		if err != nil {
			return nil, nil, fmt.Errorf("failed to encrypt internal service password: %v", err)
		}

		i.EncryptedServiceKey = encryptedServiceKey

		// base64 decode encrypted service password
		encryptedServicePassword, err := base64.RawStdEncoding.DecodeString(i.EncryptedServiceKey)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to decode encrypted service password: %v", err)
		}

		params = append(params, "password =?")
		values = append(values, i.Password)

		params = append(params, "encrypted_service_key =?")
		values = append(values, encryptedServicePassword)
	}
	if email != nil {
		i.Email = *email
		params = append(params, "email =?")
		values = append(values, *email)
	}
	if phone != nil {
		i.Phone = *phone
		params = append(params, "phone =?")
		values = append(values, *phone)
	}
	if userStatus != nil {
		i.UserStatus = *userStatus
		params = append(params, "user_status =?")
		values = append(values, *userStatus)
	}
	if bio != nil {
		i.Bio = *bio
		params = append(params, "bio =?")
		values = append(values, *bio)
	}
	if badges != nil {
		i.Badges = badges
		params = append(params, "badges =?")
		values = append(values, badges)
	}
	if savedPosts != nil {
		i.SavedPosts = savedPosts
		params = append(params, "saved_posts =?")
		values = append(values, savedPosts)
	}
	if firstName != nil {
		i.FirstName = *firstName
		params = append(params, "first_name =?")
		values = append(values, *firstName)
	}
	if lasName != nil {
		i.LastName = *lasName
		params = append(params, "last_name =?")
		values = append(values, *lasName)
	}
	if giteaID != nil {
		i.GiteaID = *giteaID
		params = append(params, "gitea_id =?")
		values = append(values, *giteaID)
	}
	if externalAuth != nil {
		i.ExternalAuth = *externalAuth
		params = append(params, "external_auth =?")
		values = append(values, *externalAuth)
	}
	if starInfo != nil {
		i.StartUserInfo = starInfo
		// marshall workspace settings to json
		startSettings, err := json.Marshal(i.StartUserInfo)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to marshall starter user settings: %v", err)
		}

		params = append(params, "start_user_info =?")
		values = append(values, startSettings)
	}
	if timezone != nil {
		i.Timezone = *timezone
		params = append(params, "timezone =?")
		values = append(values, *timezone)
	}
	if avatar != nil {
		i.AvatarSettings = avatar
		// marshall workspace settings to json
		avatarSettings, err := json.Marshal(i.AvatarSettings)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to marshall starter user settings: %v", err)
		}
		params = append(params, "avatar_settings =?")
		values = append(values, avatarSettings)
	}
	if broadcastThreshold != nil {
		i.BroadcastThreshold = *broadcastThreshold
		params = append(params, "broadcast_threshold =?")
		values = append(values, *broadcastThreshold)
	}

	values = append(values, i.ID)

	query := "UPDATE users SET " + strings.Join(params, ", ") + " WHERE id =?"
	sqlStatements := &SQLInsertStatement{
		Statement: query,
		Values:    values,
	}

	return i, sqlStatements, nil

}

func UserFromSQLNative(db *ti.Database, rows *sql.Rows) (*User, error) {
	// create new user object to load into
	userSQL := new(UserSQL)

	// scan row into user object
	err := sqlstruct.Scan(userSQL, rows)
	if err != nil {
		return nil, err
	}

	ctx, span := otel.Tracer("gigo-core").Start(context.Background(), "gigo-lib")
	callerName := "UserFromSQLNative"
	badgeRows, err := db.QueryContext(ctx, &span, &callerName, "select badge_id from user_badges where user_id = ?", userSQL.ID)
	if err != nil {
		return nil, err
	}

	defer badgeRows.Close()

	badges := make([]int64, 0)

	for badgeRows.Next() {
		var badge int64
		err = badgeRows.Scan(&badge)
		if err != nil {
			return nil, err
		}
		badges = append(badges, badge)
	}

	ctx, span = otel.Tracer("gigo-core").Start(context.Background(), "gigo-lib")
	savedPostsRows, err := db.QueryContext(ctx, &span, &callerName, "select post_id from user_saved_posts where user_id = ?", userSQL.ID)
	if err != nil {
		return nil, err
	}

	defer savedPostsRows.Close()

	savedPosts := make([]int64, 0)

	for savedPostsRows.Next() {
		var savedPost int64
		err = savedPostsRows.Scan(&savedPost)
		if err != nil {
			return nil, err
		}
		savedPosts = append(savedPosts, savedPost)
	}

	// create variable to decode workspace settings into
	var workspaceSettings *WorkspaceSettings

	// unmarshall workspace setting from json buffer to WorkspaceSettings type
	err = json.Unmarshal(userSQL.WorkspaceSettings, &workspaceSettings)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall workspace settings: %v", err)
	}

	// create variable to decode workspace settings into
	var userStart *UserStart

	// unmarshall workspace setting from json buffer to WorkspaceSettings type
	err = json.Unmarshal(userSQL.StartUserInfo, &userStart)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall user start: %v", err)
	}

	// create variable to decode workspace settings into
	var avatarSetting *AvatarSettings

	// unmarshall workspace setting from json buffer to WorkspaceSettings type
	err = json.Unmarshal(userSQL.AvatarSettings, &avatarSetting)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall avatar settings: %v", err)
	}

	// create variable to decode tutorials into
	var tutorials *UserTutorial

	// unmarshall workspace setting from json buffer to WorkspaceSettings type
	if len(userSQL.Tutorials) > 0 {
		err = json.Unmarshal(userSQL.Tutorials, &tutorials)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshall tutorials: %v", err)
		}
	} else {
		tutorials = &DefaultUserTutorial
	}

	// create new user for the output
	user := &User{
		ID:                  userSQL.ID,
		UserName:            userSQL.UserName,
		Password:            userSQL.Password,
		Email:               userSQL.Email,
		Phone:               userSQL.Phone,
		UserStatus:          userSQL.UserStatus,
		Bio:                 userSQL.Bio,
		XP:                  userSQL.XP,
		Level:               userSQL.Level,
		Tier:                userSQL.Tier,
		Rank:                userSQL.Rank,
		Badges:              badges,
		SavedPosts:          savedPosts,
		Otp:                 userSQL.Otp,
		OtpValidated:        userSQL.OtpValidated,
		AuthRole:            userSQL.AuthRole,
		StripeUser:          userSQL.StripeUser,
		StripeAccount:       userSQL.StripeAccount,
		FirstName:           userSQL.FirstName,
		LastName:            userSQL.LastName,
		GiteaID:             userSQL.GiteaID,
		CreatedAt:           userSQL.CreatedAt,
		ExternalAuth:        userSQL.ExternalAuth,
		StripeSubscription:  userSQL.StripeSubscription,
		WorkspaceSettings:   workspaceSettings,
		EncryptedServiceKey: base64.RawStdEncoding.EncodeToString(userSQL.EncryptedServiceKey),
		FollowerCount:       userSQL.FollowerCount,
		StartUserInfo:       userStart,
		HighestScore:        userSQL.HighestScore,
		Timezone:            userSQL.Timezone,
		AvatarSettings:      avatarSetting,
		BroadcastThreshold:  userSQL.BroadcastThreshold,
		AvatarReward:        userSQL.AvatarReward,
		ExclusiveAgreement:  userSQL.ExclusiveAgreement,
		ResetToken:          userSQL.ResetToken,
		HasBroadcast:        userSQL.HasBroadcast,
		HolidayThemes:       userSQL.HolidayThemes,
		Tutorials:           tutorials,
		IsEphemeral:         userSQL.IsEphemeral,
		ReferredBy:          userSQL.ReferredBy,
	}

	return user, nil
}

func (i *User) ToSearch() *UserSearch {
	return &UserSearch{
		ID:       i.ID,
		UserName: i.UserName,
	}
}

func (i *User) ToFrontend() (*UserFrontend, error) {
	badges := make([]string, 0)

	for _, b := range i.Badges {
		badges = append(badges, fmt.Sprintf("%d", b))
	}

	savedPosts := make([]string, 0)

	for p := range i.SavedPosts {
		savedPosts = append(savedPosts, fmt.Sprintf("%d", p))
	}

	var avatarReward *string = nil
	if i.AvatarReward != nil {
		reward := fmt.Sprintf("%d", *i.AvatarReward)
		avatarReward = &reward
	}

	// create new user frontend
	mf := &UserFrontend{
		ID:                 fmt.Sprintf("%d", i.ID),
		PFPPath:            fmt.Sprintf("/static/user/pfp/%v", i.ID),
		UserName:           i.UserName,
		Email:              i.Email,
		Phone:              i.Phone,
		UserStatus:         i.UserStatus,
		UserStatusString:   i.UserStatus.String(),
		Bio:                i.Bio,
		XP:                 i.XP,
		Level:              i.Level,
		Tier:               i.Tier,
		Rank:               i.Rank,
		Coffee:             i.Coffee,
		SavedPosts:         savedPosts,
		FirstName:          i.FirstName,
		LastName:           i.LastName,
		CreatedAt:          i.CreatedAt,
		FollowerCount:      i.FollowerCount,
		HighestScore:       i.HighestScore,
		BroadcastThreshold: i.BroadcastThreshold,
		AvatarReward:       avatarReward,
		ExclusiveAgreement: i.ExclusiveAgreement,
		ResetToken:         i.ResetToken,
		HasBroadcast:       i.HasBroadcast,
		HolidayThemes:      i.HolidayThemes,
	}

	return mf, nil
}

func (i *User) ToSQLNative() ([]*SQLInsertStatement, error) {

	sqlStatements := make([]*SQLInsertStatement, 0)

	if len(i.Badges) > 0 {
		for _, b := range i.Badges {
			badgeStatement := SQLInsertStatement{
				Statement: "insert ignore into user_badges(user_id, badge_id) values(?, ?);",
				Values:    []interface{}{i.ID, b},
			}

			sqlStatements = append(sqlStatements, &badgeStatement)
		}
	}

	if len(i.SavedPosts) > 0 {
		for _, p := range i.SavedPosts {
			savedPostsStatement := SQLInsertStatement{
				Statement: "insert ignore into user_saved_posts(user_id, post_id) values(?, ?);",
				Values:    []interface{}{i.ID, p},
			}

			sqlStatements = append(sqlStatements, &savedPostsStatement)
		}
	}

	// marshall workspace settings to json
	workspaceSettings, err := json.Marshal(i.WorkspaceSettings)
	if err != nil {
		return nil, fmt.Errorf("failed to marshall workspace settings: %v", err)
	}

	// marshall workspace settings to json
	startSettings, err := json.Marshal(i.StartUserInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to marshall starter user settings: %v", err)
	}

	// marshall workspace settings to json
	avatarSettings, err := json.Marshal(i.AvatarSettings)
	if err != nil {
		return nil, fmt.Errorf("failed to marshall starter user settings: %v", err)
	}

	// marshall tutorials to json
	var tutorials []byte
	if i.Tutorials != nil {
		tutorials, err = json.Marshal(i.Tutorials)
		if err != nil {
			return nil, fmt.Errorf("failed to marshall tutorials: %v", err)
		}
	}

	// base64 decode encrypted service password
	encryptedServicePassword, err := base64.RawStdEncoding.DecodeString(i.EncryptedServiceKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode encrypted service password: %v", err)
	}

	sqlStatements = append(sqlStatements, &SQLInsertStatement{
		Statement: "insert ignore into users(_id, email, phone, user_status, user_name, password, bio, xp, level, tier, user_rank, coffee, first_name, last_name, gitea_id, external_auth, created_at, stripe_user, stripe_subscription, workspace_settings, encrypted_service_key, follower_count, start_user_info, highest_score, timezone, avatar_settings, broadcast_threshold, avatar_reward, stripe_account, exclusive_agreement, reset_token, has_broadcast, holiday_themes, tutorials, is_ephemeral, referred_by) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);",
		Values:    []interface{}{i.ID, i.Email, i.Phone, i.UserStatus, i.UserName, i.Password, i.Bio, i.XP, i.Level, i.Tier, i.Rank, i.Coffee, i.FirstName, i.LastName, i.GiteaID, i.ExternalAuth, i.CreatedAt, i.StripeUser, i.StripeSubscription, workspaceSettings, encryptedServicePassword, i.FollowerCount, startSettings, i.HighestScore, i.Timezone, avatarSettings, i.BroadcastThreshold, i.AvatarReward, i.StripeAccount, i.ExclusiveAgreement, i.ResetToken, i.HasBroadcast, i.HolidayThemes, tutorials, i.IsEphemeral, i.ReferredBy},
	})

	// create insertion statement and return
	return sqlStatements, nil
}

func (i *User) GenerateUserOtpUri(db *ti.Database) (map[string]interface{}, error) {
	// ensure that the user has not already set up otp verification
	if i.OtpValidated != nil && *i.OtpValidated {
		return map[string]interface{}{"message": "user has already set up 2fa"}, fmt.Errorf("user with validated otp attempted to re-generate otp secret")
	}

	// generate a 64 byte (256 bit) random secret key
	secret := gotp.RandomSecret(64)

	// create an otp instance derived from the secret key
	otp := gotp.NewDefaultTOTP(secret)

	// generate a url that can be used for linking otp apps
	otpUri := otp.ProvisioningUri(i.UserName, "Gage")

	// update the user object in the database to contain the secret and mark the validation as false pending the first confirmation
	st := SQLInsertStatement{
		Statement: "update users set otp_validated = false, otp = ? where _id = ?;",
		Values:    []interface{}{fmt.Sprintf("%v", secret), i.ID},
	}

	ctx, span := otel.Tracer("gigo-core").Start(context.Background(), "gigo-lib")
	callerName := "GenerateUserOtpUri"
	res, err := db.ExecContext(ctx, &span, &callerName, st.Statement, st.Values...)
	if err != nil {
		return nil,
			errors.New(fmt.Sprintf("failed to update opt validated field in user: %v, err: %v", i.ID, err))
	}

	rowsAffect, err := res.RowsAffected()
	if err != nil {
		return nil,
			errors.New(fmt.Sprintf("failed to update opt validated/otp and retrieve affected rows for field in user: %v, err: %v", i.ID, err))
	}

	if rowsAffect < 1 {
		return nil,
			errors.New(fmt.Sprintf("failed to update opt validated/otp and retrieve affected rows for field in user: %v, err: no rows affected", i.ID))

	}

	// return the otp uri to the frontend
	return map[string]interface{}{"otp_uri": otpUri}, nil
}

func (i *User) VerifyUserOtp(db *ti.Database, storageEngine storage.Storage, otp string, ip string) (map[string]interface{}, string, error) {
	// ensure that the otp has been initialized
	if i.Otp == nil {
		return map[string]interface{}{"message": "user has not setup 2fa"}, "", fmt.Errorf("otp was nil in user during verify otp call")
	}

	// use the user secret to create a new otp instance and validate the otp code
	valid := gotp.NewDefaultTOTP(*i.Otp).Verify(otp, time.Now().Unix())

	// create an empty string to hold the token
	token := ""

	// conditionally create a valid token for the user session
	if valid {
		// create a token for the user session
		t, err := utils.CreateExternalJWT(storageEngine, fmt.Sprintf("%d", i.ID), ip, 12, 0, map[string]interface{}{
			"auth_role": i.AuthRole,
			"temporary": false,
			"init_temp": false,
			"otp":       true,
			"otp_valid": true,
		})
		if err != nil {
			return nil, "", err
		}

		// assign the token to the outer scope variable
		token = t

		// conditionally update the user if this is the first time they are verifying their otp login
		if i.OtpValidated != nil && !*i.OtpValidated {
			// update user marking their otp login as validated

			st := SQLInsertStatement{
				Statement: "update users set otp_validated = true where _id = ?;",
				Values:    []interface{}{i.ID},
			}
			ctx, span := otel.Tracer("gigo-core").Start(context.Background(), "gigo-lib")
			callerName := "VerifyUserOtp"
			res, err := db.ExecContext(ctx, &span, &callerName, st.Statement, st.Values...)
			if err != nil {
				return nil,
					"",
					errors.New(fmt.Sprintf("failed to update opt validated field in user: %v, err: %v", i.ID, err))
			}

			rowsAffect, err := res.RowsAffected()
			if err != nil {
				return nil,
					"",
					errors.New(fmt.Sprintf("failed to update opt validated and retrieve affected rows for field in user: %v, err: %v", i.ID, err))
			}

			if rowsAffect < 1 {
				return nil,
					"",
					errors.New(fmt.Sprintf("failed to update opt validated and retrieve affected rows for field in user: %v, err: no rows affected", i.ID))

			}

		}
	}

	// return the authentication and token to the frontend
	return map[string]interface{}{
		"auth":  valid,
		"token": token,
	}, token, nil
}
