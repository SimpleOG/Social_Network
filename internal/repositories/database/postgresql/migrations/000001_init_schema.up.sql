CREATE TABLE "users" (
                        "id" SERIAL NOT NULL UNIQUE primary key ,
                        "username" VARCHAR(255) UNIQUE NOT NULL ,

                        "password" VARCHAR(255)  NOT NULL,
                        "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX "user_index_0"
    ON "users" ("username");

