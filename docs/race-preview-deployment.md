# Local Preview And Deployment Workflow

This document covers how to run the full race visualization stack locally, what production hosting assumes, and where future asset-generation pipelines should plug in.

---

## Local Preview

Race pages require two servers running concurrently: Hugo's built-in dev server for the site, and Martin for the vector chart tiles.

### 1. Start The Tile Server

Vector tiles are served from the MBTiles catalog in `tiles/mbtiles/`. Martin reads its configuration from `tiles/martin-config`:

```sh
martin --config tiles/martin-config
```

Martin binds to `http://127.0.0.1:3000` by default. The Hugo site config in `config/_default/hugo.toml` tells race pages to use that address as the tile endpoint in development builds.

If Martin is not available, install it via:

```sh
cargo install martin
# or via Homebrew:
brew install martin
```

### 2. Start The Hugo Dev Server

In a second terminal from the repo root:

```sh
hugo server
```

The dev server automatically uses `HUGO_ENVIRONMENT=development`, which switches the race embed tile endpoint from the production URL to `http://127.0.0.1:3000`.

### 3. Open A Race Page

Navigate to any race page, for example:

```
http://localhost:1313/races/dan-byrne-2025/bishop-rock-race/
```

The map should render the ENC-style nautical chart and load the course and boat track overlays.

### Tile Server Not Running

If Martin is not running, the map background will fail to load while course and boat overlays may still render. The race embed emits a `data-race-viz-map-fallback-tile-endpoint` attribute so the visualization can surface a warning rather than silently failing. For map style development without the tile server, `tiles/index.html` remains the quickest standalone sandbox.

---

## Chart Tile Catalog

The local tile catalog lives in `tiles/mbtiles/`. The default tile set for V1 race pages is `combined_socal.mbtiles`, which covers Southern California ENC chart data.

Additional MBTiles files in that directory are available as named tile sets. The active tile set for any race embed is `combined_socal` unless overridden in Hugo site config under `params.raceViz.tiles.sets`.

To add chart coverage for a new region, place the MBTiles file in `tiles/mbtiles/` and add a corresponding entry in `params.raceViz.tiles.sets`:

```toml
[[params.raceViz.tiles.sets]]
    id      = 'socal_north'
    title   = 'SoCal North Coverage'
    path    = 'socal_north'
    region  = 'southern-california-north'
    source  = 'tiles/mbtiles/socal_north.mbtiles'
    bounds  = [-120.5, 33.5, -117.0, 35.5]
    center  = [-118.8, 34.5]
    zoom    = 9
```

The Hugo-published `race-viz/tile-manifest.json` is the forward-compatible place for the frontend to discover available tile sets. Adding a new set here does not require any changes to visualization code.

---

## Production Hosting

The published site at `https://topaz-racing.com/` treats vector tiles as a separate static service rather than Hugo output.

### URL Contract

Production race pages point their tile endpoint at `https://topaz-racing.com/tiles`. The tile host must serve tiles at paths of the form:

```
/{tileset}/{z}/{x}/{y}
```

For example, the default Southern California coverage is served from:

```
/tiles/combined_socal/{z}/{x}/{y}
```

### Tile Deployment

V1 does not generate tiles as part of the Hugo build. Tiles must be deployed independently to the production host. The deployment mechanism can be any static-file host, CDN-backed object storage, or tile service that exposes the URL contract above.

The production tile host is intentionally separate from the Hugo static output tree so tile files (which can be large) do not bloat the site build.

### Hugo Build

The standard Hugo build command:

```sh
hugo
```

Produces static HTML, CSS, JS, and the published map-style JSON and tile-manifest JSON under `public/`. No tile generation or tile copying happens during the Hugo build.

---

## Published Map Artifacts

Hugo emits two tile-related artifacts on race pages that keep the frontend decoupled from hard-coded configuration:

| Path | Purpose |
|------|---------|
| `race-viz/map-style.json` | MapLibre style pointing at the active tile endpoint |
| `race-viz/tile-manifest.json` | Machine-readable hosting contract: serving mode, default tile set, available sets, coverage metadata |

The visualization reads these at runtime rather than embedding configuration directly in JavaScript. This means switching tile endpoints or adding coverage only requires a Hugo config change, not a code change.

---

## Future Asset-Generation Pipelines

The following pipelines are out of scope for V1 but should target the interfaces documented here rather than requiring visualization rewrites.

### Chart Tile Generation

Future tooling that generates ENC-derived vector tile files should produce MBTiles artifacts compatible with Martin and deposit them in `tiles/mbtiles/`. The Martin config at `tiles/martin-config` is the integration point for exposing new tile sets locally.

The tile manifest contract (`race-viz/tile-manifest.json`) and the `params.raceViz.tiles.sets` config table are the stable interfaces for registering new coverage with the frontend. A tile generator does not need to touch visualization code — it only needs to produce a valid MBTiles file and add an entry to hugo.toml.

### Race Data Import

Future import tooling that converts Jibeset, Garmin, GPX, or other sources into the Topaz Racing internal format should target the `schemas/race-boats-v1.schema.json` contract. The output is a `boats.json` file dropped into the target race bundle directory.

The `source` field on each boat object (currently `"hand-authored"`) is the designated place for import provenance metadata. An import pipeline can set `"source": "garmin-fit"` or `"source": "jibeset"` without changing the visualization or the schema structure.

The visualization engine does not need modification when new import sources appear — it reads the stable internal contract regardless of where the data originated.
