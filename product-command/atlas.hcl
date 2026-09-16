env "local" {
    url = "postgres://postgres:postgres@postgres-main:5432/product?sslmode=disable"
    dev = "postgres://atlas:atlas@postgres-dev:5433/atlas?sslmode=disable"

    migration {
        dir = "file://./internal/database/migrations"
        format = "atlas"
    }

    schema{
        src = ["internal/database/schema/schema.sql"]
    }
}