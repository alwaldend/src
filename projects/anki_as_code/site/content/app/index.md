---
title: Anki as code
linkTitle: Browser app
description: Edit Anki collections and configuration files in your browser
layout: landing
---

Edit the example below, or open an Anki collection package or configuration JSON.

<div id="anki-app">
  <div class="d-flex flex-wrap gap-2 mb-3">
    <button id="open-collection" type="button"
            class="btn btn-outline-primary" disabled>Open collection</button>
    <input id="import-collection" type="file" accept=".colpkg"
           aria-label="Open Anki collection package" hidden disabled>
    <button id="open-config" type="button"
            class="btn btn-outline-primary" disabled>Open JSON</button>
    <input id="import-config" type="file" accept=".json,application/json"
           aria-label="Open configuration JSON" hidden disabled>
    <button id="check-config" type="button"
            class="btn btn-primary" disabled>Check</button>
    <button id="download-config" type="button"
            class="btn btn-outline-primary" disabled>
      Download JSON
    </button>
  </div>
  <div class="d-flex flex-wrap gap-2 mb-2">
    <button id="plan-collection" type="button"
            class="btn btn-outline-primary" disabled>Plan</button>
    <button id="apply-collection" type="button"
            class="btn btn-primary" disabled>Apply</button>
    <button id="download-collection" type="button"
            class="btn btn-outline-primary" disabled>Download collection</button>
    <button id="save-workspace" type="button"
            class="btn btn-outline-primary" disabled>Save locally</button>
    <button id="restore-workspace" type="button"
            class="btn btn-outline-primary" disabled>Restore saved</button>
  </div>
  <p id="collection-status" class="small text-body-secondary">
    Open an Anki collection package to plan and apply changes.
  </p>
  <label class="form-label" for="workspace-files">File</label>
  <select id="workspace-files" class="form-select mb-3" disabled>
    <option value="">Loading workspace…</option>
  </select>
  <label class="form-label" for="file-content">TOML</label>
  <textarea id="file-content" class="form-control font-monospace mb-2"
            rows="12" spellcheck="false" aria-describedby="file-help"
            disabled></textarea>
  <p id="app-status" class="text-body-secondary text-break" role="status"
     aria-live="polite" aria-atomic="true">Loading workspace…</p>
  <div id="collection-report" class="small text-break" hidden></div>
  <p id="file-help" class="small text-body-secondary">
    Save locally keeps your configuration and collection in this browser across
    reloads. Apply updates the collection; saving alone keeps edits unapplied.
    Browser data can be cleared, so keep downloaded backups.
    Restore saved replaces current work with your last local save.
  </p>
  <p class="small text-body-secondary">
    Open collection replaces the collection and configuration. Open JSON replaces
    configuration and keeps the open collection. Work stays on this device.
  </p>
</div>

<noscript>
  Enable JavaScript to use the browser collection workspace.
</noscript>
<script type="module" src="app.js"></script>
