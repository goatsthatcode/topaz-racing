# Race Visualization Usage

The existing `race-viz` component is page-bundle driven.

Author a race as a Hugo leaf bundle:

```text
content/races/my-season/my-race/
  index.md
  course.json
  boats.json
  events.json
```

`index.md` owns the prose and embeds the component:

```md
+++
title = "My Race"
date = "2026-05-11T09:00:00-07:00"
+++

{{< race-viz >}}
```

You can also embed a race bundle from another page by race slug:

```md
{{< race-viz race="dan-byrne-2025/bishop-rock-race" mode="replay" title="Bishop Rock Replay" class="race-viz-feature" >}}
```

Supported shortcode parameters:

- `race`: optional race bundle path under `content/races/`; omit it when embedding from the race page itself
- `mode`: `replay` by default; use `course` to render the chart and course without replay controls
- `title`: optional accessible title override for the stage label
- `class`: optional extra class on the outer `<figure>`

Example `course.json`:

```json
{
  "id": "my-race-course",
  "name": "My Race",
  "elements": [
    { "id": "start", "type": "start_line", "lat": 33.9769, "lon": -118.4451, "name": "Start", "rounding": "none" },
    { "id": "weather", "type": "mark", "lat": 33.82, "lon": -118.61, "name": "Weather Mark", "rounding": "port" },
    { "id": "finish", "type": "finish_line", "lat": 33.74, "lon": -118.32, "name": "Finish", "rounding": "none" }
  ]
}
```

Example `boats.json`:

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
        { "time": "2026-05-11T16:00:00Z", "lat": 33.9769, "lon": -118.4451 },
        { "time": "2026-05-11T17:00:00Z", "lat": 33.82, "lon": -118.61 },
        { "time": "2026-05-11T18:00:00Z", "lat": 33.74, "lon": -118.32 }
      ]
    }
  ]
}
```

Example `events.json`:

```json
{
  "events": [
    {
      "id": "set-up-run",
      "type": "gybe",
      "time": "2026-05-11T17:20:00Z",
      "lat": 33.8,
      "lon": -118.57,
      "label": "Set up for the run",
      "description": "Topaz gybes onto port for the run home."
    }
  ]
}
```

Local preview uses the same component contract as production. If `martin` is not running at the configured development tile endpoint, the component now falls back to the production tile host so the map can still render in `hugo server`.
