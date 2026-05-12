# TASKS.md

## Purpose
This document turns [SPEC.md](/Users/lorenbrindze/Projects/topaz-racing/SPEC.md:1) into an implementation backlog for the Topaz Racing race visualization system.

The intended execution style is:
- work top-down by priority unless blocked
- keep the site deployable as a mostly static Hugo project
- favor small end-to-end slices over broad unfinished scaffolding
- preserve future support for live tracking, import tooling, and richer storytelling overlays

Repository hygiene for task completion:
- treat [SPEC.md](/Users/lorenbrindze/Projects/topaz-racing/SPEC.md:1) plus the source tree as the source of truth
- when a task is marked done, ensure its authored files in `content/`, `layouts/`, `assets/`, `schemas/`, `docs/`, and tests are staged or committed
- do not treat generated `public/` output as a required deliverable unless a task explicitly says to publish build artifacts

Test organization guidance:
- keep pure data-contract tests together
- keep Hugo-render integration tests together
- prefer shared helpers over repeating `hugo` build setup in each root-level `_test.go` file

## Milestone 0: Repository And Architecture Baseline

### Task 0.1
Decide and document the on-disk content structure for race pages and race data under `content/`.

- [x] Done on 2026-05-11: documented the canonical race bundle shape in `docs/race-content-structure.md` and added `content/races/dan-byrne-2025/bishop-rock-race/` as the reference pattern.

Deliverables:
- one canonical race directory shape
- clear separation between prose content and JSON assets
- one example race directory chosen as the reference pattern

Acceptance criteria:
- a new race can be added by copying a known folder structure
- the structure is compatible with Hugo page rendering and static asset publication

### Task 0.2
Choose the Hugo integration approach for embedding race visualizations.

- [x] Done on 2026-05-11: added the `race-viz` shortcode in `layouts/shortcodes/race-viz.html`, documented the embed contract and race ID convention in `docs/race-embed-approach.md`, and covered the rendered output with `race_embed_shortcode_test.go`.

Deliverables:
- one selected embed mechanism, preferably shortcode-based
- a documented race identifier convention

Acceptance criteria:
- a race page can reference a visualization without inlining raw JSON in markdown

### Task 0.3
Define the frontend asset strategy for the visualization code.

- [x] Done on 2026-05-11: kept race visualization assets in project-owned Hugo Pipes files under `assets/js/` and `assets/css/`, documented the conditional loading strategy in `docs/frontend-asset-strategy.md`, and covered race-page-only asset emission with `race_asset_strategy_test.go`.

Deliverables:
- chosen location for JS/CSS assets in the Hugo project
- plan for bundling/loading map code in race pages only when needed

Acceptance criteria:
- visualization assets are not tied to theme internals unnecessarily
- the site remains maintainable without editing vendor theme code for core app logic

### Task 0.4
Commit to a shared visualization engine and data model for V1 modes.

- [x] Done on 2026-05-11: documented the shared engine in `docs/race-visualization-architecture.md`, updated `layouts/shortcodes/race-viz.html` and `assets/js/race-viz.js` to use one engine contract plus stable layer slots for both modes, and covered the shared-mode contract with `race_visualization_architecture_test.go`.

Deliverables:
- one architecture note or implementation decision covering shared map, overlay, and state primitives
- explicit layering plan for `course` and `replay` modes on top of the same engine
- boundaries between shared visualization logic and mode-specific behavior

Acceptance criteria:
- `course` and `replay` do not diverge into separate unrelated component systems
- shared data contracts and rendering primitives are identified before substantial mode-specific UI work begins

## Milestone 1: Data Contracts

### Task 1.1
Define the V1 JSON schema for race course data.

- [x] Done on 2026-05-11: defined the machine-readable course contract in `schemas/race-course-v1.schema.json`, documented the authoring semantics in `docs/race-course-schema.md`, and added `race_course_schema_test.go` to keep the reference course aligned with the V1 schema.

Deliverables:
- course file format
- ordered course element list
- element types for `mark`, `start_line`, and `finish_line`
- required fields for each course element: `id`, `type`, `lat`, `lon`
- optional `name` field
- `rounding` semantics with values such as `port`, `starboard`, and `none`
- support for optional manual route-shaping control points

Acceptance criteria:
- a single race course can be described completely in JSON
- the schema is sufficient for a standalone `course` mode render without any track data present
- the schema leaves room for future exclusion zones and auto-routing

### Task 1.2
Define the V1 JSON schema for boat track data.

- [x] Done on 2026-05-11: defined the machine-readable boat-track contract in `schemas/race-boats-v1.schema.json`, documented the authoring model in `docs/race-boats-schema.md`, and added `race_boats_schema_test.go` to keep the reference `boats.json` aligned with the V1 replay payload.

Deliverables:
- per-boat metadata contract
- required boat metadata fields: `id`, `name`, `color`, `boatType`, `source`, `isSelf`
- per-point time/lat/lon contract
- expected file organization for one or more boats in a race

Acceptance criteria:
- one self boat and multiple competitor boats can be described cleanly
- the model is independent of Jibeset, Garmin, GPX, or other upstream formats

### Task 1.3
Define the V1 JSON schema for event annotations.

- [x] Done on 2026-05-11: defined the machine-readable event-annotation contract in `schemas/race-events-v1.schema.json`, documented time and position anchoring in `docs/race-events-schema.md`, and added `race_events_schema_test.go` to keep the reference `events.json` aligned with the V1 annotation payload.

Deliverables:
- annotation types
- time-anchored and/or position-anchored event model
- optional label/description fields

Acceptance criteria:
- a race can include notable moments like gybes, wipeouts, or sail changes

### Task 1.4
Create one complete sample race dataset.

- [x] Done on 2026-05-11: completed the Bishop Rock reference bundle with `course.json`, `boats.json`, and `events.json`, and added `race_sample_dataset_test.go` to verify it covers the full V1 course, boats, and annotation flow.

Deliverables:
- one real or representative course file
- one `isSelf` boat track
- at least one competitor track
- at least one event annotation

Acceptance criteria:
- the dataset is sufficient to develop and validate the full V1 visualization flow

## Milestone 2: Map Foundation

### Task 2.1
Audit and adapt the existing `tiles/` prototype into a reusable map foundation for the Hugo site.

- [x] Done on 2026-05-11: documented the prototype-to-site contract in `docs/race-map-foundation.md`, published a site-owned style resource through `layouts/partials/race-viz/map-foundation.html` and `assets/race-viz/map/style.json.tmpl`, surfaced the map contract on `race-viz` embeds, and added `race_map_foundation_test.go`.

Deliverables:
- documented relationship between `tiles/index.html`, `tiles/style.json`, and site integration
- decision on how the site loads the map style in local preview and published environments

Acceptance criteria:
- the map can render inside the site without depending on the standalone prototype page

### Task 2.2
Create a reusable map component that renders the ENC-style vector chart background in a Hugo page.

- [x] Done on 2026-05-11: added a concrete map canvas to `layouts/shortcodes/race-viz.html`, loaded page-scoped MapLibre assets through `layouts/partials/head/css/race-viz.html` and `layouts/partials/head/js/race-viz.html`, initialized the published chart style from `assets/js/race-viz.js`, and covered the embeddable map contract in `race_map_component_test.go` plus updated asset gating in `race_asset_strategy_test.go`.

Deliverables:
- embeddable map container
- map initialization code
- style loading from the existing vector tile configuration

Acceptance criteria:
- one race page can render the chart map reliably in local development

### Task 2.3
Resolve static-hosting implications of vector tiles.

- [x] Done on 2026-05-11: documented the hosting workflow in `docs/race-tile-hosting-strategy.md`, published a machine-readable tile contract via `assets/race-viz/map/tile-manifest.json.tmpl` and `layouts/partials/race-viz/map-foundation.html`, exposed the hosting metadata on `race-viz` embeds, and added `race_tile_hosting_strategy_test.go`.

Deliverables:
- explicit V1 approach for tile serving or tile publishing
- local preview workflow
- production deployment assumption
- note on how the chosen approach can expand to chart coverage beyond Southern California without redesigning the app

Acceptance criteria:
- the chosen approach is simple enough to repeat
- the project is not forced into raster-only charts

## Milestone 3: Course Rendering

### Task 3.1
Render course elements on the map.

- [x] Done on 2026-05-11: loaded `course.json` into the shared MapLibre bootstrap, rendered connected course geometry plus dedicated mark/start-finish layers from `assets/js/race-viz.js`, exposed the stable layer contract in `layouts/shortcodes/race-viz.html`, and covered the slice with `race_course_rendering_test.go`.

Deliverables:
- marks
- start line
- finish line
- connecting course geometry

Acceptance criteria:
- a viewer can visually understand the intended race course from the map alone

### Task 3.2
Support course styling consistent with the site aesthetic.

- [x] Done on 2026-05-11: added the named `signal-v1` course palette and label-layer contract in `layouts/shortcodes/race-viz.html`, upgraded `assets/js/race-viz.js` with glow/casing/label treatment for course overlays, refined the editorial framing in `assets/css/race-viz.css`, and covered the styling slice with updated `race_course_rendering_test.go`.

Deliverables:
- route/mark visual treatment
- dark retro-digital overlay palette aligned with the site style

Acceptance criteria:
- overlays are legible on the nautical chart background
- styling feels intentional rather than generic

### Task 3.3
Implement the V1 land-crossing fallback.

- [x] Done on 2026-05-11: added `controlPointsToNext` support to the shared course route expansion in `assets/js/race-viz.js`, documented the authored fallback in `docs/race-course-schema.md`, and covered both straight-leg and manual-shaping behavior in `race_course_land_fallback_test.go`.

Deliverables:
- support for manual intermediate control points in rendered course geometry

Acceptance criteria:
- a course can be drawn around islands without visually crossing land when manual shaping points are provided

### Task 3.4
Implement standalone `course` mode as a first-class view.

- [x] Done on 2026-05-11: allowed `layouts/shortcodes/race-viz.html` to render `mode="course"` without `boats.json`, kept the shared-layer contract in `assets/js/race-viz.js`, and covered both course-only embeds and shared-engine behavior in `race_embed_shortcode_test.go` plus `race_visualization_architecture_test.go`.

Deliverables:
- embeddable/view-only course visualization using the shared engine
- rendering path that requires course data but no boat track data
- mode selection or configuration path shared with replay embeds

Acceptance criteria:
- a page can embed the visualization in `course` mode without providing track JSON
- the same shared component system is used for both `course` and `replay`

## Milestone 4: Replay Engine

### Task 4.1
Render static boat tracks for all boats in a race.

- [x] Done on 2026-05-11: loaded `boats.json` into the shared replay bootstrap, rendered per-boat static track layers plus sidebar legend hooks from `assets/js/race-viz.js` and `layouts/shortcodes/race-viz.html`, refined the fleet presentation in `assets/css/race-viz.css`, and covered the replay track contract with `race_replay_tracks_test.go`.

Deliverables:
- polylines for self and competitors
- per-boat color support
- legend/sidebar integration hooks

Acceptance criteria:
- all included tracks can be shown on the course simultaneously

### Task 4.2
Implement time-based replay state.

- [x] Done on 2026-05-11: normalized replay time bounds and per-boat interpolation in `assets/js/race-viz.js`, exposed stable replay clock state through `layouts/shortcodes/race-viz.html`, and covered the clock contract plus interpolation-ready sample data in `race_replay_tracks_test.go`.

Deliverables:
- normalized replay clock
- time bounds based on track data
- interpolation between recorded points

Acceptance criteria:
- boats move smoothly enough for replay even when raw points are not extremely dense

### Task 4.3
Implement replay controls.

- [x] Done on 2026-05-11: added a replay control panel in `layouts/shortcodes/race-viz.html`, wired play/pause, speed, scrub, and reset behavior through the shared clock state in `assets/js/race-viz.js`, refined the control styling in `assets/css/race-viz.css`, and covered the control contract in `race_replay_tracks_test.go`.

Deliverables:
- play/pause
- speed selector
- timeline scrubber
- reset to start

Acceptance criteria:
- a user can start at time 0 and replay the race at accelerated speed

### Task 4.4
Implement initial load behavior.

- [x] Done on 2026-05-11: initialized replay at time 0 in `loadBoats` via `state.replay.currentTimeMs = timeline.startTimeMs`, added `enterPrePlayMode` which filters the static track layers to show only the isSelf boat’s full route before playback begins, and covered the pre-play contract in `race_initial_load_test.go`.

Deliverables:
- replay initialized at time 0
- optional pre-play display of the completed `isSelf` route before user playback

Acceptance criteria:
- race pages load in a sensible pre-replay state

### Task 4.5
Render moving boats and persistent full track tails.

- [x] Done on 2026-05-11: added `buildTrackTailCoordinates` and `buildReplayTailFeatures` to compute per-boat tails trimmed to the current replay time, added `buildBoatMarkerFeatures` and `renderBoatMarkerLayers` to show current-position circles, wired both into `renderReplayFrame` which is called from `setReplayTime` on every clock tick, and covered the moving-boats contract in `race_moving_boats_test.go`.

Deliverables:
- current boat markers
- full historical trail through current replay time

Acceptance criteria:
- each boat’s progress is readable during playback

## Milestone 5: Interaction Layer

### Task 5.1
Implement boat legend and visibility toggles.

- [x] Done on 2026-05-11: added per-boat `hiddenBoatIds` Set to `createRaceVizState`, updated `renderBoatLegend` to include a toggle button per item, added `syncBoatLegendVisibility`, `attachBoatLegendToggles`, and `applyBoatVisibilityToLayers` in `assets/js/race-viz.js`, passed `hiddenBoatIds` through `renderReplayFrame` to `buildReplayTailFeatures` and `buildBoatMarkerFeatures`, added toggle button and hidden-state CSS to `assets/css/race-viz.css`, and covered the contract in `race_boat_legend_test.go`.

Deliverables:
- visible boat list
- color keys
- show/hide per boat

Acceptance criteria:
- a viewer can reduce clutter and focus on selected boats

### Task 5.2
Implement hover interactions.

- [x] Done on 2026-05-11: added `hover: { activeTooltip }` to `createRaceVizState`, added `attachBoatMarkerHoverInteractions` (cursor pointer + popup with boat name and formatted replay time on mouseenter, cleanup on mouseleave) and `attachTrackHoverInteractions` (cursor crosshair + boat name popup for static tracks and replay tails) in `assets/js/race-viz.js`; added `.race-viz-hover-tooltip` and related CSS in `assets/css/race-viz.css`; wired both after `renderBoatMarkerLayers` in `loadBoats`; and covered the contract in `race_hover_interactions_test.go`.

Deliverables:
- timestamp tooltip behavior
- useful hover state on track or replay marker

Acceptance criteria:
- users can inspect where a boat was at a given time in the replay

### Task 5.3
Implement event annotations.

- [x] Done on 2026-05-11: added `buildEventFeatures` (resolves positions from lat/lon or interpolates from isSelf track), `renderEventLayers` (amber circle halo/fill + text label layers), `attachEventInteractions` (cursor change on hover, maplibregl.Popup click handler), `setEventsState`, and `loadEvents` (waits for boats before building features so timeline is available for interpolation) in `assets/js/race-viz.js`; added event popup CSS in `assets/css/race-viz.css`; added `data-race-viz-events-state="idle"` to `layouts/shortcodes/race-viz.html`; and covered the contract in `race_event_annotations_test.go`.

Deliverables:
- annotation markers or callouts on the map
- event hover/click treatment

Acceptance criteria:
- notable race moments can be visually highlighted on the map

## Milestone 6: Hugo Page Integration

### Task 6.1
Build the shortcode or embed wrapper for race visualizations.

- [x] Done on 2026-05-11: built the `race-viz` shortcode in `layouts/shortcodes/race-viz.html`, supported race lookup by bundle ID, and covered canonical race-page and prose-page embeds in `race_embed_shortcode_test.go`.

Deliverables:
- Hugo integration that references a race dataset by ID
- page-side mounting point for the visualization

Acceptance criteria:
- race visualizations can be embedded into canonical race pages and other prose pages

### Task 6.2
Create the first canonical race page layout.

- [x] Done on 2026-05-11: added `layouts/races/single.html` which always loads race-viz assets and wraps race pages in a `data-race-page` article with a compact `race-page-header` above and a full-width `race-page-map` section below; added race-page map-first CSS to `assets/css/race-viz.css` making the stage viewport-height on desktop and narrower on mobile; restructured `content/races/dan-byrne-2025/bishop-rock-race/index.md` to place the shortcode first so the map leads the content; and covered the composition contract in `race_canonical_page_layout_test.go`.

Deliverables:
- map-first page composition
- prose section placement below or secondary to the visualization

Acceptance criteria:
- on desktop the map is the primary draw
- on mobile the map dominates first and text is reachable by scrolling

### Task 6.3
Ensure the visualization can coexist with ordinary blog content.

- [x] Done on 2026-05-11: kept race assets page-scoped in `layouts/_default/single.html` plus `layouts/partials/head/css/race-viz.html` and `layouts/partials/head/js/race-viz.html`, and verified race pages load the assets while ordinary pages stay clean in `race_asset_strategy_test.go`.

Deliverables:
- non-race posts remain unaffected
- race-specific assets load only where needed

Acceptance criteria:
- the broader Hugo site still functions as a normal blog

## Milestone 7: Styling And UX Refinement

### Task 7.1
Extend site-level styling to support the intended visual direction.

- [x] Done on 2026-05-11: created `assets/css/custom.css` (overrides theme's placeholder) with `--ac-dark: #7ef5ec` (aligns dark-mode accent with race-viz cyan), `--site-navy: #08111f` (deep navy body background), and shared identity tokens (`--site-cyan`, `--site-cyan-dim`, `--site-border-subtle`); added dark-mode retro-digital gradient treatment on the site header, race-page header, and article headings; covered the contract in `race_site_styling_test.go`.

Deliverables:
- dark theme refinements
- retro-digital gradient accents
- visual consistency between prose pages and race overlays

Acceptance criteria:
- the visualization feels native to the Topaz Racing site identity

### Task 7.2
Tune mobile usability.

- [x] Done on 2026-05-11: added a dedicated mobile refinements section to `assets/css/race-viz.css` under `max-width: 42rem` raising button and boat-toggle touch targets to 2.75rem (44px WCAG minimum), capping the fleet legend with `max-height` + `overflow-y: auto`, constraining event and hover popup widths to 14rem, and reducing race-page header padding on narrow screens; covered the contract in `race_site_styling_test.go`.

Deliverables:
- map sizing behavior
- control layout adjustments
- legend and annotation usability on narrow screens

Acceptance criteria:
- replay remains usable on a phone without breaking the map-first design

## Milestone 8: Documentation And Authoring Workflow

### Task 8.1
Document the race authoring workflow.

- [x] Done on 2026-05-11: documented the end-to-end race authoring flow in `docs/race-authoring-workflow.md`, covering directory setup, course/boats/events JSON authoring, the shortcode embed syntax, and the reference race pointer; covered the authoring contract in `race_authoring_workflow_test.go`.

Deliverables:
- how to create a new race
- where to put course JSON
- where to put boat JSON
- how to embed the visualization in markdown

Acceptance criteria:
- a new race can be authored without reverse-engineering code

### Task 8.2
Document the local preview and deployment workflow for map assets.

- [x] Done on 2026-05-11: documented the local Martin + Hugo preview workflow, production tile URL contract, Hugo-published map artifacts, and future boundaries for tile-generation and race-data import pipelines in `docs/race-preview-deployment.md`; covered the map stack contract in `race_preview_deployment_test.go`.

Deliverables:
- local instructions for previewing the chart map
- production assumptions for static hosting and tile access
- future-facing note describing boundaries for later chart-tile generation and normalized race-data import pipelines

Acceptance criteria:
- the map stack is understandable and repeatable
- future asset-generation or normalization tooling can target documented interfaces instead of requiring visualization rewrites

## Milestone 9: Import Tooling

### Task 9.1
Implement GPX-to-boats.json import tool.

- [x] Done on 2026-05-11: implemented `gpximport` package (`gpximport/converter.go`) with `ConvertGPX` and `MergeBoat`, added CLI at `cmd/gpx-import/main.go`, and covered all conversion and merge cases in `gpximport/converter_test.go` (14 tests).

Deliverables:
- `gpximport` Go package with `ConvertGPX` and `MergeBoat` functions
- `cmd/gpx-import` CLI accepting `--id`, `--name`, `--color`, `--boat-type`, `--self`, `--merge`, and `--output` flags
- Supports GPX 1.0 and 1.1, multi-segment, and multi-track files
- Merge flag allows adding a new boat to an existing `boats.json`

Acceptance criteria:
- a GPX file from a Garmin or phone can be converted to `boats.json` format without hand-authoring track points
- multiple boats can be accumulated into one `boats.json` by running the tool once per boat with `--merge`

## Deferred Backlog
These are explicitly out of scope for initial implementation, but should remain visible for future work:

- live race mode with periodic refresh
- import/conversion tooling for Jibeset and Garmin Connect export formats
- auto-routing around land/islands
- weather overlays
- exclusion zone rendering and logic
- synchronized prose-to-map storytelling effects
- offshore/remote posting workflows

## Milestone Checklist
- [x] Milestone 0: Repository And Architecture Baseline
- [x] Milestone 1: Data Contracts
- [x] Milestone 2: Map Foundation
- [x] Milestone 3: Course Rendering
- [x] Milestone 4: Replay Engine
- [x] Milestone 5: Interaction Layer
- [x] Milestone 6: Hugo Page Integration
- [x] Milestone 7: Styling And UX Refinement
- [x] Milestone 8: Documentation And Authoring Workflow
- [x] Milestone 9: Import Tooling

## Definition Of Done For V1
V1 is complete when:
- a canonical race page can be created in Hugo
- that page can render the ENC-style chart map
- a race course can be loaded from JSON in standalone `course` mode
- multiple boats can be loaded from JSON
- the race can be replayed with controls
- boats can be shown/hidden from a legend
- event annotations can be displayed
- `course` and `replay` share one visualization engine and data model
- the page remains mostly static-hostable
- the workflow for adding another race is documented
