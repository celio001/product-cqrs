-- Modify "brands" table
ALTER TABLE "public"."brands" ADD COLUMN "updated_at" timestamptz NOT NULL DEFAULT now();
-- Modify "categories" table
ALTER TABLE "public"."categories" ADD COLUMN "updated_at" timestamptz NOT NULL DEFAULT now();
-- Modify "product_fiscal_data" table
ALTER TABLE "public"."product_fiscal_data" ADD COLUMN "updated_at" timestamptz NOT NULL DEFAULT now();
