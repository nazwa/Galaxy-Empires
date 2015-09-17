package middleware

import (
	"github.com/DavidHuie/gomigrate"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	DBKey string = "DB"
)

func RunMigrations(db *sqlx.DB) error {
	migrator, err := gomigrate.NewMigrator(db.DB, gomigrate.Mysql{}, "./migrations")
	if err != nil {
		return err
	}
	if err := migrator.Migrate(); err != nil {
		return err
	}
	return nil
}

// This method collects all the errors and submits them to Rollbar
func Database(connString string) gin.HandlerFunc {
	db := sqlx.MustConnect("mysql", connString)
	if err := RunMigrations(db); err != nil {
		panic(err)
	}

	return func(c *gin.Context) {
		c.Set(DBKey, db)
		c.Next()
	}
}
