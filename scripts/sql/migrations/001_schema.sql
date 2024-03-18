-- +migrate Up

CREATE TABLE "users" (
    "id" uuid DEFAULT gen_random_uuid(),
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "name" varchar(256),
    "email" varchar(512),
    "password" varchar(256),
    "gender" gender DEFAULT gender('unknown'),
    PRIMARY KEY ("id")
);

CREATE TABLE "types" (
    "id" uuid DEFAULT gen_random_uuid(),
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "name" varchar(64),
    PRIMARY KEY ("id")
);

CREATE TABLE "subtypes" (
    "id" uuid DEFAULT gen_random_uuid(),
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "name" varchar(64),
    PRIMARY KEY ("id")
);

CREATE TABLE "styles" (
    "id" uuid DEFAULT gen_random_uuid(),
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "name" varchar(64),
    PRIMARY KEY ("id")
);

CREATE TABLE "clothes" (
    "id" uuid DEFAULT gen_random_uuid(),
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "name" varchar(128),
    "note" varchar(512),
    "image" varchar(256),
    "user_id" uuid,
    "style_id" uuid DEFAULT null,
    "type_id" uuid DEFAULT null,
    "subtype_id" uuid DEFAULT null,
    "color" char(7),
    "seasons" season [],
    PRIMARY KEY ("id"),
    CONSTRAINT "fk_clothes_user" FOREIGN KEY ("user_id") REFERENCES "users"("id"),
    CONSTRAINT "fk_clothes_style" FOREIGN KEY ("style_id") REFERENCES "styles"("id"),
    CONSTRAINT "fk_clothes_type" FOREIGN KEY ("type_id") REFERENCES "types"("id"),
    CONSTRAINT "fk_clothes_subtype" FOREIGN KEY ("subtype_id") REFERENCES "subtypes"("id")
);

CREATE TABLE "tags" (
    "id" uuid DEFAULT gen_random_uuid(),
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "name" varchar(64),
    PRIMARY KEY ("id")
);

CREATE TABLE "clothes_tags" (
    "clothes_model_id" uuid DEFAULT gen_random_uuid(),
    "tag_id" uuid DEFAULT gen_random_uuid(),
    PRIMARY KEY ("clothes_model_id", "tag_id"),
    CONSTRAINT "fk_clothes_tags_clothes_model" FOREIGN KEY ("clothes_model_id") REFERENCES "clothes"("id"),
    CONSTRAINT "fk_clothes_tags_tag" FOREIGN KEY ("tag_id") REFERENCES "tags"("id")
);

CREATE TABLE "user_images" (
    "id" uuid DEFAULT gen_random_uuid(),
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "user_id" text,
    "image" text,
    PRIMARY KEY ("id")
);

CREATE TABLE "try_on_results" (
    "id" uuid DEFAULT gen_random_uuid(),
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "image" text,
    "rating" bigint,
    "user_id" uuid,
    "clothes_model_id" uuid,
    PRIMARY KEY ("id"),
    CONSTRAINT "fk_try_on_results_user" FOREIGN KEY ("user_id") REFERENCES "users"("id"),
    CONSTRAINT "fk_try_on_results_clothes_model" FOREIGN KEY ("clothes_model_id") REFERENCES "clothes"("id")
);

-- +migrate Down
DROP TABLE "users";

DROP TABLE "types";

DROP TABLE "subtypes";

DROP TABLE "styles";

DROP TABLE "clothes";

DROP TABLE "tags";

DROP TABLE "clothes_tags";

DROP TABLE "user_images";

DROP TABLE "try_on_results";
