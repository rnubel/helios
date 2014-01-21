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

func (o *ORM) CreateEvent(name, expectedFrequency string) (*Event, error) {
  stmt, err := o.db.Query(
    `INSERT INTO events (name, expected_frequency)
     VALUES (?, ?)`,
     name, expectedFrequency)

  if err != nil {
    return nil, err
  }

  eventId := db.Driver().LastInsertId();
  return &e, nil
}
