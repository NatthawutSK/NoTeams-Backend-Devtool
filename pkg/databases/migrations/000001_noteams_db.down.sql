BEGIN;

DROP TRIGGER IF EXISTS set_updated_at_timestamp_User_table ON "User";
DROP TRIGGER IF EXISTS set_updated_at_timestamp_Oauth_table ON "Oauth";
DROP TRIGGER IF EXISTS set_updated_at_timestamp_Team_table ON "Team";
DROP TRIGGER IF EXISTS set_updated_at_timestamp_Permission_table ON "Permission";
DROP TRIGGER IF EXISTS set_updated_at_timestamp_Task_table ON "Task";



DROP FUNCTION IF EXISTS set_updated_at_column();

DROP TABLE IF EXISTS "User" CASCADE;
DROP TABLE IF EXISTS "Oauth" CASCADE;
DROP TABLE IF EXISTS "Team" CASCADE;
DROP TABLE IF EXISTS "TeamMember" CASCADE;
DROP TABLE IF EXISTS "Permission" CASCADE;
DROP TABLE IF EXISTS "File" CASCADE;
DROP TABLE IF EXISTS "Task" CASCADE;

DROP SEQUENCE IF EXISTS user_id_seq;
DROP SEQUENCE IF EXISTS team_id_seq;

DROP TYPE IF EXISTS "task_status_enum";
DROP TYPE IF EXISTS "role_enum";

COMMIT;