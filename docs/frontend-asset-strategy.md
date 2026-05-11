# Race Visualization Frontend Asset Strategy

## Decision
Keep race visualization frontend code in project-owned Hugo Pipes assets and load it only on single pages whose rendered content contains a `race-viz` embed.

Selected locations:
- `assets/js/race-viz.js` for the shared visualization bootstrap and later shared engine code
- `assets/css/race-viz.css` for visualization-specific styling
- `layouts/partials/head/js/race-viz.html` and `layouts/partials/head/css/race-viz.html` for asset publication
- `layouts/_default/single.html` for page-scoped loading logic

## Loading Strategy
- The `race-viz` shortcode remains responsible only for stable DOM and data URLs.
- The single-page template checks rendered `.Content` for `data-race-viz`.
- When a page includes at least one visualization embed, Hugo emits the race-specific CSS and JS in the page head.
- Pages without a visualization embed do not load the race frontend bundle.

This keeps the visualization code independent from vendor theme asset bundles while still using Hugo Pipes for minification and fingerprinting.

## Why This Shape
- Avoids editing theme-owned asset concatenation for app code that belongs to the project.
- Works for canonical race pages and for ordinary prose pages that embed a race by reference.
- Gives the visualization a stable repo-owned entrypoint before the shared engine exists.
- Keeps ordinary posts free of map/replay assets until they actually render an embed.

## Future Expansion
- `assets/js/race-viz.js` is the shared bootstrap entrypoint for both `course` and `replay` modes.
- Mode-specific behavior should branch from DOM data attributes instead of separate page templates.
- If the bundle grows, the same partials can swap to a compiled pipeline without changing authoring or shortcode contracts.
