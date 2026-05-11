# PROMPT.md

## Project
Topaz Racing is a Hugo-based personal website and blog for Loren Brindze's sailing and racing campaigns. It includes ordinary blog content plus a custom race visualization system built around nautical chart maps and replayable race tracks.

The existing repo already contains:
- a Hugo site for general content and posts
- an early custom vector-tile marine map prototype under `tiles/`
- ENC-style chart rendering work that should be reused rather than replaced

## Primary Goal
Build a reusable race visualization system that can be embedded into Hugo content and also serve as the main feature of dedicated race pages.

V1 should allow Loren to:
- add a new race page
- define a race course on top of the custom marine chart map
- add one or more boat tracks as hand-authored JSON
- replay the race over time with animation and controls
- embed that visualization inside longer prose content when desired

## Non-Goals For V1
Do not implement yet:
- live race ingestion or polling
- remote authoring by email or offshore workflows
- weather overlays
- automatic data import/conversion tools
- full commentary synchronization with scroll position
- exclusion zones as a rendered feature
- automatic route-around-land if it becomes too complex for first pass

These should be considered in the design so V1 does not block them later.

## Content Model
The site supports two broad content types:
- normal blog posts not tied to races
- race pages that combine prose with a custom map-based visualization

Each race should have a canonical Hugo page, likely `index.md`, which owns both the prose and the embedded visualization.

Race-related data should live under `content/` for now.

## Visualization System
Create a reusable race visualization component system with a shared data model.

V1 modes:
- `course`: show chart background plus race course geometry
- `replay`: show course plus time-based animated boat tracks

Future mode:
- `live`: reuse the replay-oriented track model and map UI, but update positions from external sources

The preferred implementation is a shared visualization engine with layered behavior rather than completely separate unrelated components.

## Map Requirements
Use the existing vector-tile marine chart approach already prototyped in `tiles/` as the baseline.

Requirements:
- preserve ENC-style marine chart direction
- keep published site as static as possible
- allow offline local authoring/preview as much as practical
- prefer vector tiles over raster tiles
- allow future expansion to additional regions beyond Southern California

Implementation flexibility:
- static hosting is preferred
- pre-generated static vector tile assets are acceptable
- if needed, a tile server can be used later
- V1 should avoid locking the project into raster-only charts

## Race Course Model
A race course is an ordered set of course elements rendered over a chart domain.

V1 course features:
- one ordered course per race
- waypoints/marks with lat/lon
- optional `name`
- rounding direction support
- start/finish line supported as a special type

Initial course element types:
- `mark`
- `start_line`
- `finish_line`

Expected waypoint fields:
- `id`
- `type`
- `lat`
- `lon`
- `name` optional
- `rounding` with values like `port`, `starboard`, `none`

The course should be rendered as connected geometry between relevant points.

## Land Handling
Preferred long-term behavior:
- route/course lines should visually go around islands or land rather than intersecting land

V1 fallback:
- support manual intermediate control points if automatic land-aware routing is not practical
- data model should leave room for future automatic route-around-land behavior
- manual route-shaping points may be treated as implementation details rather than first-class user-facing marks

## Boat Track Model
Boat tracks are time-indexed sequences of positions.

Per-point shape:

```json
{ "time": "2025-02-11T18:35:00Z", "lat": 33.123, "lon": -118.456 }
```

Each boat should also support metadata:
- `id`
- `name`
- `color`
- `boatType`
- `source`
- `isSelf`

Tracks will be hand-authored JSON in V1.

The system should be designed so future import tools can normalize Jibeset, Garmin, GPX, or other sources into this internal format.

## Replay Behavior
Replay mode is required in V1.

Required behavior:
- animate boats over time from track data
- provide play/pause
- provide speed control
- provide timeline scrubber
- allow show/hide per boat
- provide boat legend/sidebar with names and colors
- leave full track tails behind boats during replay in V1
- support competitor boats in addition to Loren's boat

Initial load behavior:
- load at time 0
- optionally show only the full route for `isSelf == true` before playback begins
- user then starts replay to animate all boats

Interpolation between recorded points should be smooth enough for animation and not require extremely dense source data.

## Interaction Requirements
V1 interactions:
- hover tooltip with timestamp information
- visible boat list / legend
- hide/show boats
- replay controls
- event annotations rendered on the map

Event annotations should support future storytelling moments such as:
- gybe
- wipeout
- sail change
- other notable race events

Commentary callouts synchronized with prose scrolling are not required in V1, but the design should not prevent later text-map coordination.

## Page Composition
Race pages should make the map the primary visual feature.

Desktop intent:
- map is the main draw
- prose is secondary but still part of the page

Mobile intent:
- map dominates first
- reader scrolls down for text
- more advanced tab/slideout patterns may come later

The visualization must also be embeddable inside other prose-driven pages or posts when needed.

## Authoring Model
V1 authoring is fully local and manual.

Assumptions:
- Loren writes Hugo markdown locally
- race/course/track data is authored as JSON
- embedding should be by reference rather than inlining raw data into markdown

Preferred embed shape:
- a reusable shortcode or equivalent reference such as a race ID

The exact Hugo integration mechanism should optimize for maintainability and clean content organization.

## Accessibility And UX
The site visual direction should align with Loren's intended aesthetic:
- dark themed
- retro-digital / 90s gradient influence
- nautical chart background preserved
- race overlays should feel intentional and editorial, not generic dashboard UI

The map should remain usable on both desktop and mobile, even if desktop gets the richer initial experience.

## Future Considerations
Design for future support of:
- live race mode with periodic data refresh
- data ingestion from Jibeset and/or Garmin
- weather overlays
- exclusion zones
- text-map synchronized storytelling
- offshore/remote posting workflows
- static asset generation pipelines for chart tiles and normalized race data

## Success Criteria For V1
V1 is successful when Loren can:
- create a new race page in Hugo
- define a course on the existing chart system
- add his own and competitors' boat tracks as JSON
- publish a page where the race can be replayed with controls
- use that same visualization as the centerpiece of a race page with supporting prose
- keep the published site mostly static
