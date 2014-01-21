package helios

import (
  "testing"
  "database/sql"
  _ "code.google.com/p/go-sqlite/go1/sqlite3"
)

var db  *sql.DB
var orm *ORM

func doSetup() {
  db, _ = sql.Open("sqlite3", ":memory:")
  orm   = NewORM(db)
}

func doCleanup() {
  db.Close()
}

func TestEventCreate(t *testing.T) {
  doSetup()

  event, err := orm.CreateEvent("test/metric", "* * * * *")

  if err != nil {
    t.Log("Error raised on creating event:", err)
    t.Fail()
  }

  if event.Name != "test/metric" {
    t.Log("Name was not mapped correctly")
    t.Fail()
  }

  if event.ExpectedFrequency != "* * * * *" {
    t.Log("Frequency was not mapped correctly")
    t.Fail()
  }

  _, err = db.Query("SELECT * FROM events WHERE event_id = ?", event.EventId)
  if err != nil {
    t.Log("Was not able to retrieve the created event; err:", err)
    t.Fail()
  }

  doCleanup()
}
