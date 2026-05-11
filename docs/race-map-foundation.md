# Race Map Foundation

Task 2.1 promotes the existing `tiles/` prototype into a Hugo-facing contract without making race pages depend on `tiles/index.html`.

## Relationship To The Prototype

- `tiles/index.html` remains the quick standalone sandbox for manual style inspection.
- `tiles/style.json` remains the human-edited prototype source that describes the ENC-style vector layer stack.
- `assets/race-viz/map/style.json.tmpl` is the site-owned publishing template derived from `tiles/style.json`.
- `layouts/partials/race-viz/map-foundation.html` is the Hugo integration point that publishes the style JSON and exposes its URL to race embeds.

The race visualization should consume the published Hugo style resource, not the standalone prototype page.

## Site Integration Contract

Every `race-viz` embed now emits:

- `data-race-viz-map-style-url` for the published MapLibre style JSON
- `data-race-viz-map-tile-endpoint` for the vector tile host prefix
- `data-race-viz-map-prototype-page` and `data-race-viz-map-prototype-style` so the provenance from `tiles/` stays explicit during development

This keeps the map foundation owned by the site layer and gives later map initialization code a stable contract.

## Environment Decision

Local preview and published builds use the same style template but different tile endpoints:

- local preview (`hugo server` or `HUGO_ENVIRONMENT=development`) points at `http://127.0.0.1:3000`
- published builds point at `https://topaz-racing.com/tiles`

That means the Hugo page loads the same race visualization contract in both environments while leaving tile serving itself to the dedicated endpoint behind the style JSON.

## Why This Is Enough For Task 2.1

The map can now be initialized inside a Hugo page by reading the race embed contract and fetching the published style JSON. The standalone prototype page is no longer required for site integration, even though it remains useful for isolated chart-style iteration.
