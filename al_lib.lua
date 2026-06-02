local M = {}

function M.plugin(t)
    local name = t.name
    config({ plugins = {  [name] = t }})
end

function M.auth(t)
    config({ auth = {t}})
end

function M.vault(t)
    config({ vaults = {t}})
end

return M
