# Race Tile Hosting Strategy

Task 2.3 fixes the remaining operational gap in the map foundation: race pages now publish an explicit tile-hosting contract for local preview and production builds instead of assuming tile serving is configured out of band.

## V1 Hosting Decision

- V1 keeps vector tiles outside the Hugo output tree and treats them as a dedicated static service.
- Hugo publishes two map-side artifacts for race pages:
- `race-viz/map-style.json` for the MapLibre style that points at the active tile endpoint.
- `race-viz/tile-manifest.json` for the machine-readable tile-serving contract, including the default tile set and coverage metadata.
- The current serving mode is `external-static-vector-host`.

This preserves a mostly static site while avoiding a premature in-app tile server.

## Local Preview Workflow

- Run Hugo normally for page rendering.
- Run `martin --config tiles/martin-config` to serve the local MBTiles catalog on `http://127.0.0.1:3000`.
- In development builds, race pages and the published map style point at that local Martin endpoint.

The race embed therefore reads the same Hugo-published style and manifest contract in local preview that it will use in production, with only the tile host changing by environment.

## Production Assumption

- Production pages point at `https://topaz-racing.com/tiles`.
- The production tile host must expose tile paths in the form `/{tileset}/{z}/{x}/{y}`.
- The default V1 race visualization uses the `combined_socal` tile set defined in site config.
- The tile host can be a CDN-backed static endpoint, object storage fronted by a tile service, or another deployment that preserves the same URL contract.

V1 therefore commits to a stable URL shape instead of a specific infrastructure vendor.

## Coverage Expansion

- Additional chart regions should be added as new MBTiles artifacts under `tiles/mbtiles/`.
- Each new region should get a corresponding entry in `params.raceViz.tiles.sets` with its path, bounds, center, and source metadata.
- The Hugo-published tile manifest is the forward-compatible place for the frontend to discover available coverage beyond Southern California.

This avoids redesigning the app when the project expands past the initial SoCal chart bundle.
