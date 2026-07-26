-- 1. ENUM para o status do produto
CREATE TYPE product_status AS ENUM ('ACTIVE', 'INACTIVE', 'DISCONTINUED');

-- 2. Tabela de Categorias
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    parent_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP, 
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- 3. Tabela de Marcas
CREATE TABLE brands (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- 4. Tabela Produtos
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    brand_id UUID REFERENCES brands(id) ON DELETE SET NULL,
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    
    name VARCHAR(255) NOT NULL,
    sku VARCHAR(50) UNIQUE NOT NULL,
    barcode_ean13 VARCHAR(13),
    short_description VARCHAR(255),
    detailed_description TEXT,
    unit_of_measure VARCHAR(10) NOT NULL DEFAULT 'UN',
    
    cost_price DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    sale_price DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    promotional_price DECIMAL(10,2),
    
    gross_weight DECIMAL(8,3),
    net_weight DECIMAL(8,3),
    height DECIMAL(8,2),
    width DECIMAL(8,2),
    length DECIMAL(8,2),
    
    status product_status NOT NULL DEFAULT 'INACTIVE',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP, 
    deleted_at TIMESTAMPTZ
);

-- 5. Tabela de Estoque
CREATE TABLE product_inventory (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID REFERENCES products(id) ON DELETE CASCADE,
    location_aisle VARCHAR(50),
    quantity_available DECIMAL(10,3) NOT NULL DEFAULT 0.000,
    minimum_stock DECIMAL(10,3) NOT NULL DEFAULT 0.000,
    maximum_stock DECIMAL(10,3),
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- 6. Tabela de Dados Fiscais
CREATE TABLE product_fiscal_data (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID UNIQUE REFERENCES products(id) ON DELETE CASCADE,
    ncm_code VARCHAR(8),
    cest_code VARCHAR(7),
    origin_code SMALLINT NOT NULL DEFAULT 0, -- 0 para Nacional por padrão
    icms_rate DECIMAL(5,2) DEFAULT 0.00,
    pis_rate DECIMAL(5,2) DEFAULT 0.00,
    cofins_rate DECIMAL(5,2) DEFAULT 0.00,
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ipi_rate DECIMAL(5,2) DEFAULT 0.00
);