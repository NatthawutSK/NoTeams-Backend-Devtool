BEGIN;

INSERT INTO "User" ("username", "password", "email", "dob", "phone", "bio")
VALUES
  ('user1', '$2a$10$agrNUd.FI/ZLzs2xfVbpR.VV/E08UWTYvWE1R6WPknAJbygG6ifLS', 'user1@example.com', '1990-01-01', '0882345678', 'I am user1'),
  ('user2', '$2a$10$agrNUd.FI/ZLzs2xfVbpR.VV/E08UWTYvWE1R6WPknAJbygG6ifLS', 'user2@example.com' , '1990-01-01', '0862345678', 'I am user2'),
  ('user3', '$2a$10$agrNUd.FI/ZLzs2xfVbpR.VV/E08UWTYvWE1R6WPknAJbygG6ifLS', 'user3@example.com' , '1990-01-01', '0822345678', 'I am user3');



COMMIT;
