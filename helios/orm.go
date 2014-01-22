package helios

import (
  "database/sql"
)

type ORM struct {
  db  *sql.DB
}

func NewORM(db *sql.DB) *ORM {
  return &ORM{db: db}
}

func (o *ORM) LastInsertId(table, key string) (id int64) {
  stmt, err := o.db.Query("SELECT MAX(" + key + ") FROM "  + table + ";")

  if err != nil {
    return 0
  }

  stmt.Next() // Needed, for some reason. Unsure why.
  stmt.Scan(&id)

  return
}

func (o *ORM) ListEvents() ([]*Event, error) {
  events := make([]*Event, 0)
  row, err := o.db.Query(`
    SELECT event_id, name, expected_frequency
    FROM events`);

  if err != nil {
    return events, err
  }

  for row.Next() {
    var e Event
    row.Scan(&e.EventId, &e.Name, &e.ExpectedFrequency)
    events = append(events, &e)
  }

  return events, nil
}

func (o *ORM) LoadEvent(eventId int64) (*Event, error) {
  row, err := o.db.Query(`
    SELECT event_id, name, expected_frequency
    FROM events
    WHERE event_id = ?
  `, eventId);

  if err != nil {
    return nil, err
  }

  e := Event{}
  row.Next()
  row.Scan(&e.EventId, &e.Name, &e.ExpectedFrequency)

  return &e, nil
}

func (o *ORM) CreateEvent(name, expectedFrequency string) (*Event, error) {
  _, err := o.db.Exec(
    `INSERT INTO events (name, expected_frequency)
     VALUES (?, ?)`,
     name, expectedFrequency)

  if err != nil {
    return nil, err
  }

  eventId := o.LastInsertId("events", "event_id")
  e, err := o.LoadEvent(eventId)

  return e, nil
}
