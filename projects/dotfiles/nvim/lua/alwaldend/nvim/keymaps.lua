--[[ Keymaps ]]
local M = {}
local k = vim.keymap

-- Telescope live_grep in git root
-- Function to find the git root directory based on the current buffer's path
local function find_git_root()
    -- Use the current buffer's path as the starting point for the git search
    local current_file = vim.api.nvim_buf_get_name(0)
    local current_dir
    local cwd = vim.fn.getcwd()
    -- If the buffer is not associated with a file, return nil
    if current_file == "" then
        current_dir = cwd
    else
        -- Extract the directory from the current file's path
        current_dir = vim.fn.fnamemodify(current_file, ":h")
    end

    -- Find the Git root directory from the current file's path
    local git_root = vim.fn.systemlist(
        "git -C "
            .. vim.fn.escape(current_dir, " ")
            .. " rev-parse --show-toplevel"
    )[1]
    if vim.v.shell_error ~= 0 then
        print("Not a git repository. Searching on current working directory")
        return cwd
    end
    return git_root
end

local function delete_all_buffers_except_current()
    local bufs = vim.api.nvim_list_bufs()
    local current_buf = vim.api.nvim_get_current_buf()
    for _, i in ipairs(bufs) do
        if i ~= current_buf then
            vim.api.nvim_buf_delete(i, {})
        end
    end
end

--- [Setup keymaps that require no dependencies]
function M.setup_no_deps()
    k.set(
        { "n", "v" },
        "<Space>",
        "<Nop>",
        { desc = "Disable space", silent = true }
    )
    k.set(
        "n",
        "k",
        "v:count == 0 ? 'gk' : 'k'",
        { desc = "Fix word wrap (k)", expr = true, silent = true }
    )
    k.set(
        "n",
        "j",
        "v:count == 0 ? 'gj' : 'j'",
        { desc = "Fix word wrap (j)", expr = true, silent = true }
    )
    k.set(
        "n",
        "[d",
        vim.diagnostic.goto_prev,
        { desc = "Go to previous diagnostic message" }
    )
    k.set(
        "n",
        "]d",
        vim.diagnostic.goto_next,
        { desc = "Go to next diagnostic message" }
    )
    k.set(
        "n",
        "<leader>e",
        vim.diagnostic.open_float,
        { desc = "Open floating diagnostic message" }
    )
    k.set(
        "n",
        "<leader>q",
        vim.diagnostic.setloclist,
        { desc = "Open diagnostics list" }
    )
    k.set(
        "n",
        "<leader>bD",
        delete_all_buffers_except_current,
        { desc = "Delete all buffers except the current" }
    )
    k.set(
        "n",
        "<leader>n",
        "<cmd>Explore<cr><cr>",
        { desc = "Open netrw", silent = true }
    )
end

function M.setup_neogen()
    local opts = { noremap = true, silent = true }
    local neogen = require("neogen")
    k.set("n", "<leader>ng", neogen.generate, opts)
end

-- [[ Telescope keymaps ]]
function M.setup_telescope()
    local builtin = require("telescope.builtin")
    local themes = require("telescope.themes")

    -- Custom live_grep function to search in git root
    local function live_grep_git_root()
        local git_root = find_git_root()
        if git_root then
            builtin.live_grep({ search_dirs = { git_root } })
        end
    end

    local function live_grep_open_files()
        builtin.live_grep({
            grep_open_files = true,
            prompt_title = "Live Grep in Open Files",
        })
    end

    k.set(
        "n",
        "<leader>sr",
        builtin.oldfiles,
        { desc = "[S]earch [R]ecently opened files" }
    )
    k.set("n", "<leader><space>", function()
        builtin.buffers({
            show_all_buffers = false,
            sort_lastused = true,
            sort_mru = true,
        })
    end, { desc = "[ ] Find existing buffers" })
    k.set("n", "<leader>/", function()
        builtin.current_buffer_fuzzy_find(themes.get_dropdown({
            winblend = 10,
            previewer = false,
        }))
    end, { desc = "[/] Fuzzily search in current buffer" })
    k.set(
        "n",
        "<leader>s/",
        live_grep_open_files,
        { desc = "[S]earch [/] in Open Files" }
    )
    k.set(
        "n",
        "<leader>ss",
        builtin.builtin,
        { desc = "[S]earch [S]elect Telescope" }
    )
    k.set(
        "n",
        "<leader>?",
        builtin.git_files,
        { desc = "[?] Search Git Files" }
    )
    k.set("n", "<leader>sf", builtin.find_files, { desc = "[S]earch [F]iles" })
    k.set("n", "<leader>sh", builtin.help_tags, { desc = "[S]earch [H]elp" })
    k.set(
        "n",
        "<leader>sw",
        builtin.grep_string,
        { desc = "[S]earch current [W]ord" }
    )
    k.set(
        "n",
        "<leader>sg",
        builtin.live_grep,
        { desc = "[S]earch by [G]rep" }
    )
    k.set(
        "n",
        "<leader>sG",
        live_grep_git_root,
        { desc = "[S]earch by [G]rep on Git Root" }
    )
    k.set(
        "n",
        "<leader>sd",
        builtin.diagnostics,
        { desc = "[S]earch [D]iagnostics" }
    )
    k.set("n", "<leader>sR", builtin.resume, { desc = "[S]earch [R]esume" })
    --k.set("n", "<c-d>", actions.delete_buffer, { desc = "Delete current buffer" })
end

--[[ LSP keymaps ]]
--  This function gets run when an LSP connects to a particular buffer.
function M.lsp_callback(event)
    local function nmap(keys, func, desc, mode)
        mode = mode or "n"
        vim.keymap.set(
            mode,
            keys,
            func,
            { buffer = event.buf, desc = "LSP: " .. desc }
        )
    end

    local function list_workspace_folders()
        print(vim.inspect(vim.lsp.buf.list_workspace_folders()))
    end

    local builtin = require("telescope.builtin")

    nmap("<leader>rn", vim.lsp.buf.rename, "[R]e[n]ame")
    nmap("<leader>ca", vim.lsp.buf.code_action, "[C]ode [A]ction")
    nmap("gd", builtin.lsp_definitions, "[G]oto [D]efinition")
    nmap("gr", builtin.lsp_references, "[G]oto [R]eferences")
    nmap("gI", builtin.lsp_implementations, "[G]oto [I]mplementation")
    nmap("<leader>D", builtin.lsp_type_definitions, "Type [D]efinition")
    nmap("<leader>ds", builtin.lsp_document_symbols, "[D]ocument [S]ymbols")
    nmap(
        "<leader>ws",
        builtin.lsp_dynamic_workspace_symbols,
        "[W]orkspace [S]ymbols"
    )
    nmap("K", vim.lsp.buf.hover, "Hover Documentation")
    nmap("<C-k>", vim.lsp.buf.signature_help, "Signature Documentation")
    nmap("gD", vim.lsp.buf.declaration, "[G]oto [D]eclaration")
    nmap(
        "<leader>wa",
        vim.lsp.buf.add_workspace_folder,
        "[W]orkspace [A]dd Folder"
    )
    nmap(
        "<leader>wr",
        vim.lsp.buf.remove_workspace_folder,
        "[W]orkspace [R]emove Folder"
    )
    nmap("<leader>wl", list_workspace_folders, "[W]orkspace [L]ist Folders")

    -- This function resolves a difference between neovim nightly (version 0.11) and stable (version 0.10)
    ---@param client vim.lsp.Client
    ---@param method vim.lsp.protocol.Method
    ---@param bufnr? integer some lsp support methods only in specific files
    ---@return boolean
    local function client_supports_method(client, method, bufnr)
        if vim.fn.has("nvim-0.11") == 1 then
            return client:supports_method(method, bufnr)
        else
            return client.supports_method(method, { bufnr = bufnr })
        end
    end

    vim.api.nvim_buf_create_user_command(
        event.buf,
        "LspFormat",
        vim.lsp.buf.format,
        {
            desc = "Format current buffer with LSP",
        }
    )

    local client = vim.lsp.get_client_by_id(event.data.client_id)

    -- The following code creates a keymap to toggle inlay hints in your
    -- code, if the language server you are using supports them
    --
    -- This may be unwanted, since they displace some of your code
    if
        client
        and client_supports_method(
            client,
            vim.lsp.protocol.Methods.textDocument_inlayHint,
            event.buf
        )
    then
        nmap("<leader>th", function()
            vim.lsp.inlay_hint.enable(
                not vim.lsp.inlay_hint.is_enabled({ bufnr = event.buf })
            )
        end, "[T]oggle Inlay [H]ints")
    end
end

--[[ CMP keymaps ]]
function M.get_cmp()
    local cmp = require("cmp")
    local luasnip = require("luasnip")
    local mapping = {
        ["<C-n>"] = cmp.mapping.select_next_item(),
        ["<C-p>"] = cmp.mapping.select_prev_item(),
        ["<C-b>"] = cmp.mapping.scroll_docs(-4),
        ["<C-f>"] = cmp.mapping.scroll_docs(4),
        ["<C-Space>"] = cmp.mapping.complete({}),
        ["<CR>"] = cmp.mapping.confirm({
            behavior = cmp.ConfirmBehavior.Replace,
            select = true,
        }),
        ["<Tab>"] = cmp.mapping(function(fallback)
            if cmp.visible() then
                cmp.select_next_item()
            elseif luasnip.expand_or_locally_jumpable() then
                luasnip.expand_or_jump()
            else
                fallback()
            end
        end, { "i", "s" }),
        ["<S-Tab>"] = cmp.mapping(function(fallback)
            if cmp.visible() then
                cmp.select_prev_item()
            elseif luasnip.locally_jumpable(-1) then
                luasnip.jump(-1)
            else
                fallback()
            end
        end, { "i", "s" }),
    }
    return cmp.mapping.preset.insert(mapping)
end

--[[ Gitsigns keymaps ]]
---@param bufnr number
function M.setup_gitsigns(bufnr)
    local gs = package.loaded.gitsigns

    local function map(mode, l, r, opts)
        opts = opts or {}
        opts.buffer = bufnr
        vim.keymap.set(mode, l, r, opts)
    end
    local function reset_hunk()
        gs.reset_hunk({ vim.fn.line("."), vim.fn.line("v") })
    end
    local function stage_hunk()
        gs.stage_hunk({ vim.fn.line("."), vim.fn.line("v") })
    end
    local function jump_to_prev()
        if vim.wo.diff then
            return "[c"
        end
        vim.schedule(function()
            gs.prev_hunk()
        end)
        return "<Ignore>"
    end
    local function jump_to_next()
        if vim.wo.diff then
            return "]c"
        end
        vim.schedule(function()
            gs.next_hunk()
        end)
        return "<Ignore>"
    end

    map(
        { "n", "v" },
        "]c",
        jump_to_next,
        { expr = true, desc = "Jump to next hunk" }
    )
    map(
        { "n", "v" },
        "[c",
        jump_to_prev,
        { expr = true, desc = "Jump to previous hunk" }
    )
    map("v", "<leader>hs", stage_hunk, { desc = "stage git hunk" })
    map("v", "<leader>hr", reset_hunk, { desc = "reset git hunk" })
    map("n", "<leader>hs", gs.stage_hunk, { desc = "git stage hunk" })
    map("n", "<leader>hr", gs.reset_hunk, { desc = "git reset hunk" })
    map("n", "<leader>hS", gs.stage_buffer, { desc = "git Stage buffer" })
    map("n", "<leader>hu", gs.undo_stage_hunk, { desc = "undo stage hunk" })
    map("n", "<leader>hR", gs.reset_buffer, { desc = "git Reset buffer" })
    map("n", "<leader>hp", gs.preview_hunk, { desc = "preview git hunk" })
    map("n", "<leader>hb", function()
        gs.blame_line({ full = false })
    end, { desc = "git blame line" })
    map("n", "<leader>hd", gs.diffthis, { desc = "git diff against index" })
    map("n", "<leader>hD", function()
        gs.diffthis("~")
    end, { desc = "git diff against last commit" })
    map(
        "n",
        "<leader>tb",
        gs.toggle_current_line_blame,
        { desc = "toggle git blame line" }
    )
    map(
        "n",
        "<leader>td",
        gs.toggle_deleted,
        { desc = "toggle git show deleted" }
    )
    map(
        { "o", "x" },
        "ih",
        ":<C-U>Gitsigns select_hunk<CR>",
        { desc = "select git hunk" }
    )
end

--[[ Telescope keymaps ]]
M.config_treesitter = {
    incremental_selection = {
        enable = true,
        keymaps = {
            init_selection = "<c-space>",
            node_incremental = "<c-space>",
            scope_incremental = "<c-s>",
            node_decremental = "<M-space>",
        },
    },
    textobjects = {
        select = {
            enable = true,
            -- Automatically jump forward to textobj, similar to targets.vim
            lookahead = true,
            keymaps = {
                -- You can use the capture groups defined in textobjects.scm
                ["aa"] = "@parameter.outer",
                ["ia"] = "@parameter.inner",
                ["af"] = "@function.outer",
                ["if"] = "@function.inner",
                ["ac"] = "@class.outer",
                ["ic"] = "@class.inner",
            },
        },
        move = {
            enable = true,
            set_jumps = true, -- whether to set jumps in the jumplist
            goto_next_start = {
                ["]m"] = "@function.outer",
                ["]]"] = "@class.outer",
            },
            goto_next_end = {
                ["]M"] = "@function.outer",
                ["]["] = "@class.outer",
            },
            goto_previous_start = {
                ["[m"] = "@function.outer",
                ["[["] = "@class.outer",
            },
            goto_previous_end = {
                ["[M"] = "@function.outer",
                ["[]"] = "@class.outer",
            },
        },
        swap = {
            enable = true,
            swap_next = {
                ["<leader>a"] = "@parameter.inner",
            },
            swap_previous = {
                ["<leader>A"] = "@parameter.inner",
            },
        },
    },
}

return M
