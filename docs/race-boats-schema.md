# Race Boats Schema

## Decision
V1 race boat-track data is defined by [`schemas/race-boats-v1.schema.json`](/Users/lorenbrindze/Projects/topaz-racing/schemas/race-boats-v1.schema.json:1).

Each race bundle keeps all replay boats in one `boats.json` file so the embed can load a single track payload for both the self boat and competitors.

## File Shape

```json
{
  "boats": [
    {
      "id": "topaz",
      "name": "Topaz",
      "color": "#4fd1ff",
      "boatType": "Express 27",
      "source": "hand-authored",
      "isSelf": true,
      "track": [
        {
          "time": "2025-02-11T18:00:00Z",
          "lat": 34.4086,
          "lon": -119.6932
        },
        {
          "time": "2025-02-11T19:15:00Z",
          "lat": 34.4310,
          "lon": -119.8940
        }
      ]
    },
    {
      "id": "wildcard",
      "name": "Wildcard",
      "color": "#ff8a5b",
      "boatType": "J/105",
      "source": "hand-authored",
      "isSelf": false,
      "track": [
        {
          "time": "2025-02-11T18:00:00Z",
          "lat": 34.4092,
          "lon": -119.6921
        },
        {
          "time": "2025-02-11T19:10:00Z",
          "lat": 34.4347,
          "lon": -119.8842
        }
      ]
    }
  ]
}
```

## Boat Contract
- `boats` is the only top-level collection and contains every replay participant for the race.
- `id`, `name`, `color`, `boatType`, `source`, `isSelf`, and `track` are required on every boat.
- `track` is an ordered time series and each point requires `time`, `lat`, and `lon`.
- `time` uses RFC 3339 / JSON Schema `date-time` strings so replay state can normalize boats onto one clock.

## File Organization
V1 keeps all boats for a race in a single `boats.json` file instead of splitting one file per boat. That keeps replay loading simple, lets the shared engine calculate global time bounds from one payload, and stays independent of any future import pipeline.

The model does not depend on Jibeset, Garmin, GPX, or any other upstream source. Later import tooling can map those sources into this one stable internal contract.

## Consequences
- One self boat and multiple competitors can be rendered from the same payload shape.
- Replay code can compute time bounds and interpolation inputs without source-specific branching.
- A later schema revision can add optional telemetry fields without redefining the core per-boat track model.
