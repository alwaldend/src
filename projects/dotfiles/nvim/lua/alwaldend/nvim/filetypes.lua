local M = {}

function M.setup()
    vim.filetype.add({
        pattern = {
            [".*/ansible/.*.ya?ml"] = "yaml.ansible",
        },
    })
end

return M
