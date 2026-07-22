-- Create enum type "product_status"
CREATE TYPE "public"."product_status" AS ENUM ('ACTIVE', 'INACTIVE', 'DISCONTINUED');
-- Create "brands" table
CREATE TABLE "public"."brands" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(100) NOT NULL,
  "created_at" timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create "categories" table
CREATE TABLE "public"."categories" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "parent_id" uuid NULL,
  "name" character varying(100) NOT NULL,
  "created_at" timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "categories_parent_id_fkey" FOREIGN KEY ("parent_id") REFERENCES "public"."categories" ("id") ON UPDATE NO ACTION ON DELETE SET NULL
);
-- Create "products" table
CREATE TABLE "public"."products" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "brand_id" uuid NULL,
  "category_id" uuid NULL,
  "name" character varying(255) NOT NULL,
  "sku" character varying(50) NOT NULL,
  "barcode_ean13" character varying(13) NULL,
  "short_description" character varying(255) NULL,
  "detailed_description" text NULL,
  "unit_of_measure" character varying(10) NOT NULL DEFAULT 'UN',
  "cost_price" numeric(10,2) NOT NULL DEFAULT 0.00,
  "sale_price" numeric(10,2) NOT NULL DEFAULT 0.00,
  "promotional_price" numeric(10,2) NULL,
  "gross_weight" numeric(8,3) NULL,
  "net_weight" numeric(8,3) NULL,
  "height" numeric(8,2) NULL,
  "width" numeric(8,2) NULL,
  "length" numeric(8,2) NULL,
  "status" "public"."product_status" NOT NULL DEFAULT 'INACTIVE',
  "created_at" timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "products_sku_key" UNIQUE ("sku"),
  CONSTRAINT "products_brand_id_fkey" FOREIGN KEY ("brand_id") REFERENCES "public"."brands" ("id") ON UPDATE NO ACTION ON DELETE SET NULL,
  CONSTRAINT "products_category_id_fkey" FOREIGN KEY ("category_id") REFERENCES "public"."categories" ("id") ON UPDATE NO ACTION ON DELETE SET NULL
);
-- Create "product_fiscal_data" table
CREATE TABLE "public"."product_fiscal_data" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "product_id" uuid NULL,
  "ncm_code" character varying(8) NULL,
  "cest_code" character varying(7) NULL,
  "origin_code" smallint NOT NULL DEFAULT 0,
  "icms_rate" numeric(5,2) NULL DEFAULT 0.00,
  "pis_rate" numeric(5,2) NULL DEFAULT 0.00,
  "cofins_rate" numeric(5,2) NULL DEFAULT 0.00,
  "ipi_rate" numeric(5,2) NULL DEFAULT 0.00,
  PRIMARY KEY ("id"),
  CONSTRAINT "product_fiscal_data_product_id_key" UNIQUE ("product_id"),
  CONSTRAINT "product_fiscal_data_product_id_fkey" FOREIGN KEY ("product_id") REFERENCES "public"."products" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create "product_inventory" table
CREATE TABLE "public"."product_inventory" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "product_id" uuid NULL,
  "location_aisle" character varying(50) NULL,
  "quantity_available" numeric(10,3) NOT NULL DEFAULT 0.000,
  "minimum_stock" numeric(10,3) NOT NULL DEFAULT 0.000,
  "maximum_stock" numeric(10,3) NULL,
  "updated_at" timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  CONSTRAINT "product_inventory_product_id_fkey" FOREIGN KEY ("product_id") REFERENCES "public"."products" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
