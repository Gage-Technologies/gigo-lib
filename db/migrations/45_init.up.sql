CREATE TABLE IF NOT EXISTS email_subscription (
    user_id BIGINT NOT NULL,
    user_email VARCHAR(280) NOT NULL,
    all_emails BOOLEAN NOT NULL,
    streak BOOLEAN NOT NULL,
    pro BOOLEAN NOT NULL,
    newsletter BOOLEAN NOT NULL,
    inactivity BOOLEAN NOT NULL,
    messages BOOLEAN NOT NULL,
    referrals BOOLEAN NOT NULL,
    promotional BOOLEAN NOT NULL,
    PRIMARY KEY (user_id, user_email)
);

INSERT INTO email_subscription (user_id, user_email, all_emails, streak, pro, newsletter, inactivity, messages, referrals, promotional)
SELECT
    _id AS user_id,
    email AS user_email,
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