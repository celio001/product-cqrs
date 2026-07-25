package config

import "os"

var config = map[string]string{
	"ENV": "development",

	//Postgres
	"POSTGRES_DB_DSN": "",
}

func GetString(k string) string {
	v := os.Getenv(k)
	if v == "" {
		return config[k]
	}
	return v
}
