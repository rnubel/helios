package helios

import (
  "database/sql"
)

// Initializes the Helios schema.
func LoadDatabaseSchema(db *sql.DB) error {
  sql := `
    CREATE TABLE events (
      event_id            INTEGER                     PRIMARY KEY AUTOINCREMENT,
      name                TEXT                        NOT NULL,
      expected_frequency  TEXT                        NOT NULL
    );
    `
  _, err := db.Exec(sql)
  if err != nil { return err }

  sql = `
    CREATE TABLE thresholds (
      threshold_id        INTEGER                     PRIMARY KEY AUTOINCREMENT,
      event_id            INTEGER                     NOT NULL,
      name                TEXT                        NOT NULL,
      severity_limit      NUMERIC(10,2)               NOT NULL,
      severity_function   TEXT                        NOT NULL
    );
    `
  _, err = db.Exec(sql)
  if err != nil { return err }

  sql = `
    CREATE TABLE event_occurrences (
      event_occurrence_id INTEGER                     PRIMARY KEY AUTOINCREMENT,
      event_id            INTEGER                     NOT NULL,
      occurred_at         TIMESTAMP NOT NULL
    );
    `
  _, err = db.Exec(sql)
  if err != nil { return err }

  sql = `
    CREATE TABLE threshold_checks (
      threshold_check_id  INTEGER                     PRIMARY KEY AUTOINCREMENT,
      threshold_id        INTEGER                     NOT NULL,
      checked_at          TIMESTAMP WITHOUT TIME ZONE NOT NULL,
      severity            NUMERIC(10,2)               NOT NULL
    );
  `

  _, err = db.Exec(sql)

  return err
}
