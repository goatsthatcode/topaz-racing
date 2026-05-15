const TILE_PATH = /^\/tiles\/([^/]+)\/(\d+)\/(\d+)\/(\d+)$/;
const CORS = { 'Access-Control-Allow-Origin': '*' };

export default {
  async fetch(request, env) {
    if (request.method === 'OPTIONS') {
      return new Response(null, {
        headers: {
          ...CORS,
          'Access-Control-Allow-Methods': 'GET, OPTIONS',
          'Access-Control-Max-Age': '86400',
        },
      });
    }

    const url = new URL(request.url);
    const match = url.pathname.match(TILE_PATH);

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
    const bytes = new Uint8Array(data);
    const isGzipped = bytes[0] === 0x1f && bytes[1] === 0x8b;

    const headers = {
      'Content-Type': 'application/vnd.mapbox-vector-tile',
      'Cache-Control': 'public, max-age=86400',
      ...CORS,
    };

    if (isGzipped) {
      const ds = new DecompressionStream('gzip');
      const stream = new Response(data).body.pipeThrough(ds);
      return new Response(stream, { headers });
    }

    return new Response(data, { headers });
  },
};
