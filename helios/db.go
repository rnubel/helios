package helios

import (
  "database/sql"
  _ "github.com/mattn/go-sqlite3"
)

func OpenDatabaseConnection() (*sql.DB, error) {
  return sql.Open("sqlite3", "helios.db")
}
