# Race Events Schema

## Decision
V1 race event annotation data is defined by [`schemas/race-events-v1.schema.json`](/Users/lorenbrindze/Projects/topaz-racing/schemas/race-events-v1.schema.json:1).

Each race bundle keeps editorial race moments in `events.json` so both `course` and `replay` embeds can highlight notable moments without coupling annotations to boat-track internals.

## File Shape

```json
{
  "events": [
    {
      "id": "weather-mark-set",
      "type": "gybe",
      "time": "2025-02-11T21:05:00Z",
      "lat": 34.4475,
      "lon": -120.0344,
      "label": "Set up for the run home",
      "description": "Topaz gybes after rounding Bishop Rock and turns back toward Santa Barbara."
    }
  ]
}
```

## Event Contract
- `events` is the only top-level collection and contains every annotation for the race.
- `id` and `type` are required on every event.
- Each event must be anchored by `time`, `lat`/`lon`, or both.
- `type` is limited to `gybe`, `wipeout`, `sail_change`, `tack`, `mark_rounding`, or `note`.
- `label` and `description` are optional editorial text fields.

## Anchoring Model
Time-anchored events support replay storytelling, while position-anchored events support persistent map callouts. Allowing both on one event lets the same annotation work in the static course view and during replay.

V1 does not require every event to know everything. A pure timing note can omit coordinates, and a pure chart callout can omit time, as long as one anchor style is present.

## Consequences
- A race can highlight moments like gybes, wipeouts, and sail changes with one stable payload shape.
- Replay UI can place annotations on the race clock without inferring event timing from prose.
- Later schema revisions can add richer event media or categories without redefining the core anchor model.
