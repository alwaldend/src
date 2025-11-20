local M = {}

function M.setup()
    vim.g.mapleader = " "
    vim.g.maplocalleader = " "
    require("alwaldend/nvim/lazy").setup()
    require("alwaldend/nvim/options").setup()
    require("alwaldend/nvim/autocmd").setup()
    require("alwaldend/nvim/filetypes").setup()
    require("alwaldend/nvim/keymaps").setup_no_deps()
end

return M
