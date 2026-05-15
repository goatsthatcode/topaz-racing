const TILE_PATH = /^\/tiles\/([^/]+)\/(\d+)\/(\d+)\/(\d+)$/;

export default {
  async fetch(request, env) {
    if (request.method === 'OPTIONS') {
      return new Response(null, {
        headers: {
          'Access-Control-Allow-Origin': '*',
          'Access-Control-Allow-Methods': 'GET, OPTIONS',
          'Access-Control-Max-Age': '86400',
        },
      });
    }

    const url = new URL(request.url);
    const match = url.pathname.match(TILE_PATH);

    const CORS = { 'Access-Control-Allow-Origin': '*' };

    if (!match) {
      return new Response('Not found', { status: 404, headers: CORS });
    }

    const [, tileset, z, x, y] = match;
    const key = `${tileset}/${z}/${x}/${y}.pbf`;

    const object = await env.TILES_BUCKET.get(key);
    if (!object) {
      return new Response('Tile not found', { status: 404, headers: CORS });
    }

    const data = await object.arrayBuffer();

    // MBTiles vector tiles are gzip-compressed; detect and declare encoding
    // so the browser/MapLibre decompresses correctly.
    const bytes = new Uint8Array(data, 0, 2);
    const isGzipped = bytes[0] === 0x1f && bytes[1] === 0x8b;

    return new Response(data, {
      headers: {
        'Content-Type': 'application/vnd.mapbox-vector-tile',
        'Cache-Control': 'public, max-age=86400',
        'Access-Control-Allow-Origin': '*',
        ...(isGzipped ? { 'Content-Encoding': 'gzip' } : {}),
      },
    });
  },
};
