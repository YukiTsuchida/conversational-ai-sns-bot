package config

import "os"

func POSTGRES_HOST() string {
	return os.Getenv("POSTGRES_HOST")
}

func POSTGRES_PORT() string {
	return os.Getenv("POSTGRES_PORT")
}

func POSTGRES_USER() string {
	return os.Getenv("POSTGRES_USER")
}

func POSTGRES_PASSWORD() string {
	return os.Getenv("POSTGRES_PASSWORD")
}

func POSTGRES_DB() string {
	return os.Getenv("POSTGRES_DB")
}
