/**
 * Knock out light/grey baked backgrounds from 3D PNG icons → RGBA.
 * Edge flood-fill only (connected to corners) — never punches holes
 * in light parts of the subject (e.g. white piggy body).
 *
 * Usage: node scripts/knockout-icon-bg.mjs
 * Prefer running once on opaque source assets.
 */
import sharp from "sharp";
import { readdir } from "fs/promises";
import { join, dirname } from "path";
import { fileURLToPath } from "url";

const __dirname = dirname(fileURLToPath(import.meta.url));
const dir = join(__dirname, "../public/icons/3d");

function colorDist(a, b) {
  const dr = a[0] - b[0];
  const dg = a[1] - b[1];
  const db = a[2] - b[2];
  return Math.sqrt(dr * dr + dg * dg + db * db);
}

async function processFile(file) {
  const path = join(dir, file);
  const { data, info } = await sharp(path)
    .ensureAlpha()
    .raw()
    .toBuffer({ resolveWithObject: true });

  const { width, height, channels } = info;
  const corners = [
    [0, 0],
    [width - 1, 0],
    [0, height - 1],
    [width - 1, height - 1],
  ];
  let sr = 0,
    sg = 0,
    sb = 0;
  for (const [x, y] of corners) {
    const i = (y * width + x) * channels;
    sr += data[i];
    sg += data[i + 1];
    sb += data[i + 2];
  }
  const bg = [sr / 4, sg / 4, sb / 4];
  const luminance = 0.2126 * bg[0] + 0.7152 * bg[1] + 0.0722 * bg[2];
  if (luminance < 140) {
    console.log(`skip (dark corners): ${file}`);
    return;
  }

  const threshold = 48;
  const idx = (x, y) => (y * width + x) * channels;
  const isBg = (x, y) => {
    const i = idx(x, y);
    return colorDist([data[i], data[i + 1], data[i + 2]], bg) < threshold;
  };

  const visited = new Uint8Array(width * height);
  const queue = [];
  for (const [x, y] of corners) {
    if (isBg(x, y)) {
      queue.push(x, y);
      visited[y * width + x] = 1;
    }
  }

  // Also seed top/bottom/left/right edges sparsely
  for (let x = 0; x < width; x += 4) {
    for (const y of [0, height - 1]) {
      const p = y * width + x;
      if (!visited[p] && isBg(x, y)) {
        visited[p] = 1;
        queue.push(x, y);
      }
    }
  }
  for (let y = 0; y < height; y += 4) {
    for (const x of [0, width - 1]) {
      const p = y * width + x;
      if (!visited[p] && isBg(x, y)) {
        visited[p] = 1;
        queue.push(x, y);
      }
    }
  }

  let qi = 0;
  while (qi < queue.length) {
    const x = queue[qi++];
    const y = queue[qi++];
    const i = idx(x, y);
    data[i + 3] = 0;
    const neigh = [
      [x - 1, y],
      [x + 1, y],
      [x, y - 1],
      [x, y + 1],
    ];
    for (const [nx, ny] of neigh) {
      if (nx < 0 || ny < 0 || nx >= width || ny >= height) continue;
      const p = ny * width + nx;
      if (visited[p]) continue;
      if (!isBg(nx, ny)) continue;
      visited[p] = 1;
      queue.push(nx, ny);
    }
  }

  // Soft edge: reduce alpha near remaining bg-like border pixels
  const soft = 22;
  for (let y = 0; y < height; y++) {
    for (let x = 0; x < width; x++) {
      const p = y * width + x;
      if (visited[p]) continue;
      const i = idx(x, y);
      const dist = colorDist([data[i], data[i + 1], data[i + 2]], bg);
      if (dist < threshold + soft) {
        // Only soften if adjacent to knocked-out pixel (edge fringe)
        let nearHole = false;
        for (const [dx, dy] of [
          [-1, 0],
          [1, 0],
          [0, -1],
          [0, 1],
        ]) {
          const nx = x + dx,
            ny = y + dy;
          if (nx < 0 || ny < 0 || nx >= width || ny >= height) continue;
          if (visited[ny * width + nx]) {
            nearHole = true;
            break;
          }
        }
        if (nearHole) {
          const t = Math.max(0, (dist - threshold) / soft);
          data[i + 3] = Math.min(data[i + 3], Math.round(255 * t));
        }
      }
    }
  }

  await sharp(data, { raw: { width, height, channels } })
    .png()
    .toFile(path);
  console.log(`ok: ${file}`);
}

const files = (await readdir(dir)).filter((f) => f.endsWith(".png"));
for (const f of files) {
  await processFile(f);
}
console.log(`done ${files.length} files`);
