package models

import (
	"database/sql"
	"github.com/kisielk/sqlstruct"
)

type EmailSubscription struct {
	UserId      int64  `json:"user_id" sql:"user_id"`
	UserEmail   string `json:"user_email" sql:"user_email"`
	AllEmails   bool   `json:"all_emails" sql:"all_emails"`
	Streak      bool   `json:"streak" sql:"streak"`
	Pro         bool   `json:"pro" sql:"pro"`
	Newsletter  bool   `json:"newsletter" sql:"newsletter"`
	Inactivity  bool   `json:"inactivity" sql:"inactivity"`
	Messages    bool   `json:"messages" sql:"messages"`
	Referrals   bool   `json:"referrals" sql:"referrals"`
	Promotional bool   `json:"promotional" sql:"promotional"`
}

type EmailSubscriptionSQL struct {
	UserId      int64  `json:"user_id" sql:"user_id"`
	UserEmail   string `json:"user_email" sql:"user_email"`
	AllEmails   bool   `json:"all_emails" sql:"all_emails"`
	Streak      bool   `json:"streak" sql:"streak"`
	Pro         bool   `json:"pro" sql:"pro"`
	Newsletter  bool   `json:"newsletter" sql:"newsletter"`
	Inactivity  bool   `json:"inactivity" sql:"inactivity"`
	Messages    bool   `json:"messages" sql:"messages"`
	Referrals   bool   `json:"referrals" sql:"referrals"`
	Promotional bool   `json:"promotional" sql:"promotional"`
}

func CreateEmailSubscription(userId int64, userEmail string, allEmail bool, streak bool, pro bool, newsletter bool, inactivity bool,
	messages bool, referrals bool, promotional bool) (*EmailSubscription, error) {
	return &EmailSubscription{
		UserId:      userId,
		UserEmail:   userEmail,
		AllEmails:   allEmail,
		Streak:      streak,
		Pro:         pro,
		Newsletter:  newsletter,
		Inactivity:  inactivity,
		Messages:    messages,
		Referrals:   referrals,
		Promotional: promotional,
	}, nil
}

func (i *EmailSubscription) ToSQLNative() []*SQLInsertStatement {
	sqlStatements := make([]*SQLInsertStatement, 0)

	sqlStatements = append(sqlStatements, &SQLInsertStatement{
		Statement: "insert ignore into email_subscription(user_id, user_email, all_emails, streak, pro, newsletter, inactivity, messages, referrals, promotional) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?);",
		Values:    []interface{}{i.UserId, i.UserEmail, i.AllEmails, i.Streak, i.Pro, i.Newsletter, i.Inactivity, i.Messages, i.Referrals, i.Promotional},
	})

	return sqlStatements
}

func EmailSubscriptionFromSQLNative(rows *sql.Rows) (*EmailSubscription, error) {
	// create new EmailSubscription object to load into
	subscriptionSQL := new(EmailSubscriptionSQL)

	// scan row into EmailSubscription object
	err := sqlstruct.Scan(subscriptionSQL, rows)
	if err != nil {
		return nil, err
	}

	// create new UserInactivity for the output
	subscription := &EmailSubscription{
		UserId:      subscriptionSQL.UserId,
		UserEmail:   subscriptionSQL.UserEmail,
		AllEmails:   subscriptionSQL.AllEmails,
		Streak:      subscriptionSQL.Streak,
		Pro:         subscriptionSQL.Pro,
		Newsletter:  subscriptionSQL.Newsletter,
		Inactivity:  subscriptionSQL.Inactivity,
		Messages:    subscriptionSQL.Messages,
		Referrals:   subscriptionSQL.Referrals,
		Promotional: subscriptionSQL.Promotional,
	}

	return subscription, nil
}
