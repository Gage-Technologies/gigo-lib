-- Add last_read_message to chat users to track the last message in the chat that was read by the user
ALTER TABLE chat_users ADD COLUMN last_read_message BIGINT;

-- Add last_message to chat to track the last message in the chat
ALTER TABLE chat ADD COLUMN last_message BIGINT;

-- Update the chat table with the last message in the chat as long as the
-- chat is not a global, regional, or challenge chat
UPDATE chat SET last_message = (SELECT MAX(_id) FROM chat_messages WHERE chat_id = chat._id) WHERE type != 0 AND type != 1 AND type != 5;