CREATE TABLE "users" (
  "username" varchar PRIMARY KEY,
  "email" varchar UNIQUE NOT NULL,
  "hashed_password" varchar NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "tasks" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "title" varchar NOT NULL,
  "description" varchar NOT NULL,
  "done" boolean NOT NULL DEFAULT 'FALSE',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "tasks" ("owner");

ALTER TABLE "tasks" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");
