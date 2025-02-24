package migrator

import (
	"flag"
)

func main(){
	var storagePath, migrationsPath, migrationsTable string

	flag.StringVar(&storagePath, "storage-path", "/tmp/storage", "Path to store the migrations")
	flag.StringVar(&migrationsPath, "migrations-path", "/tmp/migrations", "Path to store the migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "Table name for the migrations")
	flag.Parse()

	if storagePath == "" {
		panic("storage-path is required")
	}
	if migrationsPath == "" {
		panic("migrations-path is required")
	}
	m, err :=
}