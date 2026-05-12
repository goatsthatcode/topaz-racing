# Race Authoring Workflow

This guide covers everything needed to publish a new race on the Topaz Racing site. No code changes are required — authoring is entirely file-based.

## Overview

Each race is a Hugo leaf bundle containing:
- `index.md` — prose, front matter, and the visualization embed
- `course.json` — the ordered race course (marks, start/finish lines)
- `boats.json` — one or more boat tracks with time-series positions
- `events.json` — notable race moments shown as map annotations

All four files live together in one directory under `content/races/`.

---

## Step 1: Create The Race Directory

Choose a season or series slug and a race slug. Both become part of the page URL.

```text
content/
  races/
    <season-slug>/
      _index.md          (one per season, already exists if adding to an existing season)
      <race-slug>/
        index.md
        course.json
        boats.json
        events.json
```

**Example**: a new race called "Channel Islands Overnight" in the 2026 Dan Byrne series would go in:

```text
content/races/dan-byrne-2026/channel-islands-overnight/
```

The race slug (`channel-islands-overnight`) becomes the canonical URL path and the race identifier used in cross-page embeds.

If the season directory is new, add an `_index.md` for it with front matter similar to the existing `dan-byrne-2025/_index.md`.

---

## Step 2: Write index.md

The page prose and the visualization embed both live here. Place the shortcode first so the map leads the content on desktop and mobile.

```markdown
+++
date = '2026-03-15T09:00:00-08:00'
title = 'Channel Islands Overnight'
+++

{{< race-viz >}}

Race recap prose goes here. The visualization above is populated from the
adjacent course.json, boats.json, and events.json files.
```

The `{{< race-viz >}}` shortcode (no arguments) embeds the visualization using the race data files in the same bundle directory.

To embed this race from a different page or post, use the explicit `race` parameter:

```markdown
{{< race-viz race="dan-byrne-2026/channel-islands-overnight" >}}
```

---

## Step 3: Author course.json

The course is an ordered list of course elements. The visualization connects them in sequence.

```json
{
  "id": "channel-islands-overnight-course",
  "name": "Channel Islands Overnight",
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
      "id": "anacapa",
      "type": "mark",
      "lat": 34.0133,
      "lon": -119.3631,
      "name": "Anacapa East End",
      "rounding": "port"
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

**Required fields on every element**: `id`, `type`, `lat`, `lon`, `rounding`.

**`type`** must be one of:
- `start_line` — the race start
- `mark` — a rounding mark mid-course
- `finish_line` — the race finish

**`rounding`** must be one of:
- `port` — round the mark leaving it to port
- `starboard` — round the mark leaving it to starboard
- `none` — for start and finish lines

**Optional `name`**: displayed as a label on the map.

### Routing Around Land

If a direct leg between two marks crosses land, add intermediate control points to the first element to shape the route:

```json
{
  "id": "anacapa",
  "type": "mark",
  "lat": 34.0133,
  "lon": -119.3631,
  "name": "Anacapa East End",
  "rounding": "port",
  "controlPointsToNext": [
    { "lat": 34.0250, "lon": -119.4100 }
  ]
}
```

`controlPointsToNext` coordinates are inserted between the current element and the next one in the rendered geometry. They are invisible as marks; they only shape the connecting line.

---

## Step 4: Author boats.json

All participating boats go in one file. Include at least one boat with `"isSelf": true` for Topaz.

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
        { "time": "2026-03-15T18:00:00Z", "lat": 34.4086, "lon": -119.6932 },
        { "time": "2026-03-15T21:30:00Z", "lat": 34.0850, "lon": -119.4200 },
        { "time": "2026-03-16T02:15:00Z", "lat": 34.4086, "lon": -119.6932 }
      ]
    },
    {
      "id": "rival",
      "name": "Rival",
      "color": "#ff8a5b",
      "boatType": "Olson 30",
      "source": "hand-authored",
      "isSelf": false,
      "track": [
        { "time": "2026-03-15T18:00:00Z", "lat": 34.4090, "lon": -119.6925 },
        { "time": "2026-03-15T21:40:00Z", "lat": 34.0862, "lon": -119.4188 },
        { "time": "2026-03-16T02:30:00Z", "lat": 34.4090, "lon": -119.6925 }
      ]
    }
  ]
}
```

**Required boat fields**: `id`, `name`, `color`, `boatType`, `source`, `isSelf`, `track`.

**`color`**: any CSS hex color. Topaz conventionally uses `#4fd1ff` (cyan). Pick contrasting colors for competitors.

**`source`**: use `"hand-authored"` for manually written tracks. Future import pipelines will use other values here.

**`isSelf`**: set `true` on Topaz, `false` for all competitors. The visualization shows the isSelf track before replay begins.

**`track`**: ordered time series. Each point requires `time` (RFC 3339, UTC), `lat`, and `lon`. Points do not need to be extremely dense — the replay engine interpolates smoothly between them.

---

## Step 5: Author events.json

Events are optional but recommended for interesting races. They appear as amber annotation markers on the map.

```json
{
  "events": [
    {
      "id": "anacapa-gybe",
      "type": "gybe",
      "time": "2026-03-15T21:35:00Z",
      "lat": 34.0250,
      "lon": -119.4100,
      "label": "Gybe at Anacapa",
      "description": "Topaz gybes around the east end of Anacapa and heads back north."
    }
  ]
}
```

**Required fields**: `id`, `type`, `time`.

**Optional fields**: `lat`, `lon` (explicit map position), `label` (short map callout), `description` (popup body text).

If `lat`/`lon` are omitted, the visualization interpolates the position from the isSelf boat track at the given `time`.

**Common `type` values**: `gybe`, `tack`, `wipeout`, `sail_change`, `note`.

---

## Step 6: Verify Locally

Start the Hugo development server and the Martin tile server (see `docs/race-preview-deployment.md` for full tile server instructions):

```sh
hugo server
```

Then open `http://localhost:1313/races/<season-slug>/<race-slug>/` to confirm the race page renders, the map loads, and the replay controls work.

---

## Reference Race

`content/races/dan-byrne-2025/bishop-rock-race/` is the canonical reference implementation. Copy its structure as a starting point for any new race.

## Schema Reference

The machine-readable contracts for all three JSON files are in `schemas/`:
- `schemas/race-course-v1.schema.json`
- `schemas/race-boats-v1.schema.json`
- `schemas/race-events-v1.schema.json`
