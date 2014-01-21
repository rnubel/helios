package helios

import (
  "database/sql"
  _ "code.google.com/p/go-sqlite/go1/sqlite3"
)

func OpenDatabaseConnection() (*sql.DB, error) {
  return sqlite3.Open("helios.db")
}
