const status = document.getElementById("status");
const romInput = document.getElementById("rom-input");
const search = document.getElementById("search");
const clearBtn = document.getElementById("clear");
const results = document.getElementById("results");
const version = document.getElementById("version");

let romBytes = null;
let ready = false;

const go = new Go();
WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject)
  .then((result) => {
    go.run(result.instance);
    ready = true;
    version.textContent = "v" + window.mudoraVersion;
    status.textContent = "Choose a ROM file. It is parsed entirely in your browser and never uploaded.";
  })
  .catch((err) => {
    status.textContent = "Failed to load WASM module: " + err;
  });

romInput.addEventListener("change", async () => {
  const file = romInput.files[0];
  if (!file) return;
  romBytes = new Uint8Array(await file.arrayBuffer());
  search.disabled = false;
  clearBtn.disabled = false;
  status.textContent = `Inspecting ${file.name}`;
  render(search.value);
});

search.addEventListener("input", () => render(search.value));
clearBtn.addEventListener("click", () => {
  search.value = "";
  render("");
});

function render(query) {
  if (!ready || !romBytes) return;

  const raw = window.mudoraInspect(romBytes, query);
  const parsed = JSON.parse(raw);
  if (parsed && parsed.error) {
    status.textContent = "Error: " + parsed.error;
    return;
  }

  results.innerHTML = "";
  const collapsed = query.trim() === "";
  for (const group of parsed) {
    results.appendChild(buildRegion(group, collapsed));
  }
}

function buildRegion(group, collapsed) {
  const section = document.createElement("div");
  section.className = "region" + (collapsed ? "" : " expanded");

  const header = document.createElement("div");
  header.className = "region-header";
  header.addEventListener("click", () => section.classList.toggle("expanded"));

  const name = document.createElement("span");
  name.className = "region-name";
  name.textContent = group.region;
  header.appendChild(name);

  for (const loc of group.locations) {
    if (loc.progression && loc.icon) {
      header.appendChild(makeIcon(loc.icon, loc.item));
    }
  }

  const rows = document.createElement("div");
  rows.className = "region-rows";
  for (const loc of group.locations) {
    rows.appendChild(buildRow(loc));
  }

  section.appendChild(header);
  section.appendChild(rows);
  return section;
}

function buildRow(loc) {
  const row = document.createElement("div");
  row.className = "row";

  const locLabel = document.createElement("span");
  locLabel.textContent = loc.location;

  const icon = loc.icon ? makeIcon(loc.icon, loc.item) : document.createElement("span");

  const itemLabel = document.createElement("span");
  itemLabel.textContent = loc.item;

  row.appendChild(locLabel);
  row.appendChild(icon);
  row.appendChild(itemLabel);
  return row;
}

function makeIcon(src, alt) {
  const img = document.createElement("img");
  img.className = "icon";
  img.src = src;
  img.alt = alt;
  return img;
}
