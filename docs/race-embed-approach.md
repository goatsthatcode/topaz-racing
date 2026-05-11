# Race Visualization Embed Approach

## Decision
Use a Hugo shortcode as the canonical V1 embed mechanism:

```text
{{< race-viz >}}
```

for a race page embedding its own visualization, and:

```text
{{< race-viz race="dan-byrne-2025/bishop-rock-race" >}}
```

for embedding the same visualization from other prose pages.

## Race Identifier Convention
The race ID is the Hugo content path under `content/races/`, excluding `index.md`.

Examples:
- `dan-byrne-2025/bishop-rock-race`
- `2026/transpac-start`

This identifier is stable across markdown embeds, template lookups, and future frontend bootstrapping because it matches the canonical race bundle location.

## Shortcode Behavior
- `race-viz` resolves the target race bundle either from the current page or from the explicit `race` parameter.
- `course.json` is always required because it is the shared baseline payload for every mode.
- `boats.json` is required only for `replay` mode.
- `events.json` is optional and is published into the embed contract only when present.
- The rendered HTML emits `data-*` attributes for the race ID, selected mode, and published JSON asset URLs.
- The map/replay frontend can later hydrate against this stable DOM contract without changing page authoring.

## Why Shortcodes
- Keeps raw JSON out of markdown prose.
- Works naturally in race pages and ordinary posts.
- Avoids editing vendor theme templates for each visualization instance.
- Preserves a mostly static Hugo publishing model because embeds resolve to published page resources.

## Authoring Notes
- Canonical race pages should usually use `{{< race-viz >}}` so the prose bundle remains self-contained.
- Cross-page references should use the explicit `race` parameter.
- `mode` is reserved as a shared configuration knob for `course` and `replay` views.
- A `course` embed can render from a race bundle that only includes `course.json`; add `boats.json` when the same race is ready for replay.
