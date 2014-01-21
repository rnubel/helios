package helios

import (
  "time"
)

type Event struct {
  EventId           int
  Name              string
  ExpectedFrequency string
}

type Threshold struct {
  ThresholdId       int
  EventId           int
  Name              string
  SeverityLimit     float32
  SeverityFunction  string
}

type EventOccurrence struct {
  EventOccurrenceId int
  EventId           int
  OccurredAt        time.Time
}

type ThresholdCheck struct {
  ThresholdCheckId  int
  ThresholdId       int
  CheckedAt         time.Time
  Severity          float32
}
