BEGIN;

--Set timezone
SET TIME ZONE 'Asia/Bangkok';

--Install uuid extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

--users_id -> U000001
--Create sequence
CREATE SEQUENCE user_id_seq START WITH 1 INCREMENT BY 1;
CREATE SEQUENCE team_id_seq START WITH 1 INCREMENT BY 1;

--Auto update
CREATE OR REPLACE FUNCTION set_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;   
END;
$$ language 'plpgsql';

CREATE TABLE "User" (
  "user_id" VARCHAR(7) PRIMARY KEY DEFAULT CONCAT('U', LPAD(NEXTVAL('user_id_seq')::TEXT, 6, '0')),
  "username" VARCHAR UNIQUE NOT NULL,
  "password" VARCHAR NOT NULL,
  "email" VARCHAR UNIQUE NOT NULL,
  "dob" VARCHAR NOT NULL,
  "phone" VARCHAR UNIQUE NOT NULL,
  "bio" TEXT,
  "avatar" VARCHAR DEFAULT 'https://www.seekpng.com/png/detail/41-410093_circled-user-icon-user-profile-icon-png.png',
  "created_at" TIMESTAMP NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE "Oauth" (
  "oauth_id" uuid NOT NULL UNIQUE PRIMARY KEY DEFAULT uuid_generate_v4(),
  "user_id" VARCHAR NOT NULL,
  "access_token" VARCHAR NOT NULL,
  "refresh_token" VARCHAR NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE "Team" (
  "team_id" VARCHAR(7) PRIMARY KEY DEFAULT CONCAT('T', LPAD(NEXTVAL('team_id_seq')::TEXT, 6, '0')),
  "team_name" VARCHAR NOT NULL,
  "team_desc" TEXT NOT NULL,
  "team_code" VARCHAR UNIQUE NOT NULL,
  "team_poster" VARCHAR DEFAULT 'https://icons.veryicon.com/png/o/miscellaneous/site-icon-library/team-28.png',
  "created_at" TIMESTAMP NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);

-- Create the enum type
CREATE TYPE role_enum AS ENUM ('OWNER', 'MEMBER');


CREATE TABLE "TeamMember" (
  "member_id" uuid NOT NULL UNIQUE PRIMARY KEY DEFAULT uuid_generate_v4(),
  "user_id" VARCHAR NOT NULL,
  "team_id" VARCHAR NOT NULL,
  "role" role_enum NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT now()  
);

CREATE TABLE "Permission" (
  "permission_id" uuid NOT NULL UNIQUE PRIMARY KEY DEFAULT uuid_generate_v4(),
  "allow_task" BOOLEAN NOT NULL DEFAULT FALSE,
  "allow_file" BOOLEAN NOT NULL DEFAULT FALSE,
  "allow_invite" BOOLEAN NOT NULL DEFAULT FALSE,
  "team_id" VARCHAR NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE "File" (
  "file_id" uuid NOT NULL UNIQUE PRIMARY KEY DEFAULT uuid_generate_v4(),
  "file_name" VARCHAR NOT NULL,
  "file_url" VARCHAR NOT NULL,
  "team_id" VARCHAR NOT NULL,
  "user_id" VARCHAR NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT now()
);

-- Create the enum type
CREATE TYPE task_status_enum AS ENUM ('TODO', 'DOING', 'DONE');


CREATE TABLE "Task" (
  "task_id" uuid NOT NULL UNIQUE PRIMARY KEY DEFAULT uuid_generate_v4(),
  "task_name" VARCHAR NOT NULL,
  "task_desc" TEXT,
  "task_status" task_status_enum  NOT NULL,
  "task_deadline" VARCHAR,
  "team_id" VARCHAR NOT NULL,
  "user_id" VARCHAR,
  "created_at" TIMESTAMP NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);


ALTER TABLE "Oauth" ADD FOREIGN KEY ("user_id") REFERENCES "User" ("user_id")  ON DELETE CASCADE;
ALTER TABLE "TeamMember" ADD FOREIGN KEY ("user_id") REFERENCES "User" ("user_id")  ON DELETE CASCADE;
ALTER TABLE "TeamMember" ADD FOREIGN KEY ("team_id") REFERENCES "Team" ("team_id")  ON DELETE CASCADE;
ALTER TABLE "Permission" ADD FOREIGN KEY ("team_id") REFERENCES "Team" ("team_id")  ON DELETE CASCADE;
ALTER TABLE "File" ADD FOREIGN KEY ("team_id") REFERENCES "Team" ("team_id")  ON DELETE CASCADE;
ALTER TABLE "File" ADD FOREIGN KEY ("user_id") REFERENCES "User" ("user_id")  ON DELETE CASCADE;
ALTER TABLE "Task" ADD FOREIGN KEY ("team_id") REFERENCES "Team" ("team_id")  ON DELETE CASCADE;
ALTER TABLE "Task" ADD FOREIGN KEY ("user_id") REFERENCES "User" ("user_id")  ON DELETE CASCADE;


CREATE TRIGGER set_updated_at_timestamp_User_table BEFORE UPDATE ON "User" FOR EACH ROW EXECUTE PROCEDURE set_updated_at_column();
CREATE TRIGGER set_updated_at_timestamp_Oauth_table BEFORE UPDATE ON "Oauth" FOR EACH ROW EXECUTE PROCEDURE set_updated_at_column();
CREATE TRIGGER set_updated_at_timestamp_Team_table BEFORE UPDATE ON "Team" FOR EACH ROW EXECUTE PROCEDURE set_updated_at_column();
CREATE TRIGGER set_updated_at_timestamp_Permission_table BEFORE UPDATE ON "Permission" FOR EACH ROW EXECUTE PROCEDURE set_updated_at_column();
CREATE TRIGGER set_updated_at_timestamp_Task_table BEFORE UPDATE ON "Task" FOR EACH ROW EXECUTE PROCEDURE set_updated_at_column();

COMMIT;