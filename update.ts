#!/usr/bin/env -S deno run --allow-read --allow-write --allow-net

import { UntarStream } from "jsr:@std/tar/untar-stream";
import { ensureDir } from "jsr:@std/fs/ensure-dir";
import { dirname, join } from "jsr:@std/path";

const version = Deno.args[0];
if (!/^\d+\.\d+\.\d+$/.test(version)) {
  throw new Error("invalid version");
}

console.log(`Downloading esbuild-${version}.tar.gz ...`);
const res = await fetch(`https://codeload.github.com/evanw/esbuild/tar.gz/refs/tags/v${version}`);
if (res.status !== 200) {
  console.error(await res.text());
  Deno.exit(1);
}

// clear dirs
for await (const entry of Deno.readDir(".")) {
  if (entry.isDirectory && !entry.name.startsWith(".") && entry.name !== "images") {
    await Deno.remove(join(Deno.cwd(), entry.name), { recursive: true });
  }
}

const entryList = res.body!.pipeThrough<Uint8Array>(new DecompressionStream("gzip")).pipeThrough(new UntarStream());

for await (const entry of entryList) {
  const fileName = entry.path.slice(`esbuild-${version}/`.length);
  if (!entry.readable) {
    continue;
  }
  if (fileName.startsWith("internal/")) {
    if (fileName.startsWith("internal/api_helpers/") || fileName.startsWith("internal/cli_helpers/")) {
      entry.readable.cancel();
      continue;
    }
    const fp = fileName.slice("internal/".length);
    let code = await (new Response(entry.readable).text());
    code = code.replaceAll("github.com/evanw/esbuild/internal", "github.com/ije/esbuild-internal");
    await ensureDir(dirname(fp));
    await Deno.writeTextFile(fp, code);
  } else if (fileName === "go.mod") {
    let code = await (new Response(entry.readable).text());
    code = code.replaceAll(
      "github.com/evanw/esbuild",
      "github.com/ije/esbuild-internal",
    );
    await Deno.writeTextFile(fileName, code);
  } else if (
    fileName === "version.txt"
    || fileName === "CHANGELOG.md"
    || fileName.startsWith("CHANGELOG-20")
    || fileName === "LICENSE.md"
    || fileName === "go.sum"
  ) {
    const f = await Deno.create(fileName);
    await entry.readable.pipeTo(f.writable);
  } else {
    entry.readable.cancel();
    continue;
  }
  console.log("write", fileName);
}

await new Deno.Command("git", { args: ["add", "--all", "."], stderr: "piped", stdout: "piped" }).output();
await new Deno.Command("git", { args: ["commit", "-m", `v${version}`], stderr: "piped", stdout: "piped" }).output();
await new Deno.Command("git", { args: ["tag", `v${version}`], stderr: "piped", stdout: "piped" }).output();
await new Deno.Command("git", { args: ["push", "origin", "--tags"], stderr: "piped", stdout: "piped" }).output();

console.log("Updated to", version);
