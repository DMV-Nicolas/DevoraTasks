CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "hashed_password" varchar NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "tasks" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "title" varchar NOT NULL,
  "description" varchar,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "tasks" ("user_id");

ALTER TABLE "tasks" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
