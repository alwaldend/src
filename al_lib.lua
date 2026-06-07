local M = {}

function M.plugin(t)
    if t.data then
        t.data = to_pb_json(t.data)
    end
    config({ plugins = { t }})
end

function M.plugin_call(t)
    if t.data then
        t.data = to_pb_json(t.data)
    end
    config({ plugin_calls = { t }})
end

function M.plugin_call_res(t)
    M.plugin_call({
        name = t.plugin .. "_" .. t.name,
        plugin = t.plugin,
        data = {
            res = {t}
        }
    })
end

function M.vault_auth(t)
    config({ vault_auth = {t}})
end

function M.vault_conn(t)
    config({ vault_conn = {t}})
end

return M
