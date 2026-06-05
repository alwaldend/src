local M = {}

function M.plugin(t)
    t.config = to_json(t.config)
    config({ plugins = { t }})
end

function M.vault_auth(t)
    config({ vault_auth = {t}})
end

function M.vault_conn(t)
    config({ vault_conn = {t}})
end

return M
