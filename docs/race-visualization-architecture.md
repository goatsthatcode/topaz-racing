# Race Visualization Architecture

## Decision
V1 uses one shared race visualization engine for both `course` and `replay` embeds.

The engine contract is defined by:
- one root embed element with race identity, mode, and asset URLs
- one shared stage scaffold with stable layer slots
- one shared runtime state object that owns data loading, map lifecycle, and replay time

`course` and `replay` are configuration modes on top of the same engine. They are not separate applications.

## Shared Primitives

### Map primitive
The map primitive owns:
- chart style loading
- viewport state
- future vector-tile source registration
- the base rendering surface used by every mode

This primitive exists even when only course geometry is shown.

### Overlay primitives
The engine reserves stable layer slots in render order:
1. `map`
2. `course`
3. `tracks`
4. `boats`
5. `events`
6. `controls`

The stage scaffold is shared across modes so later milestones can add rendering behavior without changing the embed contract.

### State primitive
The shared state object owns:
- embed metadata such as race ID and mode
- URLs for course, boat, and event JSON
- loaded payload caches
- replay clock state
- enabled layer list for the current mode

Mode-specific behavior reads from this state instead of inventing separate state trees.

## Mode Layering

### `course`
`course` mode enables:
- `map`
- `course`
- `events`

It does not require boat-track rendering or replay controls.

### `replay`
`replay` mode enables:
- `map`
- `course`
- `tracks`
- `boats`
- `events`
- `controls`

Replay-specific behavior extends the shared state with time controls and interpolation, but still renders inside the same stage and layer order.

## Boundaries
- The shortcode is responsible only for emitting the engine contract and resource URLs.
- The shared engine is responsible for parsing config, creating shared state, reserving layer slots, and dispatching mode-specific activation.
- Future map adapters, course renderers, and replay controls plug into the shared engine instead of owning their own bootstrap path.

## Consequences
- Milestone 2 can build the reusable map component once and attach it to the `map` layer.
- Milestone 3 can add course rendering without needing replay data.
- Milestone 4 can add replay behavior by activating more layers and extending shared state, not by replacing the course implementation.
