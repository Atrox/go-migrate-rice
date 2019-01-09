# Rice Source Driver for golang-migrate/migrate

## Installation
```bash
go get github.com/atrox/go-migrate-rice
```

## Example
```golang
import (
	rice "github.com/GeertJohan/go.rice"
	migraterice "github.com/atrox/go-migrate-rice"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
)

// get rice box
migrationsBox := rice.MustFindBox("migrations")

// create source driver with specified rice box
sourceDriver, _ := migraterice.WithInstance(migrationsBox)

// create database driver
dbDriver, _ := postgres.WithInstance(openDatabase(), &postgres.Config{})

// create migration instance
m, _ := migrate.NewWithInstance("rice", sourceDriver, "postgres", dbDriver)

// migrate
m.Up()
```
