package helios

import (
  "testing"
  "database/sql"
  _ "github.com/mattn/go-sqlite3"
)

func createTestDatabase() (*sql.DB, error) {
  return sql.Open("sqlite3", ":memory:")
}

func TestSchemaLoad(t *testing.T) {
  db, err := createTestDatabase();

  if (err != nil) {
    t.Log(err)
    t.Fail()
  }

  err = LoadDatabaseSchema(db);
  if (err != nil) {
    t.Log(err)
    t.Fail()
  }

  // Test a successful insert sequence
  _, err = db.Exec(`
    INSERT INTO events (name, expected_frequency)                               VALUES ('helios/events/test', '* * * * *');
    INSERT INTO event_occurrences (event_id, occurred_at)                       VALUES (1, current_timestamp);
    INSERT INTO thresholds (event_id, name, severity_limit, severity_function)  VALUES (1, 'helios/thresholds/test', 0.9, 'quadratic');
    INSERT INTO threshold_checks (threshold_id, checked_at, severity)           VALUES (1, current_timestamp, 0.54);
  `)
  if (err != nil) {
    t.Log("Expected err to be nil, but was:", err)
    t.Fail()
  }

  // Test some of the not-null constraints
  _, err = db.Exec("INSERT INTO threshold_checks (threshold_id) VALUES (NULL)")
  if (err == nil) {
    t.Log("Expected err to not be nil, but was!")
    t.Fail()
  }

  _, err = db.Exec("INSERT INTO threshold_checks (threshold_id, checked_at) VALUES (1, NULL)")
  if (err == nil) {
    t.Log("Expected err to not be nil, but was!")
    t.Fail()
  }

  _, err = db.Exec("INSERT INTO threshold_checks (threshold_id, checked_at, severity) VALUES (1, current_timestamp, NULL)")
  if (err == nil) {
    t.Log("Expected err to not be nil, but was!")
    t.Fail()
  }

  // Test uniqueness
  _, err = db.Exec("INSERT INTO events (name, expected_frequency)                               VALUES ('helios/events/test', '* * * * *');")
  if (err != nil) {
    t.Log("Duplicate names allowed in events table!")
    t.Fail()
  }

  db.Close()
}


