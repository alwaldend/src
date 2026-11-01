// ==UserScript==
// @name         Download leetcode submissions
// @namespace    http://tampermonkey.net/
// @version      0.0.0
// @description  Download leetcode submissions
// @author       simeonwarren
// @match        https://leetcode.com/problemset/*
// @icon         https://www.google.com/s2/favicons?sz=64&domain=tampermonkey.net
// @grant        GM_xmlhttpRequest

// ==/UserScript==

class App {
    constructor() {
        /** Set to true when when the download is ongoing */
        this.downloading = false;
    }

    /**
     * Return a timeout promise
     * @param {number} ms - Timeout ms
     * */
    sleep(ms) {
        return new Promise((resolve) => setTimeout(resolve, ms));
    }

    /**
     * Save submissions to Downloads
     * */
    saveSubmissions(submissions, button, text) {
        this.saveData(JSON.stringify(submissions), "submissions.json");
        button.textContent = `${text}: saved submissions`;
    }

    /**
     * Save data to downloads
     * @param {any} data - Data to save
     * @param {string} fileName - Filename to use
     * */
    saveData(data, fileName) {
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

    /**
     * Recursive function fetchinhg submissions and adding them to the submissions object
     * */
    fetchSubmissions(submissions, text, button, offset, limit, lastKey) {
        button.textContent = `${text}: Fetching ${offset}-${offset + limit}`;
        const app = this;
        GM_xmlhttpRequest({
            url: `https://leetcode.com/api/submissions/?offset=${offset}&limit=${limit}&lastkey=${lastKey}`,
            onload: async function (response) {
                if (response.status !== 200) {
                    button.textContent = `${text}: failed to fetch ${offset}-${offset + limit}: ${response.statusText || response.responseText}`;
                    return;
                }
                button.textContent = `${text}: finished ${offset}-${offset + limit}`;
                const json = JSON.parse(response.responseText);
                submissions.submissions_dump =
                    submissions.submissions_dump.concat(json.submissions_dump);
                if (!json.has_next || !app.downloading) {
                    app.saveSubmissions(submissions, button, text);
                    return;
                }
                await app.sleep(5 * 1000);
                app.fetchSubmissions(
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

    /** Button click event listener */
    onClick(btn) {
        if (this.downloading) {
            this.downloading = false;
            return;
        }
        this.downloading = true;
        const text = "Download submissions";
        btn.textContent = text;
        this.fetchSubmissions(
            {
                submissions_dump: [],
            },
            text,
            btn,
            0,
            20,
            "",
        );
    }

    /** Create the download button */
    createButton() {
        const btn = document.createElement("button");
        btn.textContent = "Download submissions";
        btn.style.position = "fixed";
        btn.style.bottom = "0";
        btn.style.right = "0";
        btn.style.zIndex = "999999";
        btn.style.padding = "0.5em 1.2em";
        btn.style.fontSize = "1rem";
        btn.addEventListener("click", () => this.onClick(btn));
        return btn;
    }

    /** Entrypoint */
    run() {
        document.body.appendChild(this.createButton());
    }
}

(function () {
    "use strict";
    const app = new App();
    app.run();
})();
