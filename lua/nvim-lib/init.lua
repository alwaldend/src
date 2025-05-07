local M = {}

function M.setup()
    vim.g.mapleader = " "
    vim.g.maplocalleader = " "
    require("alwaldend/nvim-lib/lazy").setup()
    require("alwaldend/nvim-lib/options").setup()
    require("alwaldend/nvim-lib/autocmd").setup()
    require("alwaldend/nvim-lib/filetypes").setup()
    require("alwaldend/nvim-lib/keymaps").setup_no_deps()
end

return M
