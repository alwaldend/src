local M = {}

function M.env(t)
    config({ env = {t}})
end

function M.auth(t)
    config({ auth = {t}})
end

function M.secret(t)
    config({ secrets = {t}})
end

function M.file(t)
    config({ files = {t}})
end

function M.vault(t)
    config({ vaults = {t}})
end

return M
