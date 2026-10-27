// The pinned Lunr library precedes this file in Hugo's worker bundle.
(() => {
    let ready;
    self.onmessage = async ({ data }) => {
        const { id, type } = data;
        try {
            if (type === "initialize") {
                ready = (async () => {
                    const response = await fetch(data.url);
                    if (!response.ok)
                        throw new Error(
                            `Search index HTTP ${response.status}`,
                        );
                    const payload = await response.json();
                    if (payload.index.version !== lunr.version)
                        throw new Error("Search index version mismatch");
                    return {
                        index: lunr.Index.load(payload.index),
                        documents: new Map(
                            payload.documents.map((document) => [
                                document.ref,
                                document,
                            ]),
                        ),
                    };
                })();
                await ready;
                self.postMessage({ id });
            } else if (type === "query") {
                if (!ready) throw new Error("Search is not initialized");
                const { index, documents } = await ready;
                const results = index
                    .query((query) => {
                        lunr.tokenizer(data.query.toLowerCase()).forEach(
                            (token) => {
                                const term = token.toString();
                                query.term(term, { boost: 100 });
                                query.term(term, {
                                    wildcard:
                                        lunr.Query.wildcard.LEADING |
                                        lunr.Query.wildcard.TRAILING,
                                    boost: 10,
                                });
                                query.term(term, { editDistance: 2 });
                            },
                        );
                    })
                    .slice(0, data.limit)
                    .map(({ ref }) => documents.get(ref));
                self.postMessage({ id, results });
            } else throw new Error("Unsupported search request");
        } catch (error) {
            self.postMessage({ id, error: error.message });
        }
    };
})();
