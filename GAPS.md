# GAPS.md

Gaps between the spec/tasks and the current implementation, plus bugs discovered during this audit. Organized by severity.

---

## Bugs

### ~~BUG-1: Replay speed state initializes to 1x despite UI showing 60x~~

- [x] Resolved on 2026-05-13 (commit c522531): Added `replaySpeed` to `createRaceVizConfig` reading `root.dataset.raceVizReplaySpeed` (defaulting to 60). `createRaceVizState` now seeds `replay.speed` from `config.replaySpeed`. The HTML attribute `data-race-viz-replay-speed="60"` is now consumed by the config object, so the replay clock starts at 60x on first play.

---

### ~~BUG-2: `waypoint` element type used in course data but not allowed by schema~~

- [x] Resolved on 2026-05-13: Added `"waypoint"` to the `type` enum in `schemas/race-course-v1.schema.json` with a description documenting its semantics (route-visible, no marker, no label). The Catalina Backside course's `catalina-south-coast` waypoint now passes schema validation. Tests added in `race_course_rendering_test.go`.

---

### BUG-3: Empty `name` string on `south-of-channel-islands` element fails schema validation

**File:** `content/races/dan-byrne-2025/meridian-400-race/course.json`

The `south-of-channel-islands` element has `"name": ""`. The schema defines `name.minLength: 1` (if the field is present, it must contain at least one character). An empty string violates this constraint.

**Fix:** Either omit the `name` field entirely (it is optional in the schema), or remove the `minLength: 1` constraint from the `name` property to allow intentionally unnamed elements, or lower it to `"minLength": 0"`.

---

### ~~BUG-4: `catalina-south-coast` waypoint with a name renders an orphaned text label~~

- [x] Resolved on 2026-05-13: Added a type guard to the `labelsLayerID` layer filter in `renderCourseLayers` so that only `mark`, `start_line`, and `finish_line` elements emit labels: `["match", ["get", "type"], ["mark", "start_line", "finish_line"], true, false]`. Waypoints with names no longer produce disembodied labels on the map. Test added in `race_course_rendering_test.go`.

---

### ~~BUG-5: `jibeset-import` compiled binary at repo root is not gitignored~~

- [x] Resolved on 2026-05-13 (commit c522531): Added `/jibeset-import` to `.gitignore` alongside `/gpx-import`.

---

## Spec Gaps

### ~~GAP-1: Track hover tooltip shows only boat name — timestamp missing~~

- [x] Resolved on 2026-05-13: Added `interpolateTimeFromPosition(boat, lngLat)` which projects the cursor onto the nearest track segment via dot-product and interpolates the timestamp linearly within that segment. `attachTrackHoverInteractions` now looks up the boat from `state.replay.timeline` by feature ID, calls `interpolateTimeFromPosition`, and displays the result as a `<time class="race-viz-hover-time">` element. Fixed together with UX-3.

---

### ~~GAP-2: Mark rounding direction is stored but never visually rendered~~

- [x] Resolved on 2026-05-13: Added `roundingLayerID` (`race-viz-course-marks-rounding`) in `renderCourseLayers`. The layer renders a larger outer circle behind each `mark` element using `roundingPortColor` (red) for port-rounding and `roundingStarboardColor` (green) for starboard-rounding marks. Waypoints and start/finish elements are excluded by the `["==", ["get", "type"], "mark"]` filter. Tests in `race_course_rendering_test.go` (`TestRaceVizBootstrapImplementsMarkRoundingDirectionLayer`).

---

### GAP-3: ~~Four race pages are course-only~~ — resolved; one outstanding data quality issue remains

**Status:** Largely resolved as of 2026-05-13. All four races (`catalina-backside-race`, `meridian-400-race`, `ship-rock-race`, `sb-island-race`) now have `boats.json` and use `mode="replay"`. The Jibeset import tool (`jibesetimport` package + `cmd/jibeset-import` CLI) was used to populate track data.

**Remaining issue:** None of the four newly populated races have an `events.json`. Replay works end-to-end, but the event annotation layer is untested beyond Bishop Rock. At minimum one of the new race pages should gain an `events.json` to exercise the full path.

---

### ~~GAP-4: Boats load failure has no user-visible error state~~

- [x] Resolved on 2026-05-13 (commit c522531): Added `renderBoatsFallback(stage, message)` in `assets/js/race-viz.js`, called from the `loadBoats` error path. A boats-load failure now shows an error banner at the bottom of the map stage, consistent with `renderCourseFallback`.

---

### ~~GAP-5: No shortcode parameter to set map `minZoom`~~

- [x] Resolved on 2026-05-13: Added `mapMinZoom` shortcode parameter. The shortcode emits `data-race-viz-map-min-zoom="{{ . }}"` when the parameter is set. `createRaceVizConfig` reads it into `config.map.minZoom` and `createMapInstance` conditionally passes it as `minZoom` to the MapLibre `Map` constructor. Tests added in `tests/race_map_zoom_config_test.go`.

---

## UX Issues

### ~~UX-1: Timeline labels show time-of-day (HH:MM:SS UTC), not elapsed race time~~

- [x] Resolved on 2026-05-13: Replaced `formatReplayClockLabel` with `formatElapsedLabel(elapsedMs)` which formats elapsed milliseconds from race start as `+HH:MM:SS`. `syncReplayControls` now feeds `currentTimeMs - startTimeMs` for the current label, `durationMs` for the end label, and `0` for the start label. Both hover tooltips (boat marker and track hover) also use elapsed time. Tests added in `tests/race_elapsed_label_test.go`.

---

### UX-2: Start and finish rendered as single point markers, not as crossing lines

**Files:** `assets/js/race-viz.js` — `renderCourseLayers`, `schemas/race-course-v1.schema.json`

In the current data model, `start_line` and `finish_line` are single lat/lon points rendered as colored circles (yellow for start, pink for finish). In actual sailing, the start and finish are defined by two marks — a line between a committee boat and a pin — that competitors must cross. The current single-point representation does not convey the geometry of these lines and may be unclear to viewers who are not already familiar with the course.

The V1 data model was intentionally limited to a single lat/lon per element (see SPEC.md §Race Course Model). This gap is by design in V1, but it should be tracked as a V2 data model extension: allow `start_line` and `finish_line` to specify a pair of coordinates defining the actual line segment.

---

### ~~UX-3: Track hover tooltip is static — does not follow the cursor along a track line~~

- [x] Resolved on 2026-05-13: Replaced the `mouseenter` handler in `attachTrackHoverInteractions` with a `mousemove` handler. The tooltip now follows the cursor: if the popup already exists it is updated via `setLngLat(lngLat).setHTML(html)` without recreation; if not, it is created fresh. Fixed together with GAP-1.

---

### UX-4: Event annotations are always visible during replay regardless of current replay time

**File:** `assets/js/race-viz.js` — `buildEventFeatures`, `loadEvents`

Event features are built once when `loadEvents` runs and remain permanently rendered on the map for the full duration of replay. An event that occurred 22 hours into a 26-hour race shows its amber marker and label from the very first second of playback.

This is defensible for V1 (always show all context) but creates an inconsistency: boat positions are time-accurate, but event annotations show the race's entire history at every moment. A viewer watching a close start mark rounding will see a label near the finish 400nm away as if it already happened.

For V2: filter event visibility in `renderReplayFrame` so annotations appear only after `state.replay.currentTimeMs` passes the event's own timestamp.

---

## Minor / Consistency Issues

### ~~MINOR-1: `fitMaxZoom` bypasses the config object~~

- [x] Resolved on 2026-05-13: Added `fitMaxZoom` to `createRaceVizConfig` reading `root.dataset.raceVizFitMaxZoom`. Updated `fitCourseBounds` to accept `config` instead of `root`, reading `config.fitMaxZoom`. The call site now passes `state.config`. Tests added in `tests/race_map_zoom_config_test.go`.

---

### ~~MINOR-2: `data-race-viz-replay-speed` attribute emitted but never consumed~~

- [x] Resolved on 2026-05-13 (commit c522531): Fixed together with BUG-1. `createRaceVizConfig` now reads `root.dataset.raceVizReplaySpeed` into `config.replaySpeed`, which is used to initialize `state.replay.speed`.

---

### MINOR-3: Generic figcaption text provides no race context

**File:** `layouts/shortcodes/race-viz.html` (lines 141–147)

The figcaption renders "Race visualization embed for `dan-byrne-2025/bishop-rock-race`." for all race pages. The spec notes that "race overlays should feel intentional and editorial." A figcaption that describes the race (e.g., the race name, date, and distance) would be more useful and consistent with the intended editorial tone.

---

### MINOR-4: Reference dataset has only one event of one type

**File:** `content/races/dan-byrne-2025/bishop-rock-race/events.json`

The reference `events.json` has a single `gybe` event. The schema supports `gybe`, `wipeout`, `sail_change`, `tack`, `mark_rounding`, and `note`. A richer reference dataset would better serve as a development test fixture and as an example for authoring new races. Task 1.4 required "at least one event annotation," which is satisfied, but the reference data does not exercise the full breadth of the annotation system.

---

### MINOR-5: Jibeset import tool implemented but not recorded in TASKS.md

**Files:** `cmd/jibeset-import/main.go`, `jibesetimport/parser.go`, `jibesetimport/parser_test.go`

The Jibeset multi-boat track import tool is fully implemented and used to populate `boats.json` for the new race pages. It is more capable than the GPX tool (multi-boat, sail-number filtering, race-date selection, interval downsampling, per-boat color/name overrides). However, TASKS.md only documents Task 9.1 (GPX import). The Jibeset tool appears in the deferred backlog as "import/conversion tooling for Jibeset and Garmin Connect export formats" but is already done. It should be added as a completed task (Task 9.2 or similar) so the backlog stays accurate and a future contributor knows the tool exists.

---

## Design Improvements

These are not bugs or spec gaps — they are design and layout improvements identified during a style audit. Each is self-contained and can be implemented independently.

---

### ~~DESIGN-1: Side-by-side layout — map and sidebar horizontally adjacent on large viewports~~

- [x] Resolved on 2026-05-13: Added `@media (min-width: 56rem)` block to `assets/css/race-viz.css`. `.race-viz-shell` becomes a two-column CSS grid (`1fr 19rem`) with `grid-template-areas` routing the stage, sidebar, and figcaption. The sidebar gets `margin-top: 0`, `border-left: 1px solid var(--site-border-subtle)`, and `overflow-y: auto`. The race-page stage `min-height` drops from `78svh` to `65svh` in this context since the sidebar no longer adds vertical footprint. Tests added in `tests/race_visual_design_test.go`.

**Files:** `assets/css/race-viz.css`, `layouts/races/single.html`

**Current behavior:** The sidebar (replay controls + fleet legend) always sits below the map stage in a full-width vertical stack. On a wide monitor, the map takes a narrow vertical rectangle and is followed by a wide-but-shallow strip of controls.

**Proposed:** At a breakpoint of approximately `56rem` (≈900px), convert `.race-viz-shell` to a two-column CSS grid. The map stage occupies the left column (`1fr`) and the sidebar occupies a fixed right column (`18rem–20rem`). Below the breakpoint, the layout reverts to map on top, sidebar below.

**Implementation sketch:**

```css
/* ≥ 56rem: side-by-side */
@media (min-width: 56rem) {
  .race-viz-shell {
    display: grid;
    grid-template-columns: 1fr 19rem;
    grid-template-rows: 1fr auto;
    grid-template-areas:
      "stage  sidebar"
      "caption caption";
    gap: 0;
  }

  .race-viz-stage  { grid-area: stage; }
  .race-viz-sidebar {
    grid-area: sidebar;
    margin-top: 0;
    border-left: 1px solid var(--site-border-subtle);
    overflow-y: auto;
    padding: 0.85rem 0.9rem;
  }

  .race-viz-shell figcaption { grid-area: caption; }
}
```

On a race page (where `.race-viz-shell` is already edge-to-edge), the 19rem sidebar column uses `align-self: stretch` so it fills the full map height. The map `min-height` can drop from `78svh` to `65svh` in side-by-side context since the sidebar no longer adds vertical footprint.

---

### ~~DESIGN-2: Replace glow aesthetic with flat/crisp container style~~

- [x] Resolved on 2026-05-13: Stripped all glow/gradient treatments from containers. `.race-viz-shell` background replaced with a single flat fill (`rgba(6, 14, 25, 0.97)`); its `box-shadow` drop-shadow removed. `.race-viz-panel` changed to `background: transparent` with a top-only hairline border (`border-top: 1px solid rgba(126, 245, 236, 0.09)`); card `border-radius` removed. `.race-viz-button` replaced gradient background with flat `rgba(8, 20, 34, 0.95)`; inset glow ring removed. `figcaption::before` glow `box-shadow` removed. `body > header` (custom.css) `box-shadow` removed. `.race-page-header` gradient (custom.css) removed. Map-layer, event-popup, and hover-tooltip styles left untouched. Tests added in `tests/race_visual_design_test.go`.

**Files:** `assets/css/race-viz.css`, `assets/css/custom.css`

**Current behavior:** The shell, panels, buttons, and figcaption icon all use layered radial gradients and `box-shadow` glows — atmosphere-building at first glance but visually noisy alongside the already-rich nautical chart. The glow competes with the map canvas, which is the intended centerpiece.

**Proposed direction:** Strip all containers to flat, dark fills with hairline borders. The cyan accent (`#7ef5ec`) remains as text color and single-pixel border color — no bloom, no radial spread, no `box-shadow` glow.

**Specific changes (leave map-layer, event-popup, and hover-tooltip styles untouched):**

- **`.race-viz-shell` background:** Replace the three-layer radial+linear gradient stack with a single flat fill: `background: rgba(6, 14, 25, 0.97)`. Remove `box-shadow: 0 1.25rem 3rem …`.
- **`.race-viz-panel` background:** Remove gradient and radial glow. Use `background: transparent` — the sidebar column's own background provides the fill. Change the `border` to a top-only hairline: `border: none; border-top: 1px solid rgba(126, 245, 236, 0.09)`. Remove `border-radius`.
- **`.race-viz-button` background:** Replace `linear-gradient(180deg, …)` with `background: rgba(8, 20, 34, 0.95)`. Remove `box-shadow: inset 0 0 0 1px …`.
- **`figcaption::before` (the caption dot):** Remove `box-shadow: 0 0 0.75rem rgba(126, 245, 236, 0.4)`. The two-layered `radial-gradient` fill is enough visual identity without the glow ring.
- **`body > header` (custom.css):** Remove `box-shadow: 0 0.5rem 2rem var(--site-cyan-dim)`. The `border-bottom` alone is the right weight for a site header.
- **`.race-page-header` gradient (custom.css):** Remove the subtle `linear-gradient(180deg, transparent, rgba(126,245,236,0.04))`. Flat transparent reads more editorial and lets the page background do the work.

---

### DESIGN-3: Sidebar section structure — sections over floating cards

**Related to DESIGN-1 and DESIGN-2.** In a 19rem sidebar column, the current panel "card" metaphor (gradient background, border-radius, full border) produces cramped floating boxes that feel generic. A narrow column reads better as labeled sections with hairline dividers — a chart instrument panel, not a dashboard widget stack.

**Proposed:** When in side-by-side mode, panels become flat sections:
- No `border-radius`, no card background — the sidebar column background provides the fill
- Section separation via `border-top: 1px solid rgba(126, 245, 236, 0.09)` (first section has no top border)
- `padding: 0.65rem 0` (remove horizontal padding — the sidebar column's own padding handles it)
- The `REPLAY` / `FLEET` sidebar titles drop to `font-size: 0.74rem` and `letter-spacing: 0.1em` (current 0.16em is wide for a narrow column)
- The replay clock readout (`race-viz-replay-time`) becomes the most prominent element in the sidebar: `font-size: 1rem; letter-spacing: 0.1em` — a single number in a tight column reads like a navigation instrument

---

### ~~DESIGN-4: Small visual polish items~~

- [x] Resolved on 2026-05-13: (1) Boat swatch changed from a `0.85rem` circle to a `1.2rem × 3px` horizontal tick (`border-radius: 2px`) to match the track-line appearance on the map. (2) Speed select `border-radius` reduced from `0.75rem` to `0.35rem` for a crisper look consistent with the flat direction. (3) `.race-viz-sidebar-title` `letter-spacing` reduced from `0.16em` to `0.1em` for better legibility in narrow columns. (4) `.race-viz-shell` `border-radius` reduced from `1rem` to `0.5rem` for a sharper embedded component appearance. Tests added in `tests/race_visual_design_test.go`.

**File:** `assets/css/race-viz.css`

Four targeted refinements that work regardless of whether DESIGN-1/2/3 land first:

1. **Boat color swatch — tick mark instead of dot.** The current `0.85rem` circle swatch (`border-radius: 999px`) does not visually match how boats appear on the map (as colored track lines). Replace with a horizontal tick: `width: 1.2rem; height: 3px; border-radius: 2px`. This reads as a track-line sample rather than a generic color dot and is more legible at small sizes.

2. **Speed select border-radius.** The speed selector has `border-radius: 0.75rem`, which is softer than the `999px` pill buttons beside it. Unify to `border-radius: 0.35rem` for a crisper select element that matches the flat direction of DESIGN-2 without fully squaring off.

3. **Section heading letter-spacing.** `.race-viz-sidebar-title` uses `letter-spacing: 0.16em`. At this tracking, `FLEET` in a narrow column reads as wide-spaced caps rather than a label. Reduce to `letter-spacing: 0.1em` for legibility.

4. **Shell `border-radius` on embedded (non-race-page) instances.** The current `1rem` border-radius on `.race-viz-shell` is generous for a component embedded in prose. `0.5rem` is sharper and fits the flatter direction; it also aligns with the `0.5rem` already used on the stage in the race-page context.

---

## Out-of-Scope / V2 Notes

These are not gaps against the V1 spec but are worth noting as they came up during the audit:

- **Land-aware routing:** The existing course data routes all legs as straight lines. The Bishop Rock course line crosses Santa Catalina Island if drawn directly. The spec calls for manual `controlPointsToNext` as the V1 fallback, which is supported in the renderer and schema, but none of the course files uses it yet.
- **Start/finish as two-point line:** Tracked as UX-2 above — a V2 data model extension.
- **No absolute date in time display:** Tracked as UX-1 — the current HH:MM:SS format is ambiguous for multi-day races.
- **Event annotation time-gating during replay:** Tracked as UX-4 above — a V2 behavior enhancement.
- **Tile coverage south of Tijuana:** The Meridian 400 southern turning mark (31.64°N) and other extreme waypoints fall outside the current ENC tile coverage. The map renders the course correctly but shows blank tiles for the southernmost area. Extending tile coverage to cover the full offshore Southern California racing region (down to ~30°N) is a future map tile generation task.
