#!/usr/bin/env bash
# Rebuilds tiles/mbtiles/enc_combined.mbtiles from source ENCs.
#
# Charts included (coarsest first — finer charts append on top):
#   US1WC07M (1:7M)  — full eastern Pacific base layer
#   US1HA02M (1:2M)  — Hawaii intermediate detail
#   US1WC01M (1:1M)  — West Coast detailed coverage
#   US1HA01M (1:1M)  — Hawaii detailed coverage
#
# Processing order ensures more-detailed charts paint over coarser ones
# in overlapping areas, producing a unified appearance across the corridor.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
MBTILES_DIR="$SCRIPT_DIR/mbtiles"
WORK_DIR="/tmp/enc_build_$$"
COMBINED_GPKG="$WORK_DIR/combined.gpkg"
GEOJSON_DIR="$WORK_DIR/geojson"

mkdir -p "$WORK_DIR" "$GEOJSON_DIR"
echo "Work dir: $WORK_DIR"

# Charts in processing order: coarsest scale first, finest last.
# This way more-detailed features land on top in the GeoPackage,
# winning in overlapping areas when rendered.
CHART_FILES=(
  "$HOME/Downloads/ENC_ROOT-3/US1WC07M/US1WC07M.000"
  "$HOME/Downloads/ENC_ROOT-6/US1HA02M/US1HA02M.000"
  "$HOME/Downloads/ENC_ROOT-4/US1WC01M/US1WC01M.000"
  "$HOME/Downloads/ENC_ROOT-5/US1HA01M/US1HA01M.000"
)

# ── Phase 1: validate inputs ──────────────────────────────────────────────────
echo
echo "=== Phase 1: Validating chart files ==="
for f in "${CHART_FILES[@]}"; do
  if [ -f "$f" ]; then
    echo "  OK  $f"
  else
    echo "  MISSING: $f"
    exit 1
  fi
done

# ── Phase 2: merge into GeoPackage ───────────────────────────────────────────
echo
echo "=== Phase 2: Merging into GeoPackage (this takes a while) ==="
first_chart=true
for f in "${CHART_FILES[@]}"; do
  name="$(basename "$(dirname "$f")")"
  echo "  $name"
  if $first_chart; then
    ogr2ogr -f GPKG \
      -t_srs EPSG:4326 \
      -nlt PROMOTE_TO_MULTI \
      -makevalid \
      "$COMBINED_GPKG" "$f" 2>/dev/null \
      && first_chart=false \
      || echo "    Warning: skipped (ogr2ogr error)"
  else
    ogr2ogr -f GPKG \
      -append -update \
      -t_srs EPSG:4326 \
      -nlt PROMOTE_TO_MULTI \
      -makevalid \
      "$COMBINED_GPKG" "$f" 2>/dev/null \
      || echo "    Warning: skipped (ogr2ogr error)"
  fi
done

# ── Phase 3: export per-layer GeoJSON ────────────────────────────────────────
echo
echo "=== Phase 3: Exporting per-layer GeoJSON ==="
layers=$(ogrinfo "$COMBINED_GPKG" 2>/dev/null | grep "^[0-9]" | awk '{print $2}' | sort -u)
for layer in $layers; do
  echo "  $layer"
  ogr2ogr -f GeoJSON "$GEOJSON_DIR/${layer}.geojson" "$COMBINED_GPKG" "$layer" 2>/dev/null \
    || echo "    Warning: export failed for $layer"
done

# ── Phase 4: tippecanoe ───────────────────────────────────────────────────────
echo
echo "=== Phase 4: Building vector tiles ==="

LOW_ZOOM_LAYERS=(
  ADMARE BCNLAT BCNSPP BOYLAT BOYSAW BOYSPP CBLARE CBLSUB COALNE CONZNE
  COSARE CTNARE DAYMAR DEPARE DEPCNT DMPGRD DSID EXEZNE FAIRWY FOGSIG
  LAKARE LIGHTS LNDARE LNDELV LNDMRK LNDRGN M_COVR M_NPUB M_NSYS MAGVAR
  MARCUL MIPARE MORFAC OBSTRN OFSPLF PILBOP PIPSOL RDOSTA RESARE RIVERS
  RTPBCN SBDARE SEAARE SILTNK SLCONS SOUNDG TESARE TOPMAR TSEZNE TSSBND
  TSSLPT UNSARE UWTROC WRECKS
)
HIGH_ZOOM_EXTRA=(BUAARE BUISGL)

low_files=()
for l in "${LOW_ZOOM_LAYERS[@]}"; do
  [ -f "$GEOJSON_DIR/${l}.geojson" ] && low_files+=("$GEOJSON_DIR/${l}.geojson")
done

high_files=("${low_files[@]}")
for l in "${HIGH_ZOOM_EXTRA[@]}"; do
  [ -f "$GEOJSON_DIR/${l}.geojson" ] && high_files+=("$GEOJSON_DIR/${l}.geojson")
done

echo "  tippecanoe zoom 0–7 (${#low_files[@]} layers)..."
tippecanoe \
  --output=/tmp/enc_z0-7.mbtiles \
  --minimum-zoom=0 --maximum-zoom=7 \
  --force --read-parallel \
  "${low_files[@]}"

echo "  tippecanoe zoom 8–9 (${#high_files[@]} layers)..."
tippecanoe \
  --output=/tmp/enc_z8-9.mbtiles \
  --minimum-zoom=8 --maximum-zoom=9 \
  --force --read-parallel \
  "${high_files[@]}"

# ── Phase 5: tile-join ────────────────────────────────────────────────────────
echo
echo "=== Phase 5: Joining into enc_combined.mbtiles ==="
tile-join \
  --output="$MBTILES_DIR/enc_combined.mbtiles" \
  --force \
  /tmp/enc_z0-7.mbtiles \
  /tmp/enc_z8-9.mbtiles

echo
echo "=== Done! ==="
echo "Output: $MBTILES_DIR/enc_combined.mbtiles"
echo "Work dir kept at $WORK_DIR — delete when satisfied:"
echo "  rm -rf $WORK_DIR"
