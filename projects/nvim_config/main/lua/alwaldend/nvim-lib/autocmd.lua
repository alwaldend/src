local M = {}

local autocmd = vim.api.nvim_create_autocmd
local augroup = vim.api.nvim_create_augroup

function M.setup()
    -- highlight on yank
    autocmd({ "TextYankPost" }, {
        callback = function()
            vim.highlight.on_yank()
        end,
        group = augroup("YankHighlight", { clear = true }),
        pattern = "*",
    })
    -- remove trailing whitespace
    --autocmd({ "BufWritePre" }, {
    --    pattern = { "*" },
    --    command = [[%s/\s\+$//e]],
    --})
end

return M
