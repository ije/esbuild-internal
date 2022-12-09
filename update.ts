import { Untar } from "https://deno.land/std@0.165.0/archive/tar.ts";
import { readAll } from "https://deno.land/std@0.165.0/streams/conversion.ts";
import { readerFromStreamReader } from "https://deno.land/std@0.167.0/streams/reader_from_stream_reader.ts";
import { ensureDir } from "https://deno.land/std@0.165.0/fs/ensure_dir.ts";
import { join, resolve } from "https://deno.land/std@0.165.0/path/mod.ts";

const version = Deno.args[0];
if (!/^\d+\.\d+\.\d+$/.test(version)) {
  throw new Error("invalid version");
}

console.log(`Downloading esbuild-${version}.tar.gz ...`);
const resp = await fetch(
  `https://codeload.github.com/evanw/esbuild/tar.gz/refs/tags/v${version}`,
);
if (resp.status !== 200) {
  console.error(await resp.text());
  Deno.exit(1);
}

const entryList = new Untar(
  readerFromStreamReader(
    resp.body!.pipeThrough<Uint8Array>(new DecompressionStream("gzip"))
      .getReader(),
  ),
);

// clear dirs
for await (const entry of Deno.readDir(".")) {
  if (
    entry.isDirectory && !entry.name.startsWith(".") && entry.name !== "images"
  ) {
    await Deno.remove(join(Deno.cwd(), entry.name), { recursive: true });
  }
}

// write template files
for await (const entry of entryList) {
  const fileName = entry.fileName.slice(`esbuild-${version}/`.length);
  if (
    fileName.startsWith("internal/api_helpers/") ||
    fileName.startsWith("internal/cli_helpers/")
  ) {
    continue;
  }
  if (
    fileName.startsWith("internal/") &&
    entry.type === "directory"
  ) {
    const fp = fileName.slice("internal/".length);
    await ensureDir(fp);
  } else if (
    fileName.startsWith("internal/") &&
    entry.type === "file"
  ) {
    const fp = fileName.slice("internal/".length);
    let code = new TextDecoder().decode(await readAll(entry));
    code = code.replaceAll(
      "github.com/evanw/esbuild/internal",
      "github.com/ije/esbuild-internal",
    );
    await Deno.writeTextFile(fp, code);
  } else if (
    fileName === "go.mod" ||
    fileName === "go.sum"
  ) {
    let code = new TextDecoder().decode(await readAll(entry));
    code = code.replaceAll(
      "github.com/evanw/esbuild",
      "github.com/ije/esbuild-internal",
    );
    await Deno.writeTextFile(fileName, code);
  } else if (
    fileName === "version.txt" ||
    fileName === "CHANGELOG.md" ||
    fileName === "LICENSE.md"
  ) {
    await Deno.writeFile(fileName, await readAll(entry));
  } else {
    continue;
  }
  console.log("write", fileName);
}

console.log("Updated to", version);
