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
     * Recursive function fetchinhg questions and adding them to the submissions object
     * */
    fetchQuestions(submissions, text, button, offset, limit) {
        button.textContent = `${text}: Fetching questions ${offset}-${offset + limit}`;
        const app = this;
        GM_xmlhttpRequest({
            url: `https://leetcode.com/graphql/`,
            method: "POST",
            data: {
                ..._QUERY_QUESTION_LIST,
                skip: offset,
                limit: limit,
            },
            headers: {
                "Content-Type": "application/json",
                Referer: "https://leetcode.com/problemset/",
            },
            onload: async function (response) {
                if (response.status !== 200) {
                    button.textContent = `${text}: failed to fetch ${offset}-${offset + limit}: ${response.statusText || response.responseText}`;
                    return;
                }
                button.textContent = `${text}: finished ${offset}-${offset + limit}`;
                const json = JSON.parse(response.responseText);
                const data = json.data.problemsetQuestionListV2;
                submissions.questions = submissions.questions.concat(
                    data.questions,
                );
                if (!data.hasMore || !app.downloading) {
                    app.saveSubmissions(submissions, button, text);
                    return;
                }
                await app.sleep(5 * 1000);
                app.fetchQuestions(
                    submissions,
                    text,
                    button,
                    offset + limit,
                    limit,
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
        const data = {
            questions: [],
            submissions: [],
        };
        this.fetchQuestions(data, text, btn, 0, 20, "");
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
        btn.style.backgroundColor = "black";
        btn.addEventListener("click", () => this.onClick(btn));
        return btn;
    }
}

(function () {
    "use strict";
    const app = new App();
    document.body.appendChild(app.createButton());
})();
