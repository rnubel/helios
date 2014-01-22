package helios

import (
  "time"
)

type Event struct {
  EventId           int64
  Name              string
  ExpectedFrequency string
}

type Threshold struct {
  ThresholdId       int64
  EventId           int64
  Name              string
  SeverityLimit     float32
  SeverityFunction  string
}

type EventOccurrence struct {
  EventOccurrenceId int64
  EventId           int64
  OccurredAt        time.Time
}

type ThresholdCheck struct {
  ThresholdCheckId  int64
  ThresholdId       int64
  CheckedAt         time.Time
  Severity          float32
}
