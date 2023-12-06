CREATE TABLE IF NOT EXISTS email_subscription (
    user_id BIGINT NOT NULL PRIMARY KEY,
    all_emails BOOLEAN NOT NULL,
    streak BOOLEAN NOT NULL,
    pro BOOLEAN NOT NULL,
    newsletter BOOLEAN NOT NULL,
    inactivity BOOLEAN NOT NULL,
    messages BOOLEAN NOT NULL,
    referrals BOOLEAN NOT NULL,
    promotional BOOLEAN NOT NULL
);

INSERT INTO email_subscription (user_id, all_emails, streak, pro, newsletter, inactivity, messages, referrals, promotional)
SELECT
    _id,
    TRUE AS all_emails,
    TRUE AS streak,
    TRUE AS pro,
    TRUE AS newsletter,
    TRUE AS inactivity,
    TRUE AS messages,
    TRUE AS referrals,
    TRUE AS promotional
FROM
    users;