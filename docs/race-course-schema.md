# Race Course Schema

## Decision
V1 race course data is defined by [`schemas/race-course-v1.schema.json`](/Users/lorenbrindze/Projects/topaz-racing/schemas/race-course-v1.schema.json:1).

Each race bundle keeps its course in `course.json`, with one ordered `elements` array as the canonical course definition for both `course` and `replay` embeds.

## File Shape

```json
{
  "id": "bishop-rock-race-course",
  "name": "Bishop Rock Race",
  "elements": [
    {
      "id": "start",
      "type": "start_line",
      "lat": 34.4086,
      "lon": -119.6932,
      "name": "Santa Barbara Start",
      "rounding": "none"
    },
    {
      "id": "bishop-rock",
      "type": "mark",
      "lat": 34.4562,
      "lon": -120.1198,
      "name": "Bishop Rock",
      "rounding": "port",
      "controlPointsToNext": [
        { "lat": 34.4481, "lon": -120.0213 }
      ]
    },
    {
      "id": "finish",
      "type": "finish_line",
      "lat": 34.4086,
      "lon": -119.6932,
      "name": "Santa Barbara Finish",
      "rounding": "none"
    }
  ]
}
```

## Element Contract
- `id`, `type`, `lat`, `lon`, and `rounding` are required on every course element.
- `name` is optional and exists for labels in the editorial UI.
- `type` is limited to `mark`, `start_line`, or `finish_line`.
- `rounding` is limited to `port`, `starboard`, or `none`.

The array order is significant. The visualization engine connects each element to the next element in sequence when it renders the course geometry.

## Manual Route Shaping
`controlPointsToNext` is optional on any element. When present, those coordinates are inserted between the current element and the next ordered element to shape the rendered leg around land or other obstacles.

The reference Bishop Rock bundle currently leaves `controlPointsToNext` empty so the rendered leg stays straight between the authored marks.

This keeps course marks as the author-facing model while leaving room for future automatic route-around-land logic. A later schema revision can add richer routing metadata without changing the main ordered course element list.

## Consequences
- A standalone `course` embed can render fully from `course.json` alone.
- `replay` mode can reuse the same course payload without inventing a second geometry format.
- Future exclusion zones can be added as adjacent top-level collections or a schema revision without redefining the core course path model.
