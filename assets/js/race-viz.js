const DEFAULT_SHARED_LAYERS = ["map", "course", "tracks", "boats", "events", "controls"];
const DEFAULT_COURSE_SOURCE_ID = "race-viz-course";
const DEFAULT_COURSE_ROUTE_LAYER_ID = "race-viz-course-route";
const DEFAULT_COURSE_MARKS_LAYER_ID = "race-viz-course-marks";
const DEFAULT_COURSE_START_FINISH_LAYER_ID = "race-viz-course-start-finish";
const DEFAULT_COURSE_LABELS_LAYER_ID = "race-viz-course-labels";
const DEFAULT_COURSE_PALETTE = "signal-v1";
const DEFAULT_MAP_FIT_PADDING = 48;
const DEFAULT_MAP_FIT_MAX_ZOOM = 9.25;
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
    startColor: "rgba(255, 211, 92, 0.98)",
    finishColor: "rgba(255, 112, 146, 0.98)",
    labelColor: "rgba(220, 248, 255, 0.96)",
    labelHaloColor: "rgba(6, 18, 28, 0.94)",
  },
};

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
      styleURL: root.dataset.raceVizMapStyleUrl ?? "",
      tileEndpoint: root.dataset.raceVizMapTileEndpoint ?? "",
      tileManifestURL: root.dataset.raceVizMapTileManifestUrl ?? "",
      tileSet: root.dataset.raceVizMapTileSet ?? "",
      servingMode: root.dataset.raceVizMapServingMode ?? "",
      prototypePage: root.dataset.raceVizMapPrototypePage ?? "",
      prototypeStyle: root.dataset.raceVizMapPrototypeStyle ?? "",
    },
    course: {
      url: root.dataset.courseUrl ?? "",
      sourceID: root.dataset.raceVizCourseSource ?? DEFAULT_COURSE_SOURCE_ID,
      routeLayerID:
        root.dataset.raceVizCourseRouteLayer ?? DEFAULT_COURSE_ROUTE_LAYER_ID,
      marksLayerID:
        root.dataset.raceVizCourseMarksLayer ?? DEFAULT_COURSE_MARKS_LAYER_ID,
      startFinishLayerID:
        root.dataset.raceVizCourseStartFinishLayer ?? DEFAULT_COURSE_START_FINISH_LAYER_ID,
      labelsLayerID:
        root.dataset.raceVizCourseLabelsLayer ?? DEFAULT_COURSE_LABELS_LAYER_ID,
      palette: root.dataset.raceVizCoursePalette ?? DEFAULT_COURSE_PALETTE,
    },
    boatsURL: root.dataset.boatsUrl ?? "",
    eventsURL: root.dataset.eventsUrl ?? "",
  };
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
    map: {
      instance: null,
      status: "idle",
    },
    replay: {
      time: 0,
      playing: false,
      speed: 1,
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

function renderMapFallback(stage, message) {
  const canvas = getMapCanvas(stage);
  if (!canvas) {
    return;
  }

  canvas.replaceChildren();
  const fallback = document.createElement("div");
  fallback.className = "race-viz-map-fallback";
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
  if (message) {
    renderMapFallback(stage, message);
  }
}

function setCourseState(root, stage, state, status, message = "") {
  state.course.status = status;
  root.dataset.raceVizCourseState = status;
  stage.dataset.raceVizCourseState = status;
  renderCourseFallback(stage, message);
}

function initializeMap(root, stage, state) {
  const canvas = getMapCanvas(stage);
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

  setMapState(root, stage, state, "booting");

  return new Promise((resolve, reject) => {
    const map = new window.maplibregl.Map({
      container: canvas,
      style: state.config.map.styleURL,
      attributionControl: false,
    });

    let settled = false;

    function fail(message) {
      if (settled) {
        return;
      }

      settled = true;
      setMapState(root, stage, state, "error", message);
      reject(new Error(message));
    }

    map.addControl(new window.maplibregl.NavigationControl(), "top-right");
    map.once("load", () => {
      if (settled) {
        return;
      }

      settled = true;
      setMapState(root, stage, state, "ready");
      map.resize();
      resolve(map);
    });
    map.on("error", () => {
      if (state.map.status !== "ready") {
        fail("Race map failed to initialize.");
      }
    });

    state.map.instance = map;
  });
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

function fitCourseBounds(map, courseFeatures) {
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

  map.fitBounds(bounds, {
    padding: DEFAULT_MAP_FIT_PADDING,
    duration: 0,
    maxZoom: DEFAULT_MAP_FIT_MAX_ZOOM,
  });

  const fittedZoom = map.getZoom();
  map.setMinZoom(fittedZoom);
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
    fitCourseBounds(map, courseFeatures);
    setCourseState(root, stage, state, "ready");
  } catch (error) {
    const message = error instanceof Error ? error.message : "Race course failed to load.";
    setCourseState(root, stage, state, "error", message);
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

  const mapReadyPromise = initializeMap(root, stage, state);
  void loadCourse(root, stage, state, mapReadyPromise);

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
