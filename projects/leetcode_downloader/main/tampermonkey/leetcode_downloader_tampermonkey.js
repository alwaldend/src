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

const _QUERY_QUESTION_LIST = {
    operationName: "problemsetQuestionListV2",
    query: `
        query problemsetQuestionListV2($filters: QuestionFilterInput, $limit: Int, $searchKeyword: String, $skip: Int, $sortBy: QuestionSortByInput, $categorySlug: String) {
            problemsetQuestionListV2(
              filters: $filters
              limit: $limit
              searchKeyword: $searchKeyword
              skip: $skip
              sortBy: $sortBy
              categorySlug: $categorySlug
            ) {
                questions {
                    id
                    titleSlug
                    title
                    translatedTitle
                    questionFrontendId
                    paidOnly
                    difficulty
                    topicTags {
                        name
                        slug
                        nameTranslated
                    }
                    status
                    isInMyFavorites
                    frequency
                    acRate
                    contestPoint
                }
                totalLength
                finishedLength
                hasMore
            }
        }
    `,
    variables: {
        categorySlug: "all-code-essentials",
        filters: {
            acceptanceFilter: {},
            companyFilter: {
                companySlugs: [],
                operator: "IS",
            },
            contestPointFilter: {
                contestPoints: [],
                operator: "IS",
            },
            difficultyFilter: {
                difficulties: [],
                operator: "IS",
            },
            filterCombineType: "ALL",
            frequencyFilter: {},
            frontendIdFilter: {},
            languageFilter: {
                languageSlugs: [],
                operator: "IS",
            },
            lastSubmittedFilter: {},
            positionFilter: {
                operator: "IS",
                positionSlugs: [],
            },
            premiumFilter: {
                operator: "IS",
                premiumStatus: [],
            },
            publishedFilter: {},
            statusFilter: {
                operator: "IS",
                questionStatuses: ["SOLVED"],
            },
            topicFilter: {
                operator: "IS",
                topicSlugs: [],
            },
        },
        filtersV2: {
            acceptanceFilter: {},
            companyFilter: {
                companySlugs: [],
                operator: "IS",
            },
            contestPointFilter: {
                contestPoints: [],
                operator: "IS",
            },
            difficultyFilter: {
                difficulties: [],
                operator: "IS",
            },
            filterCombineType: "ALL",
            frequencyFilter: {},
            frontendIdFilter: {},
            languageFilter: {
                languageSlugs: [],
                operator: "IS",
            },
            lastSubmittedFilter: {},
            positionFilter: {
                operator: "IS",
                positionSlugs: [],
            },
            premiumFilter: {
                operator: "IS",
                premiumStatus: [],
            },
            publishedFilter: {},
            statusFilter: {
                operator: "IS",
                questionStatuses: ["SOLVED"],
            },
            topicFilter: {
                operator: "IS",
                topicSlugs: [],
            },
        },
        limit: 100,
        searchKeyword: "",
        skip: 0,
        sortBy: {
            sortField: "LAST_SUBMITTED_TIME",
            sortOrder: "DESCENDING",
        },
    },
};

const _QUERY_SUBMISSIONS = {
    operationName: "submissionList",
    query: `
        query submissionList($offset: Int!, $limit: Int!, $lastKey: String, $questionSlug: String!, $lang: Int, $status: Int) {
            questionSubmissionList(
                offset: $offset
                limit: $limit
                lastKey: $lastKey
                questionSlug: $questionSlug
                lang: $lang
                status: $status
            ) {
                lastKey
                hasNext
                submissions {
                    id
                    title
                    titleSlug
                    status
                    statusDisplay
                    lang
                    langName
                    runtime
                    timestamp
                    url
                    isPending
                    memory
                    hasNotes
                    notes
                    flagType
                    frontendId
                    topicTags {
                        id
                    }
                }
            }
        }
    `,
    variables: {
        lastKey: null,
        limit: 20,
        offset: 0,
        questionSlug: "",
    },
};

class App {
    constructor() {
        /** Set to true when when the download is ongoing */
        this.downloading = false;
        const btn = document.createElement("button");
        btn.textContent = "Download submissions";
        btn.style.position = "fixed";
        btn.style.bottom = "0";
        btn.style.right = "0";
        btn.style.zIndex = "999999";
        const text = "Download submissions";
        btn.textContent = text;
        this.textPrefix = text;
        btn.style.padding = "0.5em 1.2em";
        btn.style.fontSize = "1rem";
        btn.style.backgroundColor = "black";
        btn.addEventListener("click", () => this.onClick());
        this.btn = btn;
        const fileInput = document.createElement("input");
        fileInput.type = "file";
        const app = this;
        fileInput.onchange = (event) => {
            const file = event.target.files[0];
            const reader = new FileReader();
            reader.readAsText(file, "UTF-8");
            reader.onload = (readerEvent) => {
                app.headers = app.parseHeadersFromHAR(
                    JSON.parse(readerEvent.target.result),
                );
                const data = {
                    questions: [],
                    submissions: [],
                };
                this.fetchQuestions(data, 0, 20, "");
            };
        };
        this.fileInput = fileInput;
        this.headers = {};
    }

    /** Start the app */
    run() {
        document.body.appendChild(this.btn);
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
    saveSubmissions(submissions) {
        this.saveData(JSON.stringify(submissions), "submissions.json");
        this.btn.textContent = `${this.textPrefix}: saved submissions`;
    }

    /**
     * Save data to downloads
     * @param {any} data - Data to save
     * @param {string} fileName - Filename to use
     * */
    saveData(data, fileName) {
        const a = document.createElement("a");
        a.style = "display: none";
        const blob = new Blob([data], {
            type: "octet/stream",
        });
        const url = window.URL.createObjectURL(blob);
        a.href = url;
        a.download = fileName;
        a.click();
        window.URL.revokeObjectURL(url);
    }

    /**
     * Recursive function fetchinhg questions and adding them to the submissions object
     * */
    fetchQuestions(submissions, offset, limit) {
        const button = this.btn;
        const prefix = this.textPrefix;
        button.textContent = `${prefix}: Fetching questions ${offset}-${offset + limit}`;
        const app = this;
        const requestData = {
            ..._QUERY_QUESTION_LIST,
            skip: offset,
            limit: limit,
        };
        GM_xmlhttpRequest({
            url: `https://leetcode.com/graphql/`,
            method: "POST",
            data: requestData,
            headers: app.headers,
            onload: async function (response) {
                if (response.status !== 200) {
                    button.textContent = `${prefix}: failed to fetch ${offset}-${offset + limit}`;
                    const res = {
                        response: response,
                        responseText: response.responseText,
                        requestHeaders: app.headers,
                        requestData: requestData,
                    };
                    app.saveData(
                        JSON.stringify(res, null, 4),
                        "failed_response.json",
                    );
                    return;
                }
                button.textContent = `${prefix}: finished ${offset}-${offset + limit}`;
                const json = JSON.parse(response.responseText);
                const data = json.data.problemsetQuestionListV2;
                submissions.questions = submissions.questions.concat(
                    data.questions,
                );
                if (!data.hasMore || !app.downloading) {
                    app.saveSubmissions(submissions);
                    return;
                }
                await app.sleep(5 * 1000);
                app.fetchQuestions(submissions, offset + limit, limit);
            },
        });
    }

    /** Button click event listener */
    onClick() {
        if (this.downloading) {
            this.downloading = false;
            return;
        }
        this.downloading = true;
        this.fileInput.click();
    }

    /** Parse headers from a HAR object */
    parseHeadersFromHAR(har) {
        const res = {};
        for (const entry of har.log.entries) {
            for (const header of entry.request.headers) {
                res[header.name] = header.value;
            }
        }
        return res;
    }
}

(function () {
    "use strict";
    const app = new App();
    app.run();
})();
