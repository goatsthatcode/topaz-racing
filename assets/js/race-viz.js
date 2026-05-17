const DEFAULT_SHARED_LAYERS = ["map", "course", "tracks", "boats", "events", "controls"];
const DEFAULT_COURSE_SOURCE_ID = "race-viz-course";
const DEFAULT_COURSE_ROUTE_LAYER_ID = "race-viz-course-route";
const DEFAULT_COURSE_MARKS_LAYER_ID = "race-viz-course-marks";
const DEFAULT_COURSE_ROUNDING_LAYER_ID = "race-viz-course-marks-rounding";
const DEFAULT_COURSE_START_FINISH_LAYER_ID = "race-viz-course-start-finish";
const DEFAULT_COURSE_LABELS_LAYER_ID = "race-viz-course-labels";
const DEFAULT_TRACKS_SOURCE_ID = "race-viz-tracks";
const DEFAULT_TRACKS_LAYER_ID = "race-viz-tracks";
const DEFAULT_REPLAY_TAILS_SOURCE_ID = "race-viz-replay-tails";
const DEFAULT_REPLAY_TAILS_LAYER_ID = "race-viz-replay-tails";
const DEFAULT_BOAT_MARKERS_SOURCE_ID = "race-viz-boat-markers";
const DEFAULT_BOAT_MARKERS_LAYER_ID = "race-viz-boat-markers";
const DEFAULT_EVENTS_SOURCE_ID = "race-viz-events";
const DEFAULT_EVENTS_LAYER_ID = "race-viz-events";
const DEFAULT_COURSE_PALETTE = "signal-v1";
const DEFAULT_MAP_FIT_PADDING = 48;
const DEFAULT_MAP_FIT_MAX_ZOOM = 7.0;
const DEFAULT_REPLAY_SPEED = 120;
const COURSE_STYLE_PRESETS = {
  "signal-v1": {
    routeCasingColor: "rgba(4, 16, 24, 0.98)",
    routeCasingWidth: 10,
    routeGlowColor: "rgba(83, 232, 218, 0.34)",
    routeGlowWidth: 7,
    routeColor: "rgba(126, 245, 236, 0.98)",
    routeWidth: 4,
    routeDasharray: [1.1, 0.8],
    markColor: "rgba(122, 226, 255, 0.98)",
    markStrokeColor: "rgba(5, 20, 32, 0.98)",
    roundingPortColor: "rgba(230, 90, 65, 0.92)",
    roundingStarboardColor: "rgba(55, 200, 100, 0.92)",
    startColor: "rgba(255, 211, 92, 0.98)",
    finishColor: "rgba(255, 112, 146, 0.98)",
    labelColor: "rgba(220, 248, 255, 0.96)",
    labelHaloColor: "rgba(6, 18, 28, 0.94)",
  },
};

function parseMapMaxBounds(value) {
  if (!value) return null;
  const parts = value.split(",").map(Number);
  if (parts.length === 4 && parts.every((n) => !isNaN(n))) {
    return parts; // [west, south, east, north]
  }
  return null;
}

function parseLayerList(value, fallback) {
  if (!value) {
    return [...fallback];
  }

  return value
    .split(/\s+/)
    .map((entry) => entry.trim())
    .filter(Boolean);
}

function createRaceVizConfig(root) {
  const sharedLayers = parseLayerList(
    root.dataset.raceVizSharedLayers,
    DEFAULT_SHARED_LAYERS,
  );
  const activeLayers = parseLayerList(root.dataset.raceVizActiveLayers, sharedLayers);

  return {
    raceId: root.dataset.raceId ?? "",
    mode: root.dataset.raceMode ?? "replay",
    engine: root.dataset.raceVizEngine ?? "shared-v1",
    sharedLayers,
    activeLayers,
    map: {
      styleURL: resolveRaceVizURL(root.dataset.raceVizMapStyleUrl ?? ""),
      tileEndpoint: root.dataset.raceVizMapTileEndpoint ?? "",
      fallbackTileEndpoint: root.dataset.raceVizMapFallbackTileEndpoint ?? "",
      tileManifestURL: resolveRaceVizURL(root.dataset.raceVizMapTileManifestUrl ?? ""),
      tileSet: root.dataset.raceVizMapTileSet ?? "",
      previewCommand: root.dataset.raceVizMapPreviewCommand ?? "",
      servingMode: root.dataset.raceVizMapServingMode ?? "",
      prototypePage: root.dataset.raceVizMapPrototypePage ?? "",
      prototypeStyle: root.dataset.raceVizMapPrototypeStyle ?? "",
      maxBounds: parseMapMaxBounds(root.dataset.raceVizMapMaxBounds ?? ""),
      minZoom: root.dataset.raceVizMapMinZoom !== undefined
        ? parseFloat(root.dataset.raceVizMapMinZoom)
        : null,
      maxZoom: root.dataset.raceVizMapMaxZoom !== undefined
        ? parseFloat(root.dataset.raceVizMapMaxZoom)
        : null,
    },
    course: {
      url: resolveRaceVizURL(root.dataset.courseUrl ?? ""),
      sourceID: root.dataset.raceVizCourseSource ?? DEFAULT_COURSE_SOURCE_ID,
      routeLayerID:
        root.dataset.raceVizCourseRouteLayer ?? DEFAULT_COURSE_ROUTE_LAYER_ID,
      marksLayerID:
        root.dataset.raceVizCourseMarksLayer ?? DEFAULT_COURSE_MARKS_LAYER_ID,
      roundingLayerID:
        root.dataset.raceVizCourseRoundingLayer ?? DEFAULT_COURSE_ROUNDING_LAYER_ID,
      startFinishLayerID:
        root.dataset.raceVizCourseStartFinishLayer ?? DEFAULT_COURSE_START_FINISH_LAYER_ID,
      labelsLayerID:
        root.dataset.raceVizCourseLabelsLayer ?? DEFAULT_COURSE_LABELS_LAYER_ID,
      palette: root.dataset.raceVizCoursePalette ?? DEFAULT_COURSE_PALETTE,
    },
    tracks: {
      sourceID: root.dataset.raceVizTracksSource ?? DEFAULT_TRACKS_SOURCE_ID,
      layerID: root.dataset.raceVizTracksLayer ?? DEFAULT_TRACKS_LAYER_ID,
    },
    replayTails: {
      sourceID: root.dataset.raceVizReplayTailsSource ?? DEFAULT_REPLAY_TAILS_SOURCE_ID,
      layerID: root.dataset.raceVizReplayTailsLayer ?? DEFAULT_REPLAY_TAILS_LAYER_ID,
    },
    boatMarkers: {
      sourceID: root.dataset.raceVizBoatMarkersSource ?? DEFAULT_BOAT_MARKERS_SOURCE_ID,
      layerID: root.dataset.raceVizBoatMarkersLayer ?? DEFAULT_BOAT_MARKERS_LAYER_ID,
    },
    events: {
      sourceID: root.dataset.raceVizEventsSource ?? DEFAULT_EVENTS_SOURCE_ID,
      layerID: root.dataset.raceVizEventsLayer ?? DEFAULT_EVENTS_LAYER_ID,
    },
    boatsURL: resolveRaceVizURL(root.dataset.boatsUrl ?? ""),
    eventsURL: resolveRaceVizURL(root.dataset.eventsUrl ?? ""),
    replaySpeed: parseInt(root.dataset.raceVizReplaySpeed ?? String(DEFAULT_REPLAY_SPEED), 10) || DEFAULT_REPLAY_SPEED,
    fitMaxZoom: parseFloat(root.dataset.raceVizFitMaxZoom ?? "") || DEFAULT_MAP_FIT_MAX_ZOOM,
  };
}

function resolveRaceVizURL(value) {
  if (!value) {
    return "";
  }

  if (/^[a-z]+:\/\//i.test(value) || value.startsWith("//")) {
    return value;
  }

  return new URL(`/${value.replace(/^\/+/, "")}`, window.location.origin).toString();
}

function getCourseStyle(config) {
  return COURSE_STYLE_PRESETS[config.course.palette] ?? COURSE_STYLE_PRESETS[DEFAULT_COURSE_PALETTE];
}

function createRaceVizState(config) {
  return {
    config,
    data: {
      course: null,
      boats: null,
      events: null,
    },
    course: {
      status: "idle",
    },
    boats: {
      status: "idle",
    },
    map: {
      instance: null,
      status: "idle",
    },
    replay: {
      status: "idle",
      startTime: "",
      endTime: "",
      startTimeMs: 0,
      endTimeMs: 0,
      durationMs: 0,
      currentTimeMs: 0,
      playing: false,
      started: false,
      speed: config.replaySpeed ?? 60,
      timeline: null,
      snapshot: null,
      animationFrameID: 0,
      lastFrameMs: 0,
    },
    visibility: {
      hiddenBoatIds: new Set(),
    },
    hover: {
      activeTooltip: null,
    },
    events: {
      status: "idle",
      activePopup: null,
    },
  };
}

function ensureLayerScaffold(stage, sharedLayers) {
  const existingLayers = new Set(
    Array.from(stage.querySelectorAll("[data-race-viz-layer]"), (node) => node.dataset.raceVizLayer),
  );

  for (const layerName of sharedLayers) {
    if (existingLayers.has(layerName)) {
      continue;
    }

    const layer = document.createElement("div");
    layer.className = `race-viz-layer race-viz-layer-${layerName}`;
    layer.dataset.raceVizLayer = layerName;
    stage.append(layer);
  }
}

function applyModeToStage(stage, state) {
  const { config } = state;
  stage.dataset.raceVizStage = "ready";
  stage.dataset.raceVizEngine = config.engine;
  stage.dataset.raceVizMode = config.mode;
  stage.dataset.raceVizActiveLayers = config.activeLayers.join(" ");

  for (const layer of stage.querySelectorAll("[data-race-viz-layer]")) {
    const enabled = config.activeLayers.includes(layer.dataset.raceVizLayer);
    layer.hidden = !enabled;
    layer.dataset.raceVizLayerEnabled = String(enabled);
  }
}

function getMapCanvas(stage) {
  const mapLayer = stage.querySelector('[data-race-viz-layer="map"]');
  if (!mapLayer) {
    return null;
  }

  let canvas = mapLayer.querySelector("[data-race-viz-map-canvas]");
  if (canvas) {
    return canvas;
  }

  canvas = document.createElement("div");
  canvas.className = "race-viz-map-canvas";
  canvas.dataset.raceVizMapCanvas = "";
  canvas.setAttribute("aria-hidden", "true");
  mapLayer.append(canvas);
  return canvas;
}

function getCourseLayer(stage) {
  return stage.querySelector('[data-race-viz-layer="course"]');
}

function getBoatLegend(root) {
  return root.querySelector("[data-race-viz-boat-legend]");
}

function getReplayControls(root) {
  return {
    panel: root.querySelector("[data-race-viz-controls]"),
    playToggle: root.querySelector("[data-race-viz-play-toggle]"),
    reset: root.querySelector("[data-race-viz-replay-reset]"),
    speedSelect: root.querySelector("[data-race-viz-replay-speed-select]"),
    timeline: root.querySelector("[data-race-viz-replay-timeline]"),
    currentLabel: root.querySelector("[data-race-viz-replay-current-label]"),
    startLabel: root.querySelector("[data-race-viz-replay-start-label]"),
    endLabel: root.querySelector("[data-race-viz-replay-end-label]"),
  };
}

function renderMapFallback(stage, message) {
  const canvas = getMapCanvas(stage);
  if (!canvas) {
    return;
  }

  const existing = canvas.querySelector("[data-race-viz-map-fallback]");
  if (!message) {
    existing?.remove();
    return;
  }

  canvas.replaceChildren();
  const fallback = document.createElement("div");
  fallback.className = "race-viz-map-fallback";
  fallback.dataset.raceVizMapFallback = "";
  fallback.textContent = message;
  canvas.append(fallback);
}

function renderCourseFallback(stage, message) {
  const courseLayer = getCourseLayer(stage);
  if (!courseLayer) {
    return;
  }

  let fallback = courseLayer.querySelector("[data-race-viz-course-fallback]");
  if (!message) {
    fallback?.remove();
    return;
  }

  if (!fallback) {
    fallback = document.createElement("div");
    fallback.className = "race-viz-course-fallback";
    fallback.dataset.raceVizCourseFallback = "";
    courseLayer.append(fallback);
  }

  fallback.textContent = message;
}

function setMapState(root, stage, state, status, message = "") {
  state.map.status = status;
  root.dataset.raceVizMapState = status;
  stage.dataset.raceVizMapState = status;
  renderMapFallback(stage, message);
}

function setCourseState(root, stage, state, status, message = "") {
  state.course.status = status;
  root.dataset.raceVizCourseState = status;
  stage.dataset.raceVizCourseState = status;
  renderCourseFallback(stage, message);
}

function renderBoatsFallback(stage, message) {
  const courseLayer = getCourseLayer(stage);
  if (!courseLayer) {
    return;
  }

  let fallback = courseLayer.querySelector("[data-race-viz-boats-fallback]");
  if (!message) {
    fallback?.remove();
    return;
  }

  if (!fallback) {
    fallback = document.createElement("div");
    fallback.className = "race-viz-boats-fallback";
    fallback.dataset.raceVizBoatsFallback = "";
    courseLayer.append(fallback);
  }

  fallback.textContent = message;
}

function setBoatsState(root, stage, state, status, message = "") {
  state.boats.status = status;
  root.dataset.raceVizBoatsState = status;
  stage.dataset.raceVizBoatsState = status;
  renderBoatsFallback(stage, message);
}

function setEventsState(root, state, status) {
  state.events.status = status;
  root.dataset.raceVizEventsState = status;
}

function setReplayClockState(root, state, status) {
  state.replay.status = status;
  root.dataset.raceVizReplayState = status;
}

function syncReplayClockDataset(root, replay) {
  root.dataset.raceVizReplayTime = String(replay.currentTimeMs ?? 0);
  root.dataset.raceVizReplaySpeed = String(replay.speed ?? 1);
  root.dataset.raceVizReplayPlaying = String(Boolean(replay.playing));

  if (!replay.startTime || !replay.endTime) {
    delete root.dataset.raceVizReplayStart;
    delete root.dataset.raceVizReplayEnd;
    delete root.dataset.raceVizReplayDuration;
    return;
  }

  root.dataset.raceVizReplayStart = replay.startTime;
  root.dataset.raceVizReplayEnd = replay.endTime;
  root.dataset.raceVizReplayDuration = String(replay.durationMs ?? 0);
}

function formatElapsedLabel(elapsedMs) {
  const ms = elapsedMs > 0 ? elapsedMs : 0;
  const totalSeconds = Math.floor(ms / 1000);
  const hours = Math.floor(totalSeconds / 3600);
  const minutes = Math.floor((totalSeconds % 3600) / 60);
  const seconds = totalSeconds % 60;
  return `+${String(hours).padStart(2, "0")}:${String(minutes).padStart(2, "0")}:${String(seconds).padStart(2, "0")}`;
}

function syncReplayControls(root, state) {
  const controls = getReplayControls(root);
  if (!controls.panel) {
    return;
  }

  const ready = state.replay.status === "ready" && state.replay.timeline;
  const disabled = !ready;

  controls.playToggle?.toggleAttribute("disabled", disabled);
  controls.reset?.toggleAttribute("disabled", disabled);
  controls.speedSelect?.toggleAttribute("disabled", disabled);
  controls.timeline?.toggleAttribute("disabled", disabled);

  if (controls.playToggle) {
    controls.playToggle.textContent = state.replay.playing ? "Pause" : "Play";
  }

  if (controls.speedSelect) {
    controls.speedSelect.value = String(state.replay.speed ?? 1);
  }

  if (controls.currentLabel) {
    controls.currentLabel.textContent = formatElapsedLabel(
      (state.replay.currentTimeMs ?? 0) - (state.replay.startTimeMs ?? 0),
    );
  }

  if (controls.startLabel) {
    controls.startLabel.textContent = formatElapsedLabel(0);
  }

  if (controls.endLabel) {
    controls.endLabel.textContent = formatElapsedLabel(state.replay.durationMs ?? 0);
  }

  if (controls.timeline) {
    controls.timeline.min = "0";
    controls.timeline.max = String(state.replay.durationMs ?? 0);
    controls.timeline.step = "1000";
    controls.timeline.value = String(
      Math.max(0, (state.replay.currentTimeMs ?? 0) - (state.replay.startTimeMs ?? 0)),
    );
  }
}

function renderBoatLegend(root, boats) {
  const legend = getBoatLegend(root);
  if (!legend) {
    return;
  }

  legend.replaceChildren();

  for (const boat of boats) {
    const boatId = boat.id ?? "";
    const boatName = boat.name || boatId || "Boat";

    const item = document.createElement("li");
    item.className = "race-viz-boat-legend-item";
    item.dataset.raceVizBoatLegendItem = boatId;
    item.dataset.raceVizBoatColor = boat.color ?? "";
    item.dataset.raceVizBoatRole = boat.isSelf ? "self" : "competitor";
    item.dataset.raceVizBoatHidden = "false";

    const swatch = document.createElement("span");
    swatch.className = "race-viz-boat-legend-swatch";
    swatch.dataset.raceVizBoatLegendSwatch = "";
    swatch.style.setProperty("--race-viz-boat-color", boat.color ?? "#ffffff");
    swatch.setAttribute("aria-hidden", "true");

    const label = document.createElement("span");
    label.className = "race-viz-boat-legend-label";
    label.textContent = boatName;

    const toggle = document.createElement("button");
    toggle.type = "button";
    toggle.className = "race-viz-boat-toggle";
    toggle.dataset.raceVizBoatToggle = boatId;
    toggle.setAttribute("aria-label", `Toggle ${boatName} visibility`);
    toggle.setAttribute("aria-pressed", "true");

    item.append(swatch, label, toggle);
    legend.append(item);
  }
}

function syncBoatLegendVisibility(root, state) {
  const legend = getBoatLegend(root);
  if (!legend) {
    return;
  }

  for (const item of legend.querySelectorAll("[data-race-viz-boat-legend-item]")) {
    const boatId = item.dataset.raceVizBoatLegendItem;
    const hidden = state.visibility.hiddenBoatIds.has(boatId);
    item.dataset.raceVizBoatHidden = String(hidden);
    const toggle = item.querySelector("[data-race-viz-boat-toggle]");
    if (toggle) {
      toggle.setAttribute("aria-pressed", String(!hidden));
    }
  }
}

function attachBoatLegendToggles(root, state) {
  const legend = getBoatLegend(root);
  if (!legend) {
    return;
  }

  legend.addEventListener("click", (event) => {
    const toggle = event.target.closest("[data-race-viz-boat-toggle]");
    if (!toggle) {
      return;
    }

    const boatId = toggle.dataset.raceVizBoatToggle;
    if (!boatId) {
      return;
    }

    if (state.visibility.hiddenBoatIds.has(boatId)) {
      state.visibility.hiddenBoatIds.delete(boatId);
    } else {
      state.visibility.hiddenBoatIds.add(boatId);
    }

    syncBoatLegendVisibility(root, state);

    if (state.map.instance && state.map.status === "ready") {
      applyBoatVisibilityToLayers(state.map.instance, state);
    }
  });
}

function applyBoatVisibilityToLayers(map, state) {
  const hiddenIds = state.visibility.hiddenBoatIds;

  if (!state.replay.started) {
    const tracksLayerID = state.config.tracks.layerID;
    const selfOnlyFilter = ["boolean", ["get", "isSelf"], false];

    let filter;
    if (hiddenIds.size === 0) {
      filter = selfOnlyFilter;
    } else {
      filter = [
        "all",
        selfOnlyFilter,
        ["!", ["in", ["get", "id"], ["literal", Array.from(hiddenIds)]]],
      ];
    }

    if (map.getLayer(tracksLayerID)) {
      map.setFilter(tracksLayerID, filter);
    }
    if (map.getLayer(`${tracksLayerID}-casing`)) {
      map.setFilter(`${tracksLayerID}-casing`, filter);
    }
  } else {
    renderReplayFrame(map, state);
  }
}

function initializeMap(root, stage, state) {
  return initializeMapWithFallback(root, stage, state);
}

function normalizeEndpoint(value) {
  return (value ?? "").replace(/\/+$/, "");
}

async function fetchMapStyleDefinition(styleURL) {
  const response = await window.fetch(styleURL, {
    headers: {
      Accept: "application/json",
    },
  });

  if (!response.ok) {
    throw new Error(`Map style request failed with status ${response.status}.`);
  }

  return response.json();
}

function replaceTileEndpointInStyle(style, currentEndpoint, nextEndpoint) {
  const sourceEndpoint = normalizeEndpoint(currentEndpoint);
  const targetEndpoint = normalizeEndpoint(nextEndpoint);
  if (!sourceEndpoint || !targetEndpoint || sourceEndpoint === targetEndpoint) {
    return style;
  }

  const nextStyle = JSON.parse(JSON.stringify(style));
  for (const source of Object.values(nextStyle.sources ?? {})) {
    if (!source || !Array.isArray(source.tiles)) {
      continue;
    }

    source.tiles = source.tiles.map((tileURL) =>
      typeof tileURL === "string" ? tileURL.replace(sourceEndpoint, targetEndpoint) : tileURL,
    );
  }

  return nextStyle;
}

async function buildMapStyleVariants(config) {
  const styleDefinition = await fetchMapStyleDefinition(config.styleURL);
  const primaryEndpoint = normalizeEndpoint(config.tileEndpoint);
  const fallbackEndpoint = normalizeEndpoint(config.fallbackTileEndpoint);
  const variants = [
    {
      kind: "primary",
      style: styleDefinition,
      tileEndpoint: primaryEndpoint,
    },
  ];

  if (fallbackEndpoint && fallbackEndpoint !== primaryEndpoint) {
    variants.push({
      kind: "fallback",
      style: replaceTileEndpointInStyle(styleDefinition, primaryEndpoint, fallbackEndpoint),
      tileEndpoint: fallbackEndpoint,
    });
  }

  return variants;
}

function createMapInstance(root, stage, state, canvas, variant) {
  if (!canvas) {
    setMapState(root, stage, state, "error", "Race map container is missing.");
    return Promise.reject(new Error("Race map container is missing."));
  }

  if (!state.config.map.styleURL) {
    setMapState(root, stage, state, "error", "Race map style URL is missing.");
    return Promise.reject(new Error("Race map style URL is missing."));
  }

  if (!window.maplibregl?.Map) {
    setMapState(root, stage, state, "error", "Map runtime failed to load.");
    return Promise.reject(new Error("Map runtime failed to load."));
  }

  return new Promise((resolve, reject) => {
    const mapOptions = {
      container: canvas,
      style: variant.style,
      attributionControl: false,
      maxBounds: state.config.map.maxBounds ?? variant.style.bounds ?? null,
    };
    if (state.config.map.minZoom != null) {
      mapOptions.minZoom = state.config.map.minZoom;
    }
    if (state.config.map.maxZoom != null) {
      mapOptions.maxZoom = state.config.map.maxZoom;
    }
    const map = new window.maplibregl.Map(mapOptions);

    let settled = false;

    function fail(message) {
      if (settled) {
        return;
      }

      settled = true;
      setMapState(root, stage, state, "error", message);
      root.dataset.raceVizMapVariant = variant.kind;
      reject(new Error(message));
    }

    map.addControl(new window.maplibregl.NavigationControl(), "top-right");
    map.once("load", () => {
      if (settled) {
        return;
      }

      settled = true;
      setMapState(root, stage, state, "ready");
      root.dataset.raceVizMapVariant = variant.kind;
      root.dataset.raceVizMapActiveTileEndpoint = variant.tileEndpoint;
      map.resize();
      resolve(map);
    });
    map.on("error", (event) => {
      if (state.map.status !== "ready") {
        const errorMessage = event?.error?.message ?? "Race map failed to initialize.";
        fail(errorMessage);
      }
    });

    state.map.instance = map;
  });
}

async function initializeMapWithFallback(root, stage, state) {
  const canvas = getMapCanvas(stage);
  if (!canvas) {
    setMapState(root, stage, state, "error", "Race map container is missing.");
    throw new Error("Race map container is missing.");
  }

  if (!state.config.map.styleURL) {
    setMapState(root, stage, state, "error", "Race map style URL is missing.");
    throw new Error("Race map style URL is missing.");
  }

  if (!window.maplibregl?.Map) {
    setMapState(root, stage, state, "error", "Map runtime failed to load.");
    throw new Error("Map runtime failed to load.");
  }

  const variants = await buildMapStyleVariants(state.config.map);
  let lastError = null;

  for (const [index, variant] of variants.entries()) {
    setMapState(root, stage, state, "booting");
    try {
      return await createMapInstance(root, stage, state, canvas, variant);
    } catch (error) {
      lastError = error;
      state.map.instance?.remove();
      state.map.instance = null;
      if (index === variants.length - 1) {
        break;
      }
    }
  }

  const previewCommand = state.config.map.previewCommand;
  const message = previewCommand
    ? `Race map failed to initialize. Start the local tile server with "${previewCommand}" or verify the fallback tile host.`
    : "Race map failed to initialize.";
  setMapState(root, stage, state, "error", message);
  throw lastError ?? new Error(message);
}

async function fetchJSON(url) {
  if (!url) {
    throw new Error("JSON URL is missing.");
  }

  if (typeof window.fetch !== "function") {
    throw new Error("Fetch runtime failed to load.");
  }

  const response = await window.fetch(url, {
    headers: {
      Accept: "application/json",
    },
  });

  if (!response.ok) {
    throw new Error(`JSON request failed with status ${response.status}.`);
  }

  return response.json();
}

function buildCourseRouteCoordinates(course) {
  const coordinates = [];

  for (const [index, element] of course.elements.entries()) {
    coordinates.push([element.lon, element.lat]);

    if (index === course.elements.length - 1) {
      continue;
    }

    for (const point of element.controlPointsToNext ?? []) {
      coordinates.push([point.lon, point.lat]);
    }
  }

  return coordinates;
}

function buildCourseFeatures(course) {
  const features = [];
  const routeCoordinates = buildCourseRouteCoordinates(course);

  if (routeCoordinates.length >= 2) {
    features.push({
      type: "Feature",
      geometry: {
        type: "LineString",
        coordinates: routeCoordinates,
      },
      properties: {
        featureType: "course-route",
        courseID: course.id ?? "",
      },
    });
  }

  for (const [index, element] of course.elements.entries()) {
    features.push({
      type: "Feature",
      geometry: {
        type: "Point",
        coordinates: [element.lon, element.lat],
      },
      properties: {
        featureType: "course-element",
        index,
        id: element.id ?? "",
        type: element.type ?? "mark",
        name: element.name ?? "",
        rounding: element.rounding ?? "none",
      },
    });
  }

  return {
    type: "FeatureCollection",
    features,
  };
}

function buildBoatTrackFeatures(payload) {
  return {
    type: "FeatureCollection",
    features: (payload.boats ?? [])
      .filter((boat) => Array.isArray(boat.track) && boat.track.length >= 2)
      .map((boat) => ({
        type: "Feature",
        geometry: {
          type: "LineString",
          coordinates: boat.track.map((point) => [point.lon, point.lat]),
        },
        properties: {
          featureType: "boat-track",
          id: boat.id ?? "",
          name: boat.name ?? "",
          color: boat.color ?? "#ffffff",
          boatType: boat.boatType ?? "",
          source: boat.source ?? "",
          isSelf: Boolean(boat.isSelf),
          pointCount: boat.track.length,
        },
      })),
  };
}

function normalizeTrackPoint(point) {
  const timestampMs = Date.parse(point.time);
  if (Number.isNaN(timestampMs)) {
    throw new Error(`Invalid track timestamp ${JSON.stringify(point.time)}.`);
  }

  return {
    time: point.time,
    timestampMs,
    lat: point.lat,
    lon: point.lon,
  };
}

function normalizeBoatReplayTrack(boat) {
  const track = (boat.track ?? []).map(normalizeTrackPoint);
  if (track.length === 0) {
    throw new Error(`Boat ${JSON.stringify(boat.id ?? "")} is missing replay track points.`);
  }

  for (let index = 1; index < track.length; index += 1) {
    if (track[index].timestampMs <= track[index - 1].timestampMs) {
      throw new Error(`Boat ${JSON.stringify(boat.id ?? "")} replay times must be strictly increasing.`);
    }
  }

  return {
    ...boat,
    track,
    startTime: track[0].time,
    endTime: track[track.length - 1].time,
    startTimeMs: track[0].timestampMs,
    endTimeMs: track[track.length - 1].timestampMs,
  };
}

function buildReplayTimeline(payload) {
  const boats = (payload.boats ?? []).map(normalizeBoatReplayTrack);
  if (boats.length === 0) {
    throw new Error("Replay payload must include at least one boat.");
  }

  let startBoat = boats[0];
  let endBoat = boats[0];
  for (const boat of boats.slice(1)) {
    if (boat.startTimeMs < startBoat.startTimeMs) {
      startBoat = boat;
    }
    if (boat.endTimeMs > endBoat.endTimeMs) {
      endBoat = boat;
    }
  }

  return {
    boats,
    startTime: startBoat.startTime,
    endTime: endBoat.endTime,
    startTimeMs: startBoat.startTimeMs,
    endTimeMs: endBoat.endTimeMs,
    durationMs: endBoat.endTimeMs - startBoat.startTimeMs,
  };
}

function interpolateBoatPosition(boat, timeMs) {
  const { track } = boat;

  if (timeMs <= boat.startTimeMs) {
    const point = track[0];
    return {
      id: boat.id ?? "",
      lat: point.lat,
      lon: point.lon,
      timeMs: boat.startTimeMs,
      segmentStartTimeMs: boat.startTimeMs,
      segmentEndTimeMs: boat.startTimeMs,
      progress: 0,
    };
  }

  if (timeMs >= boat.endTimeMs) {
    const point = track[track.length - 1];
    return {
      id: boat.id ?? "",
      lat: point.lat,
      lon: point.lon,
      timeMs: boat.endTimeMs,
      segmentStartTimeMs: boat.endTimeMs,
      segmentEndTimeMs: boat.endTimeMs,
      progress: 1,
    };
  }

  for (let index = 1; index < track.length; index += 1) {
    const previous = track[index - 1];
    const next = track[index];

    if (timeMs > next.timestampMs) {
      continue;
    }

    const segmentDuration = next.timestampMs - previous.timestampMs;
    const progress = segmentDuration === 0 ? 0 : (timeMs - previous.timestampMs) / segmentDuration;

    return {
      id: boat.id ?? "",
      lat: previous.lat + ((next.lat - previous.lat) * progress),
      lon: previous.lon + ((next.lon - previous.lon) * progress),
      timeMs,
      segmentStartTimeMs: previous.timestampMs,
      segmentEndTimeMs: next.timestampMs,
      progress,
    };
  }

  const point = track[track.length - 1];
  return {
    id: boat.id ?? "",
    lat: point.lat,
    lon: point.lon,
    timeMs: boat.endTimeMs,
    segmentStartTimeMs: boat.endTimeMs,
    segmentEndTimeMs: boat.endTimeMs,
    progress: 1,
  };
}

// Inverse of interpolateBoatPosition: given a map position, find the closest
// point on the boat's track and return the interpolated timestamp at that point.
function interpolateTimeFromPosition(boat, lngLat) {
  const track = boat.track;
  if (track.length === 0) return boat.startTimeMs;
  if (track.length === 1) return track[0].timestampMs;

  const curLon = lngLat.lng;
  const curLat = lngLat.lat;

  let bestTimeMs = track[0].timestampMs;
  let bestDistSq = Infinity;

  for (let i = 1; i < track.length; i++) {
    const p0 = track[i - 1];
    const p1 = track[i];

    const dx = p1.lon - p0.lon;
    const dy = p1.lat - p0.lat;
    const segLenSq = dx * dx + dy * dy;

    let t = 0;
    if (segLenSq > 0) {
      t = ((curLon - p0.lon) * dx + (curLat - p0.lat) * dy) / segLenSq;
      t = Math.max(0, Math.min(1, t));
    }

    const closestLon = p0.lon + t * dx;
    const closestLat = p0.lat + t * dy;
    const distSq = (curLon - closestLon) ** 2 + (curLat - closestLat) ** 2;

    if (distSq < bestDistSq) {
      bestDistSq = distSq;
      bestTimeMs = p0.timestampMs + t * (p1.timestampMs - p0.timestampMs);
    }
  }

  return bestTimeMs;
}

function buildReplaySnapshot(timeline, requestedTimeMs) {
  const timeMs = Math.min(
    Math.max(requestedTimeMs, timeline.startTimeMs),
    timeline.endTimeMs,
  );

  return {
    timeMs,
    boats: timeline.boats.map((boat) => ({
      id: boat.id ?? "",
      name: boat.name ?? boat.id ?? "Boat",
      color: boat.color ?? "#ffffff",
      isSelf: Boolean(boat.isSelf),
      position: interpolateBoatPosition(boat, timeMs),
    })),
  };
}

function stopReplayPlayback(state) {
  if (state.replay.animationFrameID) {
    window.cancelAnimationFrame(state.replay.animationFrameID);
  }

  state.replay.animationFrameID = 0;
  state.replay.lastFrameMs = 0;
  state.replay.playing = false;
}

function setReplayTime(root, state, requestedTimeMs) {
  if (!state.replay.timeline) {
    return;
  }

  state.replay.currentTimeMs = Math.min(
    Math.max(requestedTimeMs, state.replay.startTimeMs),
    state.replay.endTimeMs,
  );
  state.replay.snapshot = buildReplaySnapshot(state.replay.timeline, state.replay.currentTimeMs);
  syncReplayClockDataset(root, state.replay);
  syncReplayControls(root, state);

  if (state.replay.started && state.map.instance && state.map.status === "ready") {
    renderReplayFrame(state.map.instance, state);
  }
}

function resetReplay(root, state) {
  stopReplayPlayback(state);
  state.replay.started = false;
  setReplayTime(root, state, state.replay.startTimeMs);
  if (state.map.instance && state.map.status === "ready") {
    enterPrePlayMode(state.map.instance, state);
  }
}

function startReplayPlayback(root, state) {
  if (!state.replay.timeline || state.replay.playing) {
    return;
  }

  if (state.replay.currentTimeMs >= state.replay.endTimeMs) {
    setReplayTime(root, state, state.replay.startTimeMs);
  }

  if (!state.replay.started) {
    state.replay.started = true;
    if (state.map.instance && state.map.status === "ready") {
      enterPlayingMode(state.map.instance, state);
    }
  }

  state.replay.playing = true;
  state.replay.lastFrameMs = 0;
  syncReplayClockDataset(root, state.replay);
  syncReplayControls(root, state);

  const step = (frameMs) => {
    if (!state.replay.playing) {
      return;
    }

    if (!state.replay.lastFrameMs) {
      state.replay.lastFrameMs = frameMs;
    }

    const elapsedMs = frameMs - state.replay.lastFrameMs;
    state.replay.lastFrameMs = frameMs;

    const nextTimeMs = state.replay.currentTimeMs + (elapsedMs * state.replay.speed);
    if (nextTimeMs >= state.replay.endTimeMs) {
      stopReplayPlayback(state);
      setReplayTime(root, state, state.replay.endTimeMs);
      return;
    }

    setReplayTime(root, state, nextTimeMs);
    state.replay.animationFrameID = window.requestAnimationFrame(step);
  };

  state.replay.animationFrameID = window.requestAnimationFrame(step);
}

function attachReplayControls(root, state) {
  const controls = getReplayControls(root);
  if (!controls.panel) {
    return;
  }

  controls.playToggle?.addEventListener("click", () => {
    if (state.replay.playing) {
      stopReplayPlayback(state);
      syncReplayClockDataset(root, state.replay);
      syncReplayControls(root, state);
      return;
    }

    startReplayPlayback(root, state);
  });

  controls.reset?.addEventListener("click", () => {
    resetReplay(root, state);
  });

  controls.speedSelect?.addEventListener("change", (event) => {
    const nextSpeed = Number.parseFloat(event.target.value);
    if (Number.isFinite(nextSpeed) && nextSpeed > 0) {
      state.replay.speed = nextSpeed;
      syncReplayClockDataset(root, state.replay);
      syncReplayControls(root, state);
    }
  });

  controls.timeline?.addEventListener("input", (event) => {
    const nextOffsetMs = Number.parseInt(event.target.value, 10);
    if (!Number.isFinite(nextOffsetMs)) {
      return;
    }

    if (!state.replay.started) {
      state.replay.started = true;
      if (state.map.instance && state.map.status === "ready") {
        enterPlayingMode(state.map.instance, state);
      }
    }

    setReplayTime(root, state, state.replay.startTimeMs + nextOffsetMs);
  });

  syncReplayControls(root, state);
}

function upsertCourseSource(map, state, data) {
  const source = map.getSource(state.config.course.sourceID);
  if (source) {
    source.setData(data);
    return;
  }

  map.addSource(state.config.course.sourceID, {
    type: "geojson",
    data,
  });
}

function ensureCourseLayer(map, layerID, layerConfig) {
  if (map.getLayer(layerID)) {
    return;
  }

  map.addLayer({
    id: layerID,
    ...layerConfig,
  });
}

function renderCourseLayers(map, state) {
  const source = state.config.course.sourceID;
  const routeLayerID = state.config.course.routeLayerID;
  const marksLayerID = state.config.course.marksLayerID;
  const roundingLayerID = state.config.course.roundingLayerID;
  const startFinishLayerID = state.config.course.startFinishLayerID;
  const labelsLayerID = state.config.course.labelsLayerID;
  const courseStyle = getCourseStyle(state.config);

  ensureCourseLayer(map, `${routeLayerID}-casing`, {
    type: "line",
    source,
    filter: ["==", ["geometry-type"], "LineString"],
    layout: {
      "line-cap": "round",
      "line-join": "round",
    },
    paint: {
      "line-color": courseStyle.routeCasingColor,
      "line-width": courseStyle.routeCasingWidth,
      "line-opacity": 0.9,
    },
  });

  ensureCourseLayer(map, `${routeLayerID}-glow`, {
    type: "line",
    source,
    filter: ["==", ["geometry-type"], "LineString"],
    layout: {
      "line-cap": "round",
      "line-join": "round",
    },
    paint: {
      "line-color": courseStyle.routeGlowColor,
      "line-width": courseStyle.routeGlowWidth,
      "line-opacity": 1,
      "line-blur": 0.45,
    },
  });

  ensureCourseLayer(map, routeLayerID, {
    type: "line",
    source,
    filter: ["==", ["geometry-type"], "LineString"],
    layout: {
      "line-cap": "round",
      "line-join": "round",
    },
    paint: {
      "line-color": courseStyle.routeColor,
      "line-width": courseStyle.routeWidth,
      "line-opacity": 0.96,
      "line-dasharray": courseStyle.routeDasharray,
    },
  });

  ensureCourseLayer(map, startFinishLayerID, {
    type: "circle",
    source,
    filter: [
      "all",
      ["==", ["geometry-type"], "Point"],
      ["match", ["get", "type"], ["start_line", "finish_line"], true, false],
    ],
    paint: {
      "circle-color": [
        "match",
        ["get", "type"],
        "start_line",
        courseStyle.startColor,
        courseStyle.finishColor,
      ],
      "circle-radius": 9,
      "circle-stroke-width": 3,
      "circle-stroke-color": courseStyle.markStrokeColor,
    },
  });

  ensureCourseLayer(map, roundingLayerID, {
    type: "circle",
    source,
    filter: [
      "all",
      ["==", ["geometry-type"], "Point"],
      ["==", ["get", "type"], "mark"],
      ["match", ["get", "rounding"], ["port", "starboard"], true, false],
    ],
    paint: {
      "circle-color": [
        "match",
        ["get", "rounding"],
        "port", courseStyle.roundingPortColor,
        "starboard", courseStyle.roundingStarboardColor,
        "rgba(0,0,0,0)",
      ],
      "circle-radius": 11,
      "circle-stroke-width": 0,
    },
  });

  ensureCourseLayer(map, marksLayerID, {
    type: "circle",
    source,
    filter: [
      "all",
      ["==", ["geometry-type"], "Point"],
      ["==", ["get", "type"], "mark"],
    ],
    paint: {
      "circle-color": courseStyle.markColor,
      "circle-radius": 7,
      "circle-stroke-width": 2,
      "circle-stroke-color": courseStyle.markStrokeColor,
    },
  });

  ensureCourseLayer(map, labelsLayerID, {
    type: "symbol",
    source,
    filter: [
      "all",
      ["==", ["geometry-type"], "Point"],
      ["!=", ["get", "name"], ""],
      ["match", ["get", "type"], ["mark", "start_line", "finish_line"], true, false],
    ],
    layout: {
      "text-field": ["upcase", ["get", "name"]],
      "text-size": 11,
      "text-letter-spacing": 0.18,
      "text-offset": [0, 1.3],
      "text-anchor": "top",
      "text-transform": "uppercase",
      "text-allow-overlap": true,
    },
    paint: {
      "text-color": courseStyle.labelColor,
      "text-halo-color": courseStyle.labelHaloColor,
      "text-halo-width": 1.25,
      "text-halo-blur": 0.2,
    },
  });
}

function upsertTracksSource(map, state, data) {
  const source = map.getSource(state.config.tracks.sourceID);
  if (source) {
    source.setData(data);
    return;
  }

  map.addSource(state.config.tracks.sourceID, {
    type: "geojson",
    data,
  });
}

function renderTrackLayers(map, state) {
  const source = state.config.tracks.sourceID;
  const layerID = state.config.tracks.layerID;

  ensureCourseLayer(map, `${layerID}-casing`, {
    type: "line",
    source,
    layout: {
      "line-cap": "round",
      "line-join": "round",
    },
    paint: {
      "line-color": "rgba(5, 14, 24, 0.92)",
      "line-width": [
        "case",
        ["boolean", ["get", "isSelf"], false],
        6,
        5,
      ],
      "line-opacity": 0.88,
    },
  });

  ensureCourseLayer(map, layerID, {
    type: "line",
    source,
    layout: {
      "line-cap": "round",
      "line-join": "round",
    },
    paint: {
      "line-color": ["coalesce", ["get", "color"], "#ffffff"],
      "line-width": [
        "case",
        ["boolean", ["get", "isSelf"], false],
        3.25,
        2.5,
      ],
      "line-opacity": [
        "case",
        ["boolean", ["get", "isSelf"], false],
        0.95,
        0.78,
      ],
    },
  });
}

function fitCourseBounds(map, courseFeatures, config) {
  const coordinates = [];

  for (const feature of courseFeatures.features) {
    if (feature.geometry.type === "Point") {
      coordinates.push(feature.geometry.coordinates);
      continue;
    }

    for (const coordinate of feature.geometry.coordinates) {
      coordinates.push(coordinate);
    }
  }

  if (coordinates.length === 0) {
    return;
  }

  if (coordinates.length === 1) {
    map.easeTo({
      center: coordinates[0],
      zoom: Math.max(map.getZoom(), 11),
      duration: 0,
    });
    return;
  }

  const bounds = coordinates.reduce(
    (accumulator, coordinate) => accumulator.extend(coordinate),
    new window.maplibregl.LngLatBounds(coordinates[0], coordinates[0]),
  );

  const fitMaxZoom = config?.fitMaxZoom ?? DEFAULT_MAP_FIT_MAX_ZOOM;
  map.fitBounds(bounds, {
    padding: DEFAULT_MAP_FIT_PADDING,
    duration: 0,
    maxZoom: fitMaxZoom,
  });

}

function emptyFeatureCollection() {
  return { type: "FeatureCollection", features: [] };
}

function buildTrackTailCoordinates(boat, timeMs) {
  const coords = [];

  for (const point of boat.track) {
    coords.push([point.lon, point.lat]);
    if (point.timestampMs >= timeMs) {
      break;
    }
  }

  const pos = interpolateBoatPosition(boat, timeMs);
  const last = coords[coords.length - 1];
  if (!last || last[0] !== pos.lon || last[1] !== pos.lat) {
    coords.push([pos.lon, pos.lat]);
  }

  return coords;
}

function buildReplayTailFeatures(timeline, timeMs, hiddenBoatIds = null) {
  const features = [];

  for (const boat of timeline.boats) {
    if (hiddenBoatIds !== null && hiddenBoatIds.has(boat.id ?? "")) {
      continue;
    }

    const coords = buildTrackTailCoordinates(boat, timeMs);
    if (coords.length < 2) {
      continue;
    }

    features.push({
      type: "Feature",
      geometry: { type: "LineString", coordinates: coords },
      properties: {
        featureType: "replay-tail",
        id: boat.id ?? "",
        name: boat.name ?? "",
        color: boat.color ?? "#ffffff",
        isSelf: Boolean(boat.isSelf),
      },
    });
  }

  return { type: "FeatureCollection", features };
}

function buildBoatMarkerFeatures(snapshot, hiddenBoatIds = null) {
  return {
    type: "FeatureCollection",
    features: snapshot.boats
      .filter((boat) => hiddenBoatIds === null || !hiddenBoatIds.has(boat.id ?? ""))
      .map((boat) => ({
        type: "Feature",
        geometry: {
          type: "Point",
          coordinates: [boat.position.lon, boat.position.lat],
        },
        properties: {
          featureType: "boat-marker",
          id: boat.id ?? "",
          name: boat.name ?? "",
          color: boat.color ?? "#ffffff",
          isSelf: Boolean(boat.isSelf),
        },
      })),
  };
}

function upsertReplayTailsSource(map, state, data) {
  const source = map.getSource(state.config.replayTails.sourceID);
  if (source) {
    source.setData(data);
    return;
  }

  map.addSource(state.config.replayTails.sourceID, { type: "geojson", data });
}

function upsertBoatMarkersSource(map, state, data) {
  const source = map.getSource(state.config.boatMarkers.sourceID);
  if (source) {
    source.setData(data);
    return;
  }

  map.addSource(state.config.boatMarkers.sourceID, { type: "geojson", data });
}

function renderReplayTailLayers(map, state) {
  const source = state.config.replayTails.sourceID;
  const layerID = state.config.replayTails.layerID;

  ensureCourseLayer(map, `${layerID}-casing`, {
    type: "line",
    source,
    layout: { "line-cap": "round", "line-join": "round", visibility: "none" },
    paint: {
      "line-color": "rgba(5, 14, 24, 0.92)",
      "line-width": ["case", ["boolean", ["get", "isSelf"], false], 6, 5],
      "line-opacity": 0.88,
    },
  });

  ensureCourseLayer(map, layerID, {
    type: "line",
    source,
    layout: { "line-cap": "round", "line-join": "round", visibility: "none" },
    paint: {
      "line-color": ["coalesce", ["get", "color"], "#ffffff"],
      "line-width": ["case", ["boolean", ["get", "isSelf"], false], 3.25, 2.5],
      "line-opacity": ["case", ["boolean", ["get", "isSelf"], false], 0.95, 0.78],
    },
  });
}

function renderBoatMarkerLayers(map, state) {
  const source = state.config.boatMarkers.sourceID;
  const layerID = state.config.boatMarkers.layerID;

  ensureCourseLayer(map, `${layerID}-halo`, {
    type: "circle",
    source,
    layout: { visibility: "none" },
    paint: {
      "circle-color": "rgba(5, 14, 24, 0.88)",
      "circle-radius": ["case", ["boolean", ["get", "isSelf"], false], 10, 9],
    },
  });

  ensureCourseLayer(map, layerID, {
    type: "circle",
    source,
    layout: { visibility: "none" },
    paint: {
      "circle-color": ["coalesce", ["get", "color"], "#ffffff"],
      "circle-radius": ["case", ["boolean", ["get", "isSelf"], false], 6, 5],
      "circle-stroke-width": 1.5,
      "circle-stroke-color": "rgba(255,255,255,0.5)",
    },
  });
}

function setLayerVisibility(map, layerID, visible) {
  if (!map.getLayer(layerID)) {
    return;
  }

  map.setLayoutProperty(layerID, "visibility", visible ? "visible" : "none");
}

function enterPrePlayMode(map, state) {
  const tracksLayerID = state.config.tracks.layerID;
  const replayTailsLayerID = state.config.replayTails.layerID;
  const boatMarkersLayerID = state.config.boatMarkers.layerID;
  const selfOnlyFilter = ["boolean", ["get", "isSelf"], false];

  if (map.getLayer(tracksLayerID)) {
    map.setFilter(tracksLayerID, selfOnlyFilter);
  }
  if (map.getLayer(`${tracksLayerID}-casing`)) {
    map.setFilter(`${tracksLayerID}-casing`, selfOnlyFilter);
  }

  setLayerVisibility(map, tracksLayerID, true);
  setLayerVisibility(map, `${tracksLayerID}-casing`, true);
  setLayerVisibility(map, replayTailsLayerID, false);
  setLayerVisibility(map, `${replayTailsLayerID}-casing`, false);
  setLayerVisibility(map, boatMarkersLayerID, false);
  setLayerVisibility(map, `${boatMarkersLayerID}-halo`, false);
}

function enterPlayingMode(map, state) {
  const tracksLayerID = state.config.tracks.layerID;
  const replayTailsLayerID = state.config.replayTails.layerID;
  const boatMarkersLayerID = state.config.boatMarkers.layerID;

  setLayerVisibility(map, tracksLayerID, false);
  setLayerVisibility(map, `${tracksLayerID}-casing`, false);
  setLayerVisibility(map, replayTailsLayerID, true);
  setLayerVisibility(map, `${replayTailsLayerID}-casing`, true);
  setLayerVisibility(map, boatMarkersLayerID, true);
  setLayerVisibility(map, `${boatMarkersLayerID}-halo`, true);
}

function renderReplayFrame(map, state) {
  if (!state.replay.snapshot || !state.replay.timeline) {
    return;
  }

  const hiddenBoatIds = state.visibility.hiddenBoatIds;
  const tailFeatures = buildReplayTailFeatures(state.replay.timeline, state.replay.currentTimeMs, hiddenBoatIds);
  upsertReplayTailsSource(map, state, tailFeatures);

  const markerFeatures = buildBoatMarkerFeatures(state.replay.snapshot, hiddenBoatIds);
  upsertBoatMarkersSource(map, state, markerFeatures);
}

function formatEventTime(isoTime) {
  if (!isoTime) {
    return "";
  }

  const ts = Date.parse(isoTime);
  if (Number.isNaN(ts)) {
    return isoTime;
  }

  return new Date(ts).toISOString().slice(11, 19);
}

function buildEventFeatures(payload, timeline) {
  const features = [];
  const selfBoat = timeline?.boats.find((b) => Boolean(b.isSelf)) ?? null;

  for (const event of payload.events ?? []) {
    let lat = event.lat ?? null;
    let lon = event.lon ?? null;

    if ((lat == null || lon == null) && event.time && selfBoat) {
      const timestampMs = Date.parse(event.time);
      if (!Number.isNaN(timestampMs)) {
        const pos = interpolateBoatPosition(selfBoat, timestampMs);
        lat = pos.lat;
        lon = pos.lon;
      }
    }

    if (lat == null || lon == null) {
      continue;
    }

    features.push({
      type: "Feature",
      geometry: {
        type: "Point",
        coordinates: [lon, lat],
      },
      properties: {
        featureType: "event-annotation",
        id: event.id ?? "",
        type: event.type ?? "note",
        time: event.time ?? "",
        label: event.label ?? "",
        description: event.description ?? "",
      },
    });
  }

  return { type: "FeatureCollection", features };
}

function upsertEventsSource(map, state, data) {
  const source = map.getSource(state.config.events.sourceID);
  if (source) {
    source.setData(data);
    return;
  }

  map.addSource(state.config.events.sourceID, { type: "geojson", data });
}

function renderEventLayers(map, state) {
  const source = state.config.events.sourceID;
  const layerID = state.config.events.layerID;

  ensureCourseLayer(map, `${layerID}-halo`, {
    type: "circle",
    source,
    paint: {
      "circle-color": "rgba(5, 14, 24, 0.88)",
      "circle-radius": 13,
    },
  });

  ensureCourseLayer(map, layerID, {
    type: "circle",
    source,
    paint: {
      "circle-color": "rgba(255, 200, 60, 0.95)",
      "circle-radius": 8,
      "circle-stroke-width": 2,
      "circle-stroke-color": "rgba(255, 240, 180, 0.88)",
    },
  });

  ensureCourseLayer(map, `${layerID}-label`, {
    type: "symbol",
    source,
    filter: ["!=", ["get", "label"], ""],
    layout: {
      "text-field": ["get", "label"],
      "text-size": 10.5,
      "text-letter-spacing": 0.06,
      "text-offset": [0, 1.6],
      "text-anchor": "top",
      "text-max-width": 12,
    },
    paint: {
      "text-color": "rgba(255, 230, 160, 0.96)",
      "text-halo-color": "rgba(6, 18, 28, 0.94)",
      "text-halo-width": 1.5,
      "text-halo-blur": 0.2,
    },
  });
}

function attachEventInteractions(map, state) {
  const layerID = state.config.events.layerID;

  map.on("mouseenter", layerID, () => {
    map.getCanvas().style.cursor = "pointer";
  });

  map.on("mouseleave", layerID, () => {
    map.getCanvas().style.cursor = "";
  });

  map.on("click", layerID, (event) => {
    const feature = event.features?.[0];
    if (!feature) {
      return;
    }

    const props = feature.properties;
    const coordinates = feature.geometry.coordinates.slice();
    const parts = [];

    if (props.label) {
      parts.push(`<strong class="race-viz-event-popup-label">${props.label}</strong>`);
    }
    if (props.description) {
      parts.push(`<span class="race-viz-event-popup-description">${props.description}</span>`);
    }
    if (props.time) {
      parts.push(`<time class="race-viz-event-popup-time">${formatEventTime(props.time)}</time>`);
    }
    if (props.type) {
      parts.push(`<span class="race-viz-event-type">${props.type.replace(/_/g, " ")}</span>`);
    }

    if (state.events.activePopup) {
      state.events.activePopup.remove();
      state.events.activePopup = null;
    }

    state.events.activePopup = new window.maplibregl.Popup({
      closeButton: true,
      closeOnClick: false,
      className: "race-viz-event-popup",
      maxWidth: "240px",
    })
      .setLngLat(coordinates)
      .setHTML(`<div class="race-viz-event-popup-content">${parts.join("")}</div>`)
      .addTo(map);
  });
}

function attachBoatMarkerHoverInteractions(map, state) {
  const layerID = state.config.boatMarkers.layerID;

  map.on("mouseenter", layerID, (event) => {
    map.getCanvas().style.cursor = "pointer";

    const feature = event.features?.[0];
    if (!feature) {
      return;
    }

    const props = feature.properties;
    const coordinates = feature.geometry.coordinates.slice();

    if (state.hover.activeTooltip) {
      state.hover.activeTooltip.remove();
      state.hover.activeTooltip = null;
    }

    const parts = [];
    if (props.name) {
      parts.push(`<strong class="race-viz-hover-name">${props.name}</strong>`);
    }
    const timeLabel = formatElapsedLabel(state.replay.currentTimeMs - (state.replay.startTimeMs ?? 0));
    parts.push(`<time class="race-viz-hover-time">${timeLabel}</time>`);

    state.hover.activeTooltip = new window.maplibregl.Popup({
      closeButton: false,
      closeOnClick: false,
      className: "race-viz-hover-tooltip",
      offset: 14,
    })
      .setLngLat(coordinates)
      .setHTML(`<div class="race-viz-hover-tooltip-content">${parts.join("")}</div>`)
      .addTo(map);
  });

  map.on("mouseleave", layerID, () => {
    map.getCanvas().style.cursor = "";

    if (state.hover.activeTooltip) {
      state.hover.activeTooltip.remove();
      state.hover.activeTooltip = null;
    }
  });
}

function attachTrackHoverInteractions(map, state) {
  const tracksLayerID = state.config.tracks.layerID;
  const tailsLayerID = state.config.replayTails.layerID;

  for (const layerID of [tracksLayerID, tailsLayerID]) {
    map.on("mousemove", layerID, (event) => {
      map.getCanvas().style.cursor = "crosshair";

      const feature = event.features?.[0];
      if (!feature) {
        return;
      }

      const props = feature.properties;
      const lngLat = event.lngLat;

      const parts = [];
      if (props.name) {
        parts.push(`<strong class="race-viz-hover-name">${props.name}</strong>`);
      }

      const boat = state.replay.timeline?.boats?.find((b) => b.id === props.id);
      if (boat) {
        const timeMs = interpolateTimeFromPosition(boat, lngLat);
        parts.push(`<time class="race-viz-hover-time">${formatElapsedLabel(timeMs - (state.replay.startTimeMs ?? 0))}</time>`);
      }

      const html = `<div class="race-viz-hover-tooltip-content">${parts.join("")}</div>`;

      if (state.hover.activeTooltip) {
        state.hover.activeTooltip.setLngLat(lngLat).setHTML(html);
      } else {
        state.hover.activeTooltip = new window.maplibregl.Popup({
          closeButton: false,
          closeOnClick: false,
          className: "race-viz-hover-tooltip",
          offset: 10,
        })
          .setLngLat(lngLat)
          .setHTML(html)
          .addTo(map);
      }
    });

    map.on("mouseleave", layerID, () => {
      map.getCanvas().style.cursor = "";

      if (state.hover.activeTooltip) {
        state.hover.activeTooltip.remove();
        state.hover.activeTooltip = null;
      }
    });
  }
}

async function loadEvents(root, stage, state, mapReadyPromise, boatsReadyPromise) {
  if (!state.config.activeLayers.includes("events") || !state.config.eventsURL) {
    return;
  }

  setEventsState(root, state, "loading");

  try {
    await boatsReadyPromise?.catch(() => {});

    const [payload, map] = await Promise.all([
      fetchJSON(state.config.eventsURL),
      mapReadyPromise,
    ]);

    const eventFeatures = buildEventFeatures(payload, state.replay.timeline);
    state.data.events = payload;
    upsertEventsSource(map, state, eventFeatures);
    renderEventLayers(map, state);
    attachEventInteractions(map, state);
    setEventsState(root, state, "ready");
  } catch (error) {
    setEventsState(root, state, "error");
  }
}

async function loadCourse(root, stage, state, mapReadyPromise) {
  if (!state.config.activeLayers.includes("course")) {
    return;
  }

  setCourseState(root, stage, state, "loading");

  try {
    const course = await fetchJSON(state.config.course.url);
    const courseFeatures = buildCourseFeatures(course);

    state.data.course = course;

    const map = await mapReadyPromise;
    upsertCourseSource(map, state, courseFeatures);
    renderCourseLayers(map, state);
    fitCourseBounds(map, courseFeatures, state.config);
    setCourseState(root, stage, state, "ready");
  } catch (error) {
    const message = error instanceof Error ? error.message : "Race course failed to load.";
    setCourseState(root, stage, state, "error", message);
  }
}

async function loadBoats(root, stage, state, mapReadyPromise) {
  if (!state.config.activeLayers.includes("tracks")) {
    return;
  }

  setBoatsState(root, stage, state, "loading");
  setReplayClockState(root, state, "loading");
  syncReplayControls(root, state);

  try {
    const payload = await fetchJSON(state.config.boatsURL);
    const trackFeatures = buildBoatTrackFeatures(payload);
    const timeline = buildReplayTimeline(payload);

    state.data.boats = payload;
    state.replay.timeline = timeline;
    state.replay.startTime = timeline.startTime;
    state.replay.endTime = timeline.endTime;
    state.replay.startTimeMs = timeline.startTimeMs;
    state.replay.endTimeMs = timeline.endTimeMs;
    state.replay.durationMs = timeline.durationMs;
    state.replay.currentTimeMs = timeline.startTimeMs;
    state.replay.snapshot = buildReplaySnapshot(timeline, timeline.startTimeMs);
    root.dataset.raceVizBoatCount = String(payload.boats?.length ?? 0);
    root.dataset.raceVizSelfBoatId =
      payload.boats?.find((boat) => boat.isSelf)?.id ?? "";
    syncReplayClockDataset(root, state.replay);

    renderBoatLegend(root, payload.boats ?? []);
    attachBoatLegendToggles(root, state);
    syncReplayControls(root, state);

    const map = await mapReadyPromise;
    upsertTracksSource(map, state, trackFeatures);
    renderTrackLayers(map, state);
    upsertReplayTailsSource(map, state, emptyFeatureCollection());
    renderReplayTailLayers(map, state);
    upsertBoatMarkersSource(map, state, emptyFeatureCollection());
    renderBoatMarkerLayers(map, state);
    attachBoatMarkerHoverInteractions(map, state);
    attachTrackHoverInteractions(map, state);
    enterPrePlayMode(map, state);
    setBoatsState(root, stage, state, "ready");
    setReplayClockState(root, state, "ready");
    syncReplayControls(root, state);
  } catch (error) {
    stopReplayPlayback(state);
    setBoatsState(root, stage, state, "error", "Could not load boat tracks.");
    setReplayClockState(root, state, "error");
    syncReplayClockDataset(root, state.replay);
    syncReplayControls(root, state);
  }
}

function bootRaceVizRoot(root) {
  if (root.dataset.raceVizBooted === "true") {
    return;
  }

  const stage = root.querySelector(".race-viz-stage");
  if (!stage) {
    return;
  }

  const config = createRaceVizConfig(root);
  const state = createRaceVizState(config);

  ensureLayerScaffold(stage, config.sharedLayers);
  applyModeToStage(stage, state);
  setCourseState(root, stage, state, "idle");
  setBoatsState(root, stage, state, "idle");
  setReplayClockState(root, state, "idle");
  setEventsState(root, state, "idle");
  syncReplayClockDataset(root, state.replay);
  attachReplayControls(root, state);

  const mapReadyPromise = initializeMap(root, stage, state);
  void loadCourse(root, stage, state, mapReadyPromise);
  const boatsReadyPromise = loadBoats(root, stage, state, mapReadyPromise);
  void loadEvents(root, stage, state, mapReadyPromise, boatsReadyPromise);

  const drawerToggle = root.querySelector("[data-race-viz-drawer-toggle]");
  if (drawerToggle) {
    drawerToggle.addEventListener("click", () => {
      const isOpen = root.dataset.raceVizSidebarOpen === "true";
      root.dataset.raceVizSidebarOpen = String(!isOpen);
      drawerToggle.setAttribute("aria-expanded", String(!isOpen));
    });
  }

  root.dataset.raceVizBooted = "true";
  root.dataset.raceVizState = "ready";
}

function bootAllRaceViz() {
  for (const root of document.querySelectorAll("[data-race-viz]")) {
    bootRaceVizRoot(root);
  }
}

if (document.readyState === "loading") {
  document.addEventListener("DOMContentLoaded", bootAllRaceViz, { once: true });
} else {
  bootAllRaceViz();
}
