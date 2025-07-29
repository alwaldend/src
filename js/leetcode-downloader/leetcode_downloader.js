// ==UserScript==
// @name         Download leetcode submissions
// @namespace    http://tampermonkey.net/
// @version      0.0.0
// @description  Download leetcode submissions
// @author       simeonwarren
// @match        https://leetcode.com/submissions/*
// @icon         https://www.google.com/s2/favicons?sz=64&domain=tampermonkey.net
// @grant        GM_xmlhttpRequest

// ==/UserScript==

function saveData(data, fileName) {
  var a = document.createElement("a");
  document.body.appendChild(a);
  a.style = "display: none";
  var blob = new Blob([data], {
    type: "octet/stream",
  });
  var url = window.URL.createObjectURL(blob);
  a.href = url;
  a.download = fileName;
  a.click();
  window.URL.revokeObjectURL(url);
}

function sleep(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

function fetchSubmissions(submissions, text, button, offset, limit, lastKey) {
  button.textContent = `${text}: Fetching ${offset}-${offset + limit}`;
  GM_xmlhttpRequest({
    url: `https://leetcode.com/api/submissions/?offset=${offset}&limit=${limit}&lastkey=${lastKey}`,
    onload: async function (response) {
      if (response.status !== 200) {
        button.textContent = `${text}: failed to fetch ${offset}-${offset + limit}: ${response.statusText || response.responseText}`;
        return;
      }
      button.textContent = `${text}: finished ${offset}-${offset + limit}`;
      const json = JSON.parse(response.responseText);
      submissions.submissions_dump = submissions.submissions_dump.concat(
        json.submissions_dump,
      );
      if (!json.has_next) {
        saveSubmissions(submissions, button, text);
        return;
      }
      await sleep(5 * 1000);
      fetchSubmissions(
        submissions,
        text,
        button,
        offset + limit,
        limit,
        json.lastKey,
      );
    },
  });
}

function saveSubmissions(submissions, button, text) {
  saveData(JSON.stringify(submissions), "submissions.json");
  button.textContent = `${text}: saved submissions`;
}

function createButton() {
  const btn = document.createElement("button");
  const text = "Download submissions";
  btn.textContent = text;
  btn.style.position = "fixed";
  btn.style.bottom = "0";
  btn.style.right = "0";
  btn.style.zIndex = "999999";
  btn.style.padding = "0.5em 1.2em";
  btn.style.fontSize = "1rem";
  btn.addEventListener("click", () => {
    fetchSubmissions(
      {
        submissions_dump: [],
      },
      text,
      btn,
      0,
      20,
      "",
    );
  });
  return btn;
}

(function () {
  "use strict";

  document.body.appendChild(createButton());
})();
