-- MySQL v5.7 - Change message column in chat from varchar(500) to text
ALTER TABLE chat_messages MODIFY message text NOT NULL;