// Shell Gym web UI. Plain ES modules, no frameworks, no build step.
// All markup lives in <template> elements in index.html; this file only
// clones templates and fills in data (no HTML-in-JS strings).

const $ = (sel, el = document) => el.querySelector(sel);

// tpl clones a <template> by id and returns its first element.
function tpl(id) {
  return document.getElementById(id).content.firstElementChild.cloneNode(true);
}

const state = {
  path: null,        // /api/path payload
  idx: -1,           // current scene index
  autoAdvance: true,
  sceneEl: null,     // current .scene element
  animating: false,
  debugOpen: false,
  live: false,
  ws: null,
};

// ---------- bootstrap -------------------------------------------------------

async function boot() {
  connectWS();
  try {
    const status = await fetchJSON('/api/status');
    state.live = !!status.live;
    if (state.live) $('#btn-debug').hidden = true;
  } catch { /* status is informational */ }
  await refreshPath();
  $('#path-title').textContent = state.path.title;
  await showScene(startSceneIndex(), 0);
  bindToolbar();
}

function startSceneIndex() {
  const scenes = state.path.scenes;
  if (state.path.current) {
    const i = scenes.findIndex((s) => s.kind === 'unit' && s.id === state.path.current);
    if (i >= 0 && scenes[i].status !== 'completed') return i;
    if (i >= 0) {
      for (let j = i + 1; j < scenes.length; j++) if (!sceneDone(scenes[j])) return j;
      return i;
    }
  }
  for (let j = 0; j < scenes.length; j++) if (!sceneDone(scenes[j])) return j;
  return scenes.length - 1;
}

function sceneDone(s) {
  return s.kind === 'unit' ? s.status === 'completed' : s.status === 'seen';
}

async function refreshPath() {
  state.path = await fetchJSON('/api/path');
  renderProgress();
}

// ---------- rendering -------------------------------------------------------

function renderProgress() {
  const { completed, total } = state.path;
  $('#progress-label').textContent = `${completed}/${total}`;
  $('#progress-fill').style.width = total ? `${(completed / total) * 100}%` : '0';
}

async function showScene(idx, dir) {
  // dir: 1 = forward, -1 = back, 0 = initial
  const scenes = state.path.scenes;
  if (idx < 0 || idx >= scenes.length || state.animating) return;
  const meta = scenes[idx];
  state.animating = true;

  let el;
  try {
    el = meta.kind === 'unit' ? await buildUnitScene(meta) : await buildModuleScene(meta);
  } catch (e) {
    el = buildErrorScene(meta, e);
  }

  const stage = $('#stage');
  const old = state.sceneEl;
  el.classList.add(dir >= 0 ? 'enter-right' : 'enter-left');
  if (dir === 0) el.classList.remove('enter-right', 'enter-left');
  stage.appendChild(el);
  el.scrollTop = 0;

  requestAnimationFrame(() => requestAnimationFrame(() => {
    el.classList.remove('enter-right', 'enter-left');
    el.classList.add('center');
    if (old) {
      old.classList.remove('center');
      old.classList.add(dir >= 0 ? 'exit-left' : 'exit-right');
      setTimeout(() => old.remove(), 420);
    }
    setTimeout(() => { state.animating = false; }, 400);
  }));

  state.idx = idx;
  state.sceneEl = el;

  if (meta.kind === 'unit' && meta.status !== 'completed') {
    fetch(`/api/activate/${meta.id}`, { method: 'POST' }).catch(() => {});
  }
  if (state.debugOpen) refreshDebug();
}

async function buildUnitScene(meta) {
  const unit = await fetchJSON(`/api/unit/${meta.id}`);
  const el = tpl('tpl-unit-scene');
  el.dataset.unit = unit.id;
  $('.kicker-module', el).textContent = moduleTitle(meta.moduleId);
  $('.scene-title', el).textContent = unit.title;
  // unit.html is server-rendered trusted content (the author's markdown)
  $('.scene-body', el).innerHTML = unit.html;
  for (const t of unit.tasks || []) {
    const box = $(`.task-box[data-task="${cssEsc(t.name)}"]`, el);
    if (!box) continue;
    box.dataset.mode = t.mode;
    setTaskStatus(box, t.status);
    if (t.hint && !isDone(t.status)) setTaskHint(box, t.hint);
  }
  if (unit.status === 'completed') el.classList.add('done');
  return el;
}

async function buildModuleScene(meta) {
  const mod = await fetchJSON(`/api/module/${meta.id}`);
  const el = tpl('tpl-module-scene');
  el.dataset.module = mod.id;
  $('.scene-body', el).innerHTML = mod.html;
  $('.module-continue', el).addEventListener('click', async () => {
    await fetch(`/api/module-seen/${mod.id}`, { method: 'POST' });
    meta.status = 'seen';
    next();
  });
  return el;
}

function buildErrorScene(meta, err) {
  const el = tpl('tpl-error-scene');
  $('.err-scene-id', el).textContent = meta.id;
  $('.err-text', el).textContent = String(err);
  return el;
}

function moduleTitle(moduleId) {
  const m = state.path.scenes.find((s) => s.kind === 'module' && s.id === moduleId);
  return m ? m.title : moduleId.replaceAll('-', ' ');
}

// ---------- task updates ----------------------------------------------------

const DONE = new Set(['completed', 'satisfied']);
const isDone = (st) => DONE.has(st);

function setTaskStatus(box, status) {
  box.dataset.status = status;
}

function setTaskHint(box, hint) {
  const h = $('.task-section-hint', box);
  if (!h) return;
  h.textContent = hint;
  h.hidden = !hint;
}

function currentUnitId() {
  return state.sceneEl?.dataset.unit || null;
}

function taskBox(name) {
  return state.sceneEl?.querySelector(`.task-box[data-task="${cssEsc(name)}"]`);
}

function onTaskEvent(d) {
  if (d.unit !== currentUnitId()) return;
  const box = taskBox(d.task);
  if (box) {
    setTaskStatus(box, d.status);
    if (isDone(d.status)) setTaskHint(box, '');
  }
}

function onHintEvent(d) {
  if (d.unit !== currentUnitId()) return;
  const box = taskBox(d.task);
  if (box && !isDone(box.dataset.status)) setTaskHint(box, d.hint);
}

function onInitEvent(d) {
  if (d.unit !== currentUnitId() || d.ok) return;
  const holder = $('.init-errors', state.sceneEl);
  if (!holder) return;
  const div = tpl('tpl-init-error');
  $('.init-task-name', div).textContent = d.task;
  holder.appendChild(div);
}

async function onUnitEvent(d) {
  const scene = state.path.scenes.find((s) => s.kind === 'unit' && s.id === d.unit);
  if (scene && d.status === 'completed' && scene.status !== 'completed') {
    scene.status = 'completed';
    state.path.completed++;
    renderProgress();
    if (d.unit === currentUnitId()) celebrate();
  } else if (scene && d.status !== scene.status && scene.status !== 'completed') {
    scene.status = d.status;
  }
}

function celebrate() {
  const el = state.sceneEl;
  el.classList.add('done');
  const ov = tpl('tpl-celebrate');
  el.appendChild(ov);
  setTimeout(() => {
    ov.remove();
    if (state.autoAdvance) next();
  }, 1400);
}

// ---------- navigation ------------------------------------------------------

function next() { showScene(state.idx + 1, 1); }
function prev() { showScene(state.idx - 1, -1); }

function bindToolbar() {
  $('#btn-next').addEventListener('click', next);
  $('#btn-prev').addEventListener('click', prev);
  $('#chk-auto').addEventListener('change', (e) => {
    state.autoAdvance = e.target.checked;
  });
  $('#btn-theme').addEventListener('click', toggleTheme);
  $('#btn-map').addEventListener('click', openMap);
  $('#btn-map-close').addEventListener('click', () => { $('#map-overlay').hidden = true; });
  $('#map-overlay').addEventListener('click', (e) => {
    if (e.target === $('#map-overlay')) $('#map-overlay').hidden = true;
  });
  $('#btn-help').addEventListener('click', () => { $('#help-overlay').hidden = false; });
  $('#btn-help-close').addEventListener('click', () => { $('#help-overlay').hidden = true; });
  $('#help-overlay').addEventListener('click', (e) => {
    if (e.target === $('#help-overlay')) $('#help-overlay').hidden = true;
  });
  $('#btn-debug').addEventListener('click', toggleDebug);
  $('#btn-debug-close').addEventListener('click', toggleDebug);
  document.addEventListener('keydown', (e) => {
    if (e.target.closest('input,textarea')) return;
    if (e.key === 'ArrowRight') next();
    else if (e.key === 'ArrowLeft') prev();
    else if (e.key === 'm') openMap();
    else if (e.key === '?') $('#help-overlay').hidden = false;
    else if (e.key === 'd' && !state.live) toggleDebug();
    else if (e.key === 'Escape') {
      $('#map-overlay').hidden = true;
      $('#help-overlay').hidden = true;
    }
  });
}

function toggleTheme() {
  const root = document.documentElement;
  const next = root.dataset.theme === 'dark' ? 'light' : 'dark';
  root.dataset.theme = next;
  localStorage.setItem('sg-theme', next);
}

// ---------- path map --------------------------------------------------------

async function openMap() {
  await refreshPath();
  const body = $('#map-body');
  body.replaceChildren();
  let lastModule = null;
  state.path.scenes.forEach((s, i) => {
    if (s.moduleId !== lastModule) {
      lastModule = s.moduleId;
      const h = tpl('tpl-map-module');
      h.textContent = moduleTitle(s.moduleId);
      body.appendChild(h);
    }
    if (s.kind === 'module') return;
    const item = tpl('tpl-map-item');
    item.classList.toggle('current', i === state.idx);
    const st = $('.st', item);
    st.classList.add(s.status);
    st.textContent = s.status === 'completed' || s.status === 'active' ? '●' : '○';
    $('.map-title', item).textContent = s.title;
    item.addEventListener('click', () => {
      $('#map-overlay').hidden = true;
      showScene(i, i >= state.idx ? 1 : -1);
    });
    body.appendChild(item);
  });
  $('#map-overlay').hidden = false;
}

// ---------- debug drawer ----------------------------------------------------

function toggleDebug() {
  if (state.live) return;
  state.debugOpen = !state.debugOpen;
  $('#debug-drawer').hidden = !state.debugOpen;
  if (state.debugOpen) refreshDebug();
}

async function refreshDebug() {
  const unit = currentUnitId();
  $('#debug-unit').textContent = unit || '(module scene)';
  const body = $('#debug-body');
  if (!unit) { body.replaceChildren(); return; }
  const tasks = await fetchJSON(`/api/debug/${unit}`).catch(() => []);
  body.replaceChildren();
  for (const t of tasks || []) {
    const h = tpl('tpl-debug-task');
    h.textContent = t.name;
    body.appendChild(h);
    for (const r of (t.runs || []).slice(-5).reverse()) body.appendChild(runRow(r));
  }
}

function runRow(r) {
  const div = tpl('tpl-debug-run');
  div.classList.add(r.exitCode === 0 ? 'ok' : 'fail');
  const when = new Date(r.startedAt).toLocaleTimeString();
  const to = r.timedOut ? ' (timeout)' : '';
  $('.meta', div).textContent =
    `${when} · ${r.kind} · exit ${r.exitCode}${to} · ${r.durationSec.toFixed(2)}s`;
  if (r.stdout) div.appendChild(pre(r.stdout, ''));
  if (r.stderr) div.appendChild(pre(r.stderr, 'err'));
  return div;
}

function pre(text, cls) {
  const p = document.createElement('pre');
  p.className = cls;
  p.textContent = text;
  return p;
}

function onRunEvent(d) {
  if (!state.debugOpen || d.unit !== currentUnitId()) return;
  refreshDebugSoon();
}

let debugTimer = null;
function refreshDebugSoon() {
  clearTimeout(debugTimer);
  debugTimer = setTimeout(refreshDebug, 300);
}

// ---------- websocket -------------------------------------------------------

function connectWS() {
  const proto = location.protocol === 'https:' ? 'wss' : 'ws';
  const ws = new WebSocket(`${proto}://${location.host}/api/events`);
  state.ws = ws;
  ws.onmessage = (msg) => {
    let ev;
    try { ev = JSON.parse(msg.data); } catch { return; }
    switch (ev.type) {
      case 'task': onTaskEvent(ev.data); break;
      case 'hint': onHintEvent(ev.data); break;
      case 'unit': onUnitEvent(ev.data); break;
      case 'init': onInitEvent(ev.data); break;
      case 'run': onRunEvent(ev.data); break;
    }
  };
  ws.onopen = () => { $('#conn-banner').hidden = true; };
  ws.onclose = () => {
    $('#conn-banner').hidden = false;
    setTimeout(connectWS, 1500);
  };
}

// ---------- utils -----------------------------------------------------------

async function fetchJSON(url) {
  const resp = await fetch(url);
  if (!resp.ok) throw new Error(`${url}: ${resp.status} ${await resp.text()}`);
  return resp.json();
}

function cssEsc(s) {
  return CSS && CSS.escape ? CSS.escape(s) : s.replace(/"/g, '\\"');
}

boot();
