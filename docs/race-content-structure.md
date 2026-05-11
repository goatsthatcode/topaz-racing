# Race Content Structure

## Decision
Race pages live as Hugo leaf bundles under `content/races/`. Each race bundle owns:

- `index.md` for prose and page metadata
- `course.json` for the ordered race course
- `boats.json` for one or more boat tracks
- `events.json` for race annotations

This keeps prose separate from visualization data while letting Hugo publish a canonical race page from a single directory.

## Canonical Directory Shape

```text
content/
  races/
    <series-or-season>/
      _index.md
      <race-slug>/
        index.md
        course.json
        boats.json
        events.json
```

## Notes
- The race bundle directory name is the canonical race slug used in URLs.
- JSON assets stay next to the race page instead of being embedded into markdown.
- Additional supporting assets for a race can be added to the same bundle later without changing the page model.

## Reference Race
`content/races/dan-byrne-2025/bishop-rock-race/` is the reference implementation for new races.
