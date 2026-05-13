# GAPS.md

Gaps between the spec/tasks and the current implementation, plus bugs discovered during this audit. Organized by severity.

---

## Bugs

### BUG-1: Replay speed state initializes to 1x despite UI showing 60x

**File:** `assets/js/race-viz.js` — `createRaceVizState` (line 131), `syncReplayControls` (line 346)

`createRaceVizState` initializes `replay.speed` to `1`. The HTML `<select>` in the shortcode has `<option value="60" selected>` as the default. At boot, `syncReplayControls` sets `speedSelect.value = "1"`, but since there is no option with value `"1"`, the browser ignores the assignment and the select visually stays on 60x. However, `state.replay.speed` remains `1`, so the first time the user clicks Play without changing the speed selector, the replay runs at 1x (real-time). A 26-hour race would take 26 hours to replay.

The `data-race-viz-replay-speed="60"` attribute emitted by the shortcode is never read by `createRaceVizConfig`.

**Fix:** Initialize `replay.speed` to `60` in `createRaceVizState`, or read the initial speed from the shortcode attribute in `createRaceVizConfig` and consume it in state initialization.

---

### BUG-2: `waypoint` element type used in course data but not allowed by schema

**Files:** `content/races/dan-byrne-2025/meridian-400-race/course.json`, `content/races/dan-byrne-2025/catalina-backside-race/course.json`, `schemas/race-course-v1.schema.json`

The Meridian 400 and Catalina Backside courses use `"type": "waypoint"` on intermediate routing elements. The schema (`race-course-v1.schema.json`) only permits `"mark"`, `"start_line"`, and `"finish_line"` in the `type` enum. This means both files fail JSON schema validation.

In the renderer, `waypoint` type elements are silently excluded from the marks layer filter (`["==", ["get", "type"], "mark"]`) and the start/finish layer filter. They do appear in the route line coordinates (since `buildCourseRouteCoordinates` iterates all elements), which may be intentional — but the behavior is undocumented and unsupported by the schema.

The proper V1 mechanism for invisible route-shaping points is `controlPointsToNext` on the preceding element. The `waypoint` type is being used as a workaround.

**Options:**
- Add `"waypoint"` to the schema's type enum and document its rendering semantics (route-visible, no marker), OR
- Replace `waypoint` elements in the existing course files with `controlPointsToNext` arrays on the preceding element.

---

### BUG-3: Empty `name` string on `south-of-channel-islands` element fails schema validation

**File:** `content/races/dan-byrne-2025/meridian-400-race/course.json`

The `south-of-channel-islands` element has `"name": ""`. The schema defines `name.minLength: 1` (if the field is present, it must contain at least one character). An empty string violates this constraint.

**Fix:** Either omit the `name` field entirely (it is optional in the schema), or remove the `minLength: 1` constraint from the `name` property to allow intentionally unnamed elements, or lower it to `"minLength": 0"`.

---

### BUG-4: `catalina-south-coast` waypoint with a name renders an orphaned text label

**File:** `content/races/dan-byrne-2025/catalina-backside-race/course.json`

The `catalina-south-coast` element has `"type": "waypoint"` (already flagged by BUG-2 as schema-invalid) and `"name": "Catalina South Coast"`. Because `renderCourseLayers` builds the labels layer with only a name-check filter (`["!=", ["get", "name"], ""]`) and no type guard, this element renders a floating "CATALINA SOUTH COAST" text label on the map with no corresponding circle beneath it. The result looks broken — a disembodied label over open water.

The `south-of-channel-islands` waypoint in the Meridian 400 had an empty name so it was silently invisible; the Catalina waypoint is actively misleading.

**Fix:** Either replace the `waypoint` element with `controlPointsToNext` on the preceding mark (the V1-correct approach per BUG-2), or add a type guard to `labelsLayerID` so that only `mark`, `start_line`, and `finish_line` elements emit labels.

---

### BUG-5: `jibeset-import` compiled binary at repo root is not gitignored

**File:** `.gitignore`

Running `go build -o jibeset-import ./cmd/jibeset-import` (or similar) produces a `jibeset-import` Mach-O binary at the repo root. The `.gitignore` includes `gpx-import` (the equivalent for the GPX tool) but not `jibeset-import`. The binary is currently untracked (`?? jibeset-import` in `git status`) but will be accidentally staged if anyone runs `git add .`.

**Fix:** Add `jibeset-import` to `.gitignore` alongside `gpx-import`.

---

## Spec Gaps

### GAP-1: Track hover tooltip shows only boat name — timestamp missing

**Spec requirement:** "hover tooltip with timestamp information"

**File:** `assets/js/race-viz.js` — `attachTrackHoverInteractions` (line 1716)

The boat marker hover (`attachBoatMarkerHoverInteractions`) shows boat name plus the current replay clock time — this satisfies the spec for live replay markers. However, the track line hover (`attachTrackHoverInteractions`) builds a tooltip that shows only the boat name, with no timestamp. When a user hovers over a static track line or a replay tail, there is no indication of when the boat was at the hovered point.

To satisfy the spec, the track hover should interpolate the closest track time for the cursor's position on the line and display it in the tooltip.

---

### GAP-2: Mark rounding direction is stored but never visually rendered

**Spec requirement:** "rounding direction support" with values `port`, `starboard`, `none`

**Files:** All `course.json` files, `assets/js/race-viz.js` — `renderCourseLayers` (line 1085)

The `rounding` field is required in the schema and is stored in GeoJSON feature properties (`buildCourseFeatures`), but `renderCourseLayers` never uses it. There is no visual indicator on the mark showing which way to round. A sailor reading the map cannot tell from the visualization alone whether to pass a mark on the port or starboard side.

A typical V1 approach: render a small arc or arrow on the mark, or apply a distinct fill color to indicate rounding direction.

---

### GAP-3: ~~Four race pages are course-only~~ — resolved; one outstanding data quality issue remains

**Status:** Largely resolved as of 2026-05-13. All four races (`catalina-backside-race`, `meridian-400-race`, `ship-rock-race`, `sb-island-race`) now have `boats.json` and use `mode="replay"`. The Jibeset import tool (`jibesetimport` package + `cmd/jibeset-import` CLI) was used to populate track data.

**Remaining issue:** None of the four newly populated races have an `events.json`. Replay works end-to-end, but the event annotation layer is untested beyond Bishop Rock. At minimum one of the new race pages should gain an `events.json` to exercise the full path.

---

### GAP-4: Boats load failure has no user-visible error state

**File:** `assets/js/race-viz.js` — `setBoatsState` (line 305), `loadBoats` (line 1814)

When `loadBoats` fails (network error, malformed JSON, etc.), `setBoatsState(root, stage, state, "error")` sets `data-race-viz-boats-state="error"` on the root and stage elements, but there is no corresponding visual fallback. The map renders the chart and course normally, the fleet panel is blank, the replay controls are disabled, and the user has no indication of what went wrong.

By contrast, course failures call `renderCourseFallback(stage, message)` which displays a visible error banner at the bottom of the map stage.

**Fix:** Add a `renderBoatsFallback(stage, message)` function (or reuse `renderCourseFallback`) so that a boats-load failure shows an actionable error in the UI.

---

### GAP-5: No shortcode parameter to set map `minZoom`

**File:** `layouts/shortcodes/race-viz.html`, `assets/js/race-viz.js` — `initializeMap` (line 591)

The MapLibre `Map` constructor is called without a `minZoom` option. The effective floor on outward zoom is controlled indirectly by `mapMaxBounds` — if the bounds are tight, MapLibre refuses to zoom out past the level where the bounds would leave the viewport. This is a blunt instrument: adjusting bounds to allow more zoom-out also changes the extent users can pan, which are two independent concerns.

A `mapMinZoom` shortcode parameter (mirroring the existing `fitMaxZoom` and `mapMaxBounds`) would let individual race pages set an explicit minimum zoom level that fits their course geography without needing to widen bounds as a workaround.

---

## UX Issues

### UX-1: Timeline labels show time-of-day (HH:MM:SS UTC), not elapsed race time

**File:** `assets/js/race-viz.js` — `formatReplayClockLabel` (line 338)

The replay clock and timeline scale labels format timestamps as UTC time-of-day (`"08:00:00"` to `"10:28:00"`). For the Bishop Rock race these happen to read intuitively because the start and end are on the same day's AM hours. But for any race that spans midnight — including the 26-hour Bishop Rock race itself (which runs from `2025-02-11T08:00:00Z` to `2025-02-12T10:28:00Z`) — the start label correctly shows `08:00:00` and the end label correctly shows `10:28:00`, but a viewer cannot tell from the labels alone that these are on different calendar days. The elapsed replay time is also not shown.

Suggested alternatives:
- Display elapsed time from race start (e.g., `+00:00:00` to `+26:28:00`), or
- Show a date prefix when the race crosses midnight (e.g., `Feb 12 10:28`).

---

### UX-2: Start and finish rendered as single point markers, not as crossing lines

**Files:** `assets/js/race-viz.js` — `renderCourseLayers`, `schemas/race-course-v1.schema.json`

In the current data model, `start_line` and `finish_line` are single lat/lon points rendered as colored circles (yellow for start, pink for finish). In actual sailing, the start and finish are defined by two marks — a line between a committee boat and a pin — that competitors must cross. The current single-point representation does not convey the geometry of these lines and may be unclear to viewers who are not already familiar with the course.

The V1 data model was intentionally limited to a single lat/lon per element (see SPEC.md §Race Course Model). This gap is by design in V1, but it should be tracked as a V2 data model extension: allow `start_line` and `finish_line` to specify a pair of coordinates defining the actual line segment.

---

### UX-3: Track hover tooltip is static — does not follow the cursor along a track line

**File:** `assets/js/race-viz.js` — `attachTrackHoverInteractions` (line 1716)

The track hover interaction binds to `mouseenter` on the tracks and replay-tails layers. The tooltip is placed at `event.lngLat` at the moment the cursor first enters the feature. As the user moves the cursor along the track line, the tooltip stays fixed at the entry point — it does not update position or content.

For point features (boat markers) this is fine. For line features (tracks and replay tails), a static tooltip at the entry edge is both positionally wrong and makes the "timestamp at position" enhancement (GAP-1) harder to wire up later, since a `mousemove` handler would be needed anyway to resolve the correct track segment.

**Fix:** Replace the `mouseenter` handler with a `mousemove` listener that updates `popup.setLngLat(event.lngLat)` on each movement, and add the timestamp interpolation from GAP-1 at the same time. These two fixes are naturally batched.

---

### UX-4: Event annotations are always visible during replay regardless of current replay time

**File:** `assets/js/race-viz.js` — `buildEventFeatures`, `loadEvents`

Event features are built once when `loadEvents` runs and remain permanently rendered on the map for the full duration of replay. An event that occurred 22 hours into a 26-hour race shows its amber marker and label from the very first second of playback.

This is defensible for V1 (always show all context) but creates an inconsistency: boat positions are time-accurate, but event annotations show the race's entire history at every moment. A viewer watching a close start mark rounding will see a label near the finish 400nm away as if it already happened.

For V2: filter event visibility in `renderReplayFrame` so annotations appear only after `state.replay.currentTimeMs` passes the event's own timestamp.

---

## Minor / Consistency Issues

### MINOR-1: `fitMaxZoom` bypasses the config object

**File:** `assets/js/race-viz.js` — `fitCourseBounds` (line 1265)

`fitCourseBounds` reads `fitMaxZoom` directly from `root.dataset.raceVizFitMaxZoom` instead of routing it through `createRaceVizConfig`. All other configuration goes through the config object for consistent access. This is inconsistent and makes the config object an incomplete representation of the component's configuration.

---

### MINOR-2: `data-race-viz-replay-speed` attribute emitted but never consumed

**Files:** `layouts/shortcodes/race-viz.html` (line 54), `assets/js/race-viz.js` — `createRaceVizConfig`

The shortcode emits `data-race-viz-replay-speed="60"` on the root element, but `createRaceVizConfig` does not read it. The attribute is therefore unused. This is also the root cause of BUG-1.

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

### DESIGN-1: Side-by-side layout — map and sidebar horizontally adjacent on large viewports

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

### DESIGN-2: Replace glow aesthetic with flat/crisp container style

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

### DESIGN-4: Small visual polish items

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
