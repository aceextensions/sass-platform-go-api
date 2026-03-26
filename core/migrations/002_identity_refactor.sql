-- Migration to switch identity focus to Email
-- 1. Make email NOT NULL (Ensure all users have an email first)
-- For existing users, if email is null, we might need a default or script, 
-- but in our system currently, Owner always has email.

-- 2. Make phone nullable
ALTER TABLE users ALTER COLUMN phone DROP NOT NULL;

-- 3. Ensure email is not null
-- Check if any nulls exist (shouldn't in new setups, but good practice)
UPDATE users SET email = phone || '@placeholder.com' WHERE email IS NULL;
ALTER TABLE users ALTER COLUMN email SET NOT NULL;
