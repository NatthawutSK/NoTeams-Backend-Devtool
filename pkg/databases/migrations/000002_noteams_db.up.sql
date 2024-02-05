BEGIN;

INSERT INTO "User" ("username", "password", "email")
VALUES
  ('user1', '$2a$10$agrNUd.FI/ZLzs2xfVbpR.VV/E08UWTYvWE1R6WPknAJbygG6ifLS', 'user1@example.com');


COMMIT;
