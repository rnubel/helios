# Helios

If something in your application blows up, you probably already get notified (I hope!). But what about the more subtle
case of "negative" events, such as when a scheduled process fails to even start? Would you know if no customers had
signed up in the past twenty minutes? Helios is a heartbeat-based monitoring system, designed to receive events over
a messaging bus and gradually raise alarms as a system becomes more and more anomalous. It is meant to be hooked up to
a monitoring system like Zenoss to actually handle the alerts.

And since I haven't written any code yet, the rest of this is a design doc.

# Design thoughts

* Messaging bus (zeromq? ampq? HTTP POSTs? should be swappable)
* Store expected events in a database
* Cron parsing for "how often should events occur"
* Configurable thresholds for allowed variance (e.g. exponential, quadratic, linear severity)
* Maybe use frequency checks, too, not just time-since-last-occurrence
* Severity system ranges from 0 to 100?
* Need a filter for things like, not on holidays, etc.
* Periodically (every 5 minutes?) check every filter (round robin)
* API exposed via HTTP for Zenoss, etc, for monitoring
* API should be able to be per-event, or per-category (event categorization or tagging would be
  useful... maybe a heirarchy?), or global.
* Front-end display with little bars showing severity, a view of the polling schedule, admin interface (over an API)

## Event receivers

* Spawn off of main process upon starting up
* Use a channel to communicate received events back to main thread

Code:

    type EventReceiver interface {
        Run(events chan) chan // Takes a channel to send events back on, returns a channel to control it via.
    }

## Monitoring daemon

* Do we push out notifications, then display them via the API?
* Or is the API responsible for checking the thresholds?
* Can we store anomalies? Maybe we should store health checks in general, giving a time-lapse picture of the
  health of the system.
* The above sounds good. Then the API just returns the latest health check's result.
* So, Helios's main daemon needs to loop repeatedly, check threshold severities, and store the health check.

Code:

    package Helios

    type EventType struct {
        name            string      // not really, but roll with it
        expectedFreq    Frequency   // like a cron thing
    }

    type Threshold struct {
        name            string
        eventType       EventType
        maxVariance     time.Duration
    }

    func CheckThreshold(th Threshold) {
        db  := MagicDatabaseThing.New()
        sev := GetThresholdSeverity(db, th)
        StoreHealthCheck(db, th, sev)
    }

    func GetThresholdSeverity(db MagicDatabaseThing, th Threshold) (severity int) {
        severity    := 0
        t0          := db.GetLatestOccurrence(th.eventType)
        t1          := th.eventType.ClosestExpectedOccurrence(time.Now())

        if Math.abs(t1 - t0) > th.maxVariance {
            severity = 100; // not really
        else {
            severity = 0;   // ... well, maybe, if there's a binary threshold limit
        }

        return
    }

    func StoreHealthCheck(db MagicDatabaseThing, th Threshold, sev int) {
        db.InsertHealthCheck(th.name, th.eventType, time.Now(), sev)
    }

## API

* Pretty straightforward implementation
* API resources:
    * `/v1/event_types`
    * `/v1/thresholds`
    * `/v1/thresholds/:threshold/status`
    * `/v1/status`