env "dev" {
    url = "postgres://postgres:postgres@localhost:5432/product?sslmode=disable"
    dev = "postgres://atlas:atlas@localhost:5433/atlas?sslmode=disable"

    migration {
        dir = "file://./internal/database/migrations"
        format = "atlas"
    }

    schema{
        src = ["internal/database/schema/schema.sql"]
    }
}