# TASKS.md

## Purpose
This document turns [SPEC.md](/Users/lorenbrindze/Projects/topaz-racing/SPEC.md:1) into an implementation backlog for the Topaz Racing race visualization system.

The intended execution style is:
- work top-down by priority unless blocked
- keep the site deployable as a mostly static Hugo project
- favor small end-to-end slices over broad unfinished scaffolding
- preserve future support for live tracking, import tooling, and richer storytelling overlays

## Milestone 0: Repository And Architecture Baseline

### Task 0.1
Decide and document the on-disk content structure for race pages and race data under `content/`.

Deliverables:
- one canonical race directory shape
- clear separation between prose content and JSON assets
- one example race directory chosen as the reference pattern

Acceptance criteria:
- a new race can be added by copying a known folder structure
- the structure is compatible with Hugo page rendering and static asset publication

### Task 0.2
Choose the Hugo integration approach for embedding race visualizations.

Deliverables:
- one selected embed mechanism, preferably shortcode-based
- a documented race identifier convention

Acceptance criteria:
- a race page can reference a visualization without inlining raw JSON in markdown

### Task 0.3
Define the frontend asset strategy for the visualization code.

Deliverables:
- chosen location for JS/CSS assets in the Hugo project
- plan for bundling/loading map code in race pages only when needed

Acceptance criteria:
- visualization assets are not tied to theme internals unnecessarily
- the site remains maintainable without editing vendor theme code for core app logic

## Milestone 1: Data Contracts

### Task 1.1
Define the V1 JSON schema for race course data.

Deliverables:
- course file format
- element types for `mark`, `start_line`, and `finish_line`
- support for optional manual route-shaping control points

Acceptance criteria:
- a single race course can be described completely in JSON
- the schema leaves room for future exclusion zones and auto-routing

### Task 1.2
Define the V1 JSON schema for boat track data.

Deliverables:
- per-boat metadata contract
- per-point time/lat/lon contract
- expected file organization for one or more boats in a race

Acceptance criteria:
- one self boat and multiple competitor boats can be described cleanly
- the model is independent of Jibeset, Garmin, GPX, or other upstream formats

### Task 1.3
Define the V1 JSON schema for event annotations.

Deliverables:
- annotation types
- time-anchored and/or position-anchored event model
- optional label/description fields

Acceptance criteria:
- a race can include notable moments like gybes, wipeouts, or sail changes

### Task 1.4
Create one complete sample race dataset.

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

Deliverables:
- documented relationship between `tiles/index.html`, `tiles/style.json`, and site integration
- decision on how the site loads the map style in local preview and published environments

Acceptance criteria:
- the map can render inside the site without depending on the standalone prototype page

### Task 2.2
Create a reusable map component that renders the ENC-style vector chart background in a Hugo page.

Deliverables:
- embeddable map container
- map initialization code
- style loading from the existing vector tile configuration

Acceptance criteria:
- one race page can render the chart map reliably in local development

### Task 2.3
Resolve static-hosting implications of vector tiles.

Deliverables:
- explicit V1 approach for tile serving or tile publishing
- local preview workflow
- production deployment assumption

Acceptance criteria:
- the chosen approach is simple enough to repeat
- the project is not forced into raster-only charts

## Milestone 3: Course Rendering

### Task 3.1
Render course elements on the map.

Deliverables:
- marks
- start line
- finish line
- connecting course geometry

Acceptance criteria:
- a viewer can visually understand the intended race course from the map alone

### Task 3.2
Support course styling consistent with the site aesthetic.

Deliverables:
- route/mark visual treatment
- dark retro-digital overlay palette aligned with the site style

Acceptance criteria:
- overlays are legible on the nautical chart background
- styling feels intentional rather than generic

### Task 3.3
Implement the V1 land-crossing fallback.

Deliverables:
- support for manual intermediate control points in rendered course geometry

Acceptance criteria:
- a course can be drawn around islands without visually crossing land when manual shaping points are provided

## Milestone 4: Replay Engine

### Task 4.1
Render static boat tracks for all boats in a race.

Deliverables:
- polylines for self and competitors
- per-boat color support
- legend/sidebar integration hooks

Acceptance criteria:
- all included tracks can be shown on the course simultaneously

### Task 4.2
Implement time-based replay state.

Deliverables:
- normalized replay clock
- time bounds based on track data
- interpolation between recorded points

Acceptance criteria:
- boats move smoothly enough for replay even when raw points are not extremely dense

### Task 4.3
Implement replay controls.

Deliverables:
- play/pause
- speed selector
- timeline scrubber
- reset to start

Acceptance criteria:
- a user can start at time 0 and replay the race at accelerated speed

### Task 4.4
Implement initial load behavior.

Deliverables:
- replay initialized at time 0
- optional pre-play display of the completed `isSelf` route before user playback

Acceptance criteria:
- race pages load in a sensible pre-replay state

### Task 4.5
Render moving boats and persistent full track tails.

Deliverables:
- current boat markers
- full historical trail through current replay time

Acceptance criteria:
- each boat’s progress is readable during playback

## Milestone 5: Interaction Layer

### Task 5.1
Implement boat legend and visibility toggles.

Deliverables:
- visible boat list
- color keys
- show/hide per boat

Acceptance criteria:
- a viewer can reduce clutter and focus on selected boats

### Task 5.2
Implement hover interactions.

Deliverables:
- timestamp tooltip behavior
- useful hover state on track or replay marker

Acceptance criteria:
- users can inspect where a boat was at a given time in the replay

### Task 5.3
Implement event annotations.

Deliverables:
- annotation markers or callouts on the map
- event hover/click treatment

Acceptance criteria:
- notable race moments can be visually highlighted on the map

## Milestone 6: Hugo Page Integration

### Task 6.1
Build the shortcode or embed wrapper for race visualizations.

Deliverables:
- Hugo integration that references a race dataset by ID
- page-side mounting point for the visualization

Acceptance criteria:
- race visualizations can be embedded into canonical race pages and other prose pages

### Task 6.2
Create the first canonical race page layout.

Deliverables:
- map-first page composition
- prose section placement below or secondary to the visualization

Acceptance criteria:
- on desktop the map is the primary draw
- on mobile the map dominates first and text is reachable by scrolling

### Task 6.3
Ensure the visualization can coexist with ordinary blog content.

Deliverables:
- non-race posts remain unaffected
- race-specific assets load only where needed

Acceptance criteria:
- the broader Hugo site still functions as a normal blog

## Milestone 7: Styling And UX Refinement

### Task 7.1
Extend site-level styling to support the intended visual direction.

Deliverables:
- dark theme refinements
- retro-digital gradient accents
- visual consistency between prose pages and race overlays

Acceptance criteria:
- the visualization feels native to the Topaz Racing site identity

### Task 7.2
Tune mobile usability.

Deliverables:
- map sizing behavior
- control layout adjustments
- legend and annotation usability on narrow screens

Acceptance criteria:
- replay remains usable on a phone without breaking the map-first design

## Milestone 8: Documentation And Authoring Workflow

### Task 8.1
Document the race authoring workflow.

Deliverables:
- how to create a new race
- where to put course JSON
- where to put boat JSON
- how to embed the visualization in markdown

Acceptance criteria:
- a new race can be authored without reverse-engineering code

### Task 8.2
Document the local preview and deployment workflow for map assets.

Deliverables:
- local instructions for previewing the chart map
- production assumptions for static hosting and tile access

Acceptance criteria:
- the map stack is understandable and repeatable

## Deferred Backlog
These are explicitly out of scope for initial implementation, but should remain visible for future work:

- live race mode with periodic refresh
- import/conversion tooling for Jibeset, Garmin, GPX, and similar formats
- auto-routing around land/islands
- weather overlays
- exclusion zone rendering and logic
- synchronized prose-to-map storytelling effects
- offshore/remote posting workflows

## Recommended Execution Order
Implement in this order unless blocked:

1. Milestone 0
2. Milestone 1
3. Milestone 2
4. Milestone 3
5. Milestone 4
6. Milestone 5
7. Milestone 6
8. Milestone 7
9. Milestone 8

## Definition Of Done For V1
V1 is complete when:
- a canonical race page can be created in Hugo
- that page can render the ENC-style chart map
- a race course can be loaded from JSON
- multiple boats can be loaded from JSON
- the race can be replayed with controls
- boats can be shown/hidden from a legend
- event annotations can be displayed
- the page remains mostly static-hostable
- the workflow for adding another race is documented
