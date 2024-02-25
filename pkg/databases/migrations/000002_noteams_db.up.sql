BEGIN;

INSERT INTO "User" ("username", "password", "email", "dob", "phone", "bio", "avatar")
VALUES
  ('user1', '$2a$10$agrNUd.FI/ZLzs2xfVbpR.VV/E08UWTYvWE1R6WPknAJbygG6ifLS', 'user1@example.com', '1990-01-01', '0882345678', 'I am user1', 'https://www.gravatar.com/avatar/205e460b479e2e5b48aec07710c08d50'),
  ('user2', '$2a$10$agrNUd.FI/ZLzs2xfVbpR.VV/E08UWTYvWE1R6WPknAJbygG6ifLS', 'user2@example.com' , '1990-01-01', '0862345678', 'I am user2', 'https://www.gravatar.com/avatar/205e460b479e2e5b48aec07710c08d50'),
  ('user3', '$2a$10$agrNUd.FI/ZLzs2xfVbpR.VV/E08UWTYvWE1R6WPknAJbygG6ifLS', 'user3@example.com' , '1990-01-01', '0822345678', 'I am user3', 'https://www.gravatar.com/avatar/205e460b479e2e5b48aec07710c08d50');

-- Mock data for Team table
INSERT INTO "Team" ("team_name", "team_desc", "team_code", "team_poster") 
VALUES 
('Development Team', 'A team focused on software development', 'DEVTEAM', ''),
('Marketing Team', 'Responsible for marketing activities', 'MARKTEAM', ''),
('Design Team', 'Handles design-related tasks', 'DESIGNTEAM', '');

-- Mock data for TeamMember table
INSERT INTO "TeamMember" ("user_id", "team_id", "role")
VALUES
('U000001', 'T000001', 'OWNER'),
('U000002', 'T000001', 'MEMBER'),
('U000003', 'T000002', 'OWNER'),
('U000001', 'T000002', 'MEMBER'),
('U000002', 'T000003', 'OWNER'),
('U000003', 'T000003', 'MEMBER');

-- Mock data for Permission table
INSERT INTO "Permission" ("allow_task", "allow_file", "allow_invite", "team_id")
VALUES
(TRUE, TRUE, TRUE, 'T000001'),
(TRUE, FALSE, TRUE, 'T000002'),
(FALSE, TRUE, FALSE, 'T000003');

-- Mock data for File table
INSERT INTO "File" ("file_name", "file_url", "team_id", "user_id")
VALUES
('File 1', 'https://www.example.com/file1', 'T000001', 'U000001'),
('File 2', 'https://www.example.com/file2', 'T000002', 'U000002'),
('File 3', 'https://www.example.com/file3', 'T000003', 'U000003');

-- Mock data for Task table
INSERT INTO "Task" ("task_name", "task_desc", "task_status", "task_deadline", "team_id", "user_id")
VALUES
('Task 1', 'Complete task 1', 'TODO', '2024-02-28', 'T000001', 'U000001'),
('Task 2', 'Complete task 2', 'DOING', '2024-03-05', 'T000001', 'U000002'),
('Task 3', 'Complete task 3', 'DONE', '2024-02-25', 'T000002', 'U000003'),
('Task 4', 'Complete task 4', 'TODO', '2024-03-10', 'T000003', 'U000001');



COMMIT;
