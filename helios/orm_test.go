package helios

import (
  "testing"
  "database/sql"
  "os"
  "time"
  _ "github.com/mattn/go-sqlite3"
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
  defer doCleanup()

  event := &Event{Name: "test/metric", ExpectedFrequency: "* * * * *"}
  err := orm.SaveEvent(event)

  if err != nil {
    t.Log("Error raised on creating event:", err)
    t.Fail()
  }

  if event == nil {
    t.Log("Event was not loaded! Null pointer!")
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

  stmt, err := db.Query("SELECT * FROM events WHERE event_id = ?", event.EventId)
  defer stmt.Close()
  if err != nil {
    t.Log("Was not able to retrieve the created event; err:", err)
    t.Fail()
  }
}

func TestEventLoad(t *testing.T) {
  doSetup()
  defer doCleanup()

  _, err := db.Exec(
    `INSERT INTO events (name, expected_frequency)
     VALUES (?, ?)`, "asdf/foo", "asdf")

  if err != nil {
    t.Log("Error when creating event:", err)
    t.Fail()
  }

  id, _ := orm.LastInsertId("events", "event_id")
  event, err := orm.LoadEvent(id)

  if err != nil {
    t.Log("Error in loading event:", err)
    t.Fail()
  }

  if event.Name != "asdf/foo" || event.ExpectedFrequency != "asdf" {
    t.Log("Columns mapped incorrectly:", event.Name, event.ExpectedFrequency)
    t.Fail()
  }
}

func TestEventList(t *testing.T) {
  doSetup()
  defer doCleanup()

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
}

func TestEventUpdate(t *testing.T) {
  doSetup()
  defer doCleanup()

  tx, _ := db.Begin()
  tx.Exec(`INSERT INTO events (name, expected_frequency)
           VALUES (?, ?)`, "asdf/foo", "asdf")
  tx.Commit()

  id, _ := orm.LastInsertId("events", "event_id")
  event, err := orm.LoadEvent(id)

  if err != nil { t.Log(err); t.Fail() }

  event.Name = "asdf/baz"
  err = orm.SaveEvent(event)

  if err != nil {
    t.Log("Error when updating event:", err)
    t.Fail()
  }

  stmt, err := db.Query(`SELECT name FROM events WHERE event_id = ?`, id)
  defer stmt.Close()
  stmt.Next()
  var test string
  stmt.Scan(&test)

  if test != "asdf/baz" {
    t.Log("Name did not update; is currently", test, "instead of asdf/baz")
    t.Fail()
  }
}

func TestEventOccurrences(t *testing.T) {
  doSetup()
  defer doCleanup()

  event := &Event{Name: "test/metric", ExpectedFrequency: "* * * * *"}
  err := orm.SaveEvent(event)
  if (err != nil) { t.Log(err); t.Fail(); return }

  occurrence := &EventOccurrence{EventId: event.EventId, OccurredAt: time.Now()}
  err = orm.SaveEventOccurrence(occurrence)

  if (err != nil) { t.Log(err); t.Fail(); return }

  if occurrence.EventOccurrenceId == 0 {
    t.Log("EventOccurenceId was not attached to object after saving; occ:", occurrence)
    t.Fail()
    return
 }

  rows, err := db.Query("SELECT * FROM event_occurrences WHERE event_id = ?", event.EventId);
  defer rows.Close()
  if !rows.Next() { // error raised; no results
    t.Log("Object was not inserted into database for real")
    t.Fail()
  }
}
