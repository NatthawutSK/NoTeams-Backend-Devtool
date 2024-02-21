BEGIN;

INSERT INTO "User" ("username", "password", "email", "dob", "phone", "bio", "avatar")
VALUES
  ('user1', '$2a$10$agrNUd.FI/ZLzs2xfVbpR.VV/E08UWTYvWE1R6WPknAJbygG6ifLS', 'user1@example.com', '1990-01-01', '0882345678', 'I am user1', 'https://www.gravatar.com/avatar/205e460b479e2e5b48aec07710c08d50'),
  ('user2', '$2a$10$agrNUd.FI/ZLzs2xfVbpR.VV/E08UWTYvWE1R6WPknAJbygG6ifLS', 'user2@example.com' , '1990-01-01', '0862345678', 'I am user2', 'https://www.gravatar.com/avatar/205e460b479e2e5b48aec07710c08d50'),
  ('user3', '$2a$10$agrNUd.FI/ZLzs2xfVbpR.VV/E08UWTYvWE1R6WPknAJbygG6ifLS', 'user3@example.com' , '1990-01-01', '0822345678', 'I am user3', 'https://www.gravatar.com/avatar/205e460b479e2e5b48aec07710c08d50');



COMMIT;
