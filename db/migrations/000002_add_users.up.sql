CREATE TABLE "users"
(
    "username"           varchar(255)   PRIMARY KEY,
    "hashed_password"    varchar(255)   NOT NULL,
    "full_name"          text           NOT NULL,
    "email"              varchar UNIQUE NOT NULL,
    "password_change_at" timestamptz    NOT NULL DEFAULT '0001-01-01 00:00:00Z',
    "created_at"         timestamptz    NOT NULL DEFAULT 'now()'
);

ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

-- CREATE UNIQUE INDEX ON "accounts" ("owner", "currency");
-- this is the same approach but here we can provide naming (it easer to migrate down file to determine index name)
ALTER TABLE "accounts" ADD CONSTRAINT "owner_currency_key" UNIQUE ("owner", "currency")