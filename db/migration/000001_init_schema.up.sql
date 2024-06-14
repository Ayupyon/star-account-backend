CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "hashed_password" varchar NOT NULL,
  "create_time" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "accounts" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "create_time" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "records" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "type" integer NOT NULL,
  "date" date NOT NULL,
  "amount" numeric NOT NULL,
  "account_id" bigint NOT NULL,
  "create_user_id" bigint NOT NULL,
  "last_modified_user_id" bigint NOT NULL,
  "create_time" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "account_access_rules" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "account_id" bigint NOT NULL,
  "role" integer NOT NULL
);

CREATE INDEX ON "users" ("name");

CREATE INDEX ON "users" ("email");

CREATE INDEX ON "accounts" ("name");

CREATE INDEX ON "records" ("name");

CREATE INDEX ON "records" ("create_user_id");

CREATE INDEX ON "records" ("last_modified_user_id");

CREATE INDEX ON "records" ("account_id");

CREATE INDEX ON "account_access_rules" ("user_id");

CREATE INDEX ON "account_access_rules" ("account_id");

CREATE UNIQUE INDEX ON "account_access_rules" ("user_id", "account_id");

ALTER TABLE "records" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "account_access_rules" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "account_access_rules" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");
