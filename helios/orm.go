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

func (o *ORM) LastInsertId(table, key string) (id int64, err error) {
  sql := "SELECT MAX(" + key + ") FROM "  + table + ";"
  stmt, err := o.db.Query(sql)
  defer stmt.Close()

  if err != nil {
    return 0, err
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
  defer row.Close()

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
  defer row.Close()

  if err != nil {
    return nil, err
  }

  e := Event{}
  row.Next()
  row.Scan(&e.EventId, &e.Name, &e.ExpectedFrequency)
  row.Close()

  return &e, nil
}

func (o *ORM) SaveEvent(event *Event) (err error) {
  update  := event.EventId != 0

  if update { // existing event
    _, err = o.db.Exec( `UPDATE events
                         SET name = ?, expected_frequency = ?
                         WHERE event_id = ?;`,
                         event.Name, event.ExpectedFrequency, event.EventId)
  } else {
    _, err = o.db.Exec( `INSERT INTO events (name, expected_frequency)
                         VALUES (?, ?);`,
                         event.Name, event.ExpectedFrequency)
  }

  if err != nil { return err }

  if !update {
    event.EventId, _ = o.LastInsertId("events", "event_id")
  }

  return nil
}


func (o *ORM) SaveEventOccurrence(occurrence *EventOccurrence) (err error) {
  update  := occurrence.EventOccurrenceId != 0

  if update { // existing occurrence
    _, err = o.db.Exec( `UPDATE event_occurrences
                         SET    event_id = ?, occurred_at = ?
                         WHERE  event_occurrence_id = ?;`,
                         occurrence.EventId, occurrence.OccurredAt, occurrence.EventOccurrenceId)
  } else {
    _, err = o.db.Exec( `INSERT INTO event_occurrences (event_id, occurred_at)
                         VALUES (?, ?)`,
                         occurrence.EventId, occurrence.OccurredAt)
  }

  if err != nil { return err }

  if !update {
    occurrence.EventOccurrenceId, err = o.LastInsertId("event_occurrences", "event_occurrence_id")
  }

  return err
}
