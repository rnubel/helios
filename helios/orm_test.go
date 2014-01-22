package helios

import (
  "testing"
  "database/sql"
  "os"
  _ "code.google.com/p/go-sqlite/go1/sqlite3"
)

var (
  db        *sql.DB
  orm       *ORM
)

func doSetup() {
  dbName  := "file:helios_test.db?cache=shared&mode=rwc"
  db, _   = sql.Open("sqlite3", dbName)
  orm   = NewORM(db)

  LoadDatabaseSchema(db)
}

func doCleanup() {
  db.Close()
  os.Remove("helios_test.db")
}

func TestEventCreate(t *testing.T) {
  doSetup()

  event := &Event{Name: "test/metric", ExpectedFrequency: "* * * * *"}
  err := orm.SaveEvent(event)

  if err != nil {
    t.Log("Error raised on creating event:", err)
    t.Fail()
  }

  if event == nil {
    t.Log("Event was not loaded! Null pointer!")
    t.Fail()
    return
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

func TestEventLoad(t *testing.T) {
  doSetup()

  _, err := db.Exec(
    `INSERT INTO events (name, expected_frequency)
     VALUES (?, ?)`, "asdf/foo", "asdf")

  if err != nil {
    t.Fail()
  }

  id := orm.LastInsertId("events", "event_id")
  event, err := orm.LoadEvent(id)

  if err != nil {
    t.Log("Error in loading event:", err)
    t.Fail()
  }

  if event.Name != "asdf/foo" || event.ExpectedFrequency != "asdf" {
    t.Log("Columns mapped incorrectly: ", event.Name, event.ExpectedFrequency)
    t.Fail()
  }

  doCleanup()
}

func TestEventList(t *testing.T) {
  doSetup()

  db.Exec(
    `INSERT INTO events (name, expected_frequency)
     VALUES (?, ?), (?, ?)`, "asdf/foo", "asdf", "asdf/bar", "basd")

  events, err := orm.ListEvents()

  if err != nil {
    t.Log("Error in loading events:", err)
    t.Fail()
  }

  if len(events) != 2 {
    t.Log("Wrong events loaded:", events)
    for i, e := range(events) { t.Log(i, *e) }
    t.Fail()
  }

  if events[0].Name != "asdf/foo" || events[1].Name != "asdf/bar" {
    // DBMS-dependent test, unless I add an ORDER BY.
    t.Log("Events loaded in wrong order:", events)
    t.Fail()
  }

  doCleanup()
}

func TestEventUpdate(t *testing.T) {
  doSetup()

  db.Exec(`INSERT INTO events (name, expected_frequency)
           VALUES (?, ?)`, "asdf/foo", "asdf")

  id := orm.LastInsertId("events", "event_id")
  event, err := orm.LoadEvent(id)

  if err != nil { t.Log(err); t.Fail() }

  event.Name = "asdf/baz"
  err = orm.SaveEvent(event)

  if err != nil {
    t.Log("Error when updating event:", err)
    t.Fail()
  }

  stmt, err := db.Query(`SELECT name FROM events WHERE event_id = ?`, id)
  stmt.Next()
  var test string
  stmt.Scan(&test)

  if test != "asdf/baz" {
    t.Log("Name did not update; is currently", test, "instead of asdf/baz")
    t.Fail()
  }

  doCleanup()
}
