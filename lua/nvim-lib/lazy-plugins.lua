local bg = "#333333"
local fg = "#ffffff"

local colors = {
    a = { bg = fg, fg = bg },
    b = { bg = bg, fg = fg },
    c = { bg = bg, fg = fg },
}
local theme = {
    normal = colors,
    insert = colors,
    visual = colors,
    replace = colors,
    command = colors,
    inactive = colors,
}

local config = {
    ["danymat/neogen"] = { version = "*" },
    ["folke/which-key.nvim"] = { version = "*" },
    ["hrsh7th/nvim-cmp"] = { version = "*" },
    ["lukas-reineke/indent-blankline.nvim"] = { version = "*" },
    ["lewis6991/gitsigns.nvim"] = { version = "*" },
    ["neovim/nvim-lspconfig"] = { version = "*" },
    ["numToStr/Comment.nvim"] = { version = "*" },
    ["nvim-lualine/lualine.nvim"] = { version = "*" },
    ["nvim-telescope/telescope.nvim"] = { version = "*" },
    ["nvim-tree/nvim-tree.lua"] = { version = "*" },
    ["nvim-treesitter/nvim-treesitter"] = { version = "*" },
    ["sainnhe/sonokai"] = { version = "*" },
    ["stevearc/conform.nvim"] = { version = "*" },
    ["tpope/vim-fugitive"] = { version = "*" },
    ["tpope/vim-sleuth"] = { version = "*" },
    ["mason-org/mason.nvim"] = { version = "^1.0.0" },
    ["mason-org/mason-lspconfig.nvim"] = { version = "^1.0.0" },
}

return {
    {
        --- lsps
        "mason-org/mason.nvim",
        version = config["mason-org/mason.nvim"].version,
    },
    {
        --- lspconfig
        "mason-org/mason-lspconfig.nvim",
        version = config["mason-org/mason-lspconfig.nvim"].version,
    },
    {
        --- annotation generator
        "danymat/neogen",
        dependencies = "nvim-treesitter/nvim-treesitter",
        version = config["danymat/neogen"].version,
        config = function()
            require("neogen").setup({ snippet_engine = "luasnip" })
            require("alwaldend/nvim-lib/keymaps").setup_neogen()
        end,
    },
    {
        -- Useful plugin to show you pending keybinds.
        "folke/which-key.nvim",
        opts = {},
        version = config["folke/which-key.nvim"].version,
        config = function()
            local which_key = require("which-key")
            which_key.add({
                { "<leader>c", group = "[C]ode" },
                { "<leader>d", group = "[D]ocument" },
                { "<leader>g", group = "[G]it" },
                { "<leader>h", group = "Git [H]unk" },
                { "<leader>r", group = "[R]ename" },
                { "<leader>s", group = "[S]earch" },
                { "<leader>t", group = "[T]oggle" },
                { "<leader>w", group = "[W]orkspace" },
                { "<leader>", group = "VISUAL <leader>", mode = "v" },
                { "<leader>h", desc = "Git [H]unk", mode = "v" },
            })
        end,
    },
    {
        -- Autocompletion
        "hrsh7th/nvim-cmp",
        version = config["hrsh7th/nvim-cmp"].version,
        dependencies = {
            -- Snippet Engine & its associated nvim-cmp source
            { "L3MON4D3/LuaSnip" },
            { "saadparwaiz1/cmp_luasnip" },

            -- Adds LSP completion capabilities
            { "hrsh7th/cmp-nvim-lsp" },
            { "hrsh7th/cmp-path" },

            -- Adds a number of user-friendly snippets
            { "rafamadriz/friendly-snippets" },
        },
        config = function()
            local cmp = require("cmp")
            local luasnip = require("luasnip")

            require("luasnip.loaders.from_vscode").lazy_load()
            require("luasnip.loaders.from_snipmate").lazy_load()
            luasnip.config.setup({})

            cmp.setup({
                snippet = {
                    expand = function(args)
                        luasnip.lsp_expand(args.body)
                    end,
                },
                mapping = require("alwaldend/nvim-lib/keymaps").get_cmp(),
                sources = {
                    { name = "nvim_lsp" },
                    { name = "luasnip" },
                    { name = "path" },
                },
            })
        end,
    },
    {
        -- Add indentation guides even on blank lines
        "lukas-reineke/indent-blankline.nvim",
        version = config["lukas-reineke/indent-blankline.nvim"].version,
        dependencies = {
            "nvim-treesitter/nvim-treesitter",
        },
        main = "ibl",
        opts = {
            scope = {
                show_start = false,
                show_end = false,
            },
        },
    },
    {
        -- Adds git related signs to the gutter, as well as utilities for managing changes
        "lewis6991/gitsigns.nvim",
        version = config["lewis6991/gitsigns.nvim"].version,
        opts = {
            -- See `:help gitsigns.txt`
            signs = {
                add = { text = "+" },
                change = { text = "~" },
                delete = { text = "_" },
                topdelete = { text = "‾" },
                changedelete = { text = "~" },
            },
            on_attach = require("alwaldend/nvim-lib/keymaps").setup_gitsigns,
        },
    },
    {
        -- LSP Configuration & Plugins
        "neovim/nvim-lspconfig",
        version = config["neovim/nvim-lspconfig"].version,
        dependencies = {
            -- Automatically install LSPs to stdpath for neovim
            { "mason-org/mason.nvim" },
            { "mason-org/mason-lspconfig.nvim" },
            -- Useful status updates for LSP
            { "j-hui/fidget.nvim", version = "*" },
            -- Additional lua configuration, makes nvim stuff amazing!
            { "folke/neodev.nvim" },
        },
        config = function()
            local servers = {
                ansiblels = {},
                gopls = {
                    gopls = {
                        gofumpt = true,
                        workspaceFiles = {
                            "**/BUILD",
                            "**/WORKSPACE",
                            "**/*.{bzl,bazel}",
                        },
                        -- fix for generated golang files
                        -- https://github.com/bazelbuild/rules_go/wiki/Editor-setup
                        env = { GOPACKAGESDRIVER = "gopackagesdriver" },
                        directoryFilters = {
                            "-bazel-bin",
                            "-bazel-out",
                            "-bazel-testlogs",
                        },
                    },
                },
                bashls = {},
                clangd = {},
                pyright = {},
                dockerls = {},
                docker_compose_language_service = {},
                gradle_ls = {},
                helm_ls = {},
                kotlin_language_server = {},
                bzl = {},
                marksman = {},
                taplo = {},
                jsonls = {},
                tflint = {},
                html = {},
                yamlls = {
                    settings = {
                        yaml = {
                            format = {
                                enable = true,
                            },
                            completion = true,
                            hover = true,
                            validate = true,
                            schemaStore = {
                                enable = false,
                            },
                            customTags = {
                                "!reference sequence",
                            },
                            schemas = {
                                ["http://json.schemastore.org/github-workflow"] = {
                                    ".github/workflows/*",
                                },
                                ["http://json.schemastore.org/github-action"] = {
                                    ".github/action.{yml,yaml}",
                                },
                                ["http://json.schemastore.org/ansible-stable-2.9"] = {
                                    "roles/tasks/*.{yml,yaml}",
                                },
                                ["http://json.schemastore.org/prettierrc"] = {
                                    ".prettierrc.{yml,yaml}",
                                },
                                ["http://json.schemastore.org/kustomization"] = {
                                    "kustomization.{yml,yaml}",
                                },
                                ["http://json.schemastore.org/ansible-playbook"] = {
                                    "*play*.{yml,yaml}",
                                },
                                ["http://json.schemastore.org/chart"] = {
                                    "Chart.{yml,yaml}",
                                },
                                ["https://json.schemastore.org/dependabot-v2"] = {
                                    ".github/dependabot.{yml,yaml}",
                                },
                                ["https://gitlab.com/gitlab-org/gitlab/-/raw/master/app/assets/javascripts/editor/schema/ci.json"] = {
                                    "*.gitlab-ci.{yml,yaml}",
                                },
                                ["https://raw.githubusercontent.com/compose-spec/compose-spec/master/schema/compose-spec.json"] = {
                                    "*docker-compose*.{yml,yaml}",
                                },
                                ["https://raw.githubusercontent.com/argoproj/argo-workflows/master/api/jsonschema/schema.json"] = {
                                    "*flow*.{yml,yaml}",
                                },
                            },
                        },
                    },
                },
                lua_ls = {
                    Lua = {
                        workspace = { checkThirdParty = false },
                        telemetry = { enable = false },
                        diagnostics = {
                            globals = { "vim" },
                            disable = { "missing-fields" },
                        },
                    },
                },
            }

            local cmp_nvim_lsp = require("cmp_nvim_lsp")
            local fidget = require("fidget")
            local mason = require("mason")
            local mason_lspconfig = require("mason-lspconfig")
            local lspconfig = require("lspconfig")
            local neodev = require("neodev")
            local keymaps = require("alwaldend/nvim-lib/keymaps")

            local capabilities = cmp_nvim_lsp.default_capabilities(
                vim.lsp.protocol.make_client_capabilities()
            )

            fidget.setup()
            neodev.setup()
            mason.setup()
            mason_lspconfig.setup({
                ensure_installed = vim.tbl_keys(servers),
                automatic_installation = false,
            })
            mason_lspconfig.setup_handlers({
                ---@param server_name string
                function(server_name)
                    lspconfig[server_name].setup({
                        capabilities = capabilities,
                        on_attach = keymaps.setup_lsp,
                        settings = servers[server_name],
                        filetypes = (servers[server_name] or {}).filetypes,
                    })
                end,
            })
        end,
    },
    {
        -- "gc" to comment visual regions/lines
        "numToStr/Comment.nvim",
        version = config["numToStr/Comment.nvim"].version,
        opts = {},
    },
    {
        -- better status line
        "nvim-lualine/lualine.nvim",
        version = config["nvim-lualine/lualine.nvim"].version,
        opts = {
            options = {
                globalstatus = true,
                theme = theme,
                icons_enabled = false,
                component_separators = "|",
                section_separators = "",
            },
        },
    },
    {
        -- Fuzzy Finder (files, lsp, etc)
        "nvim-telescope/telescope.nvim",
        version = config["nvim-telescope/telescope.nvim"].version,
        dependencies = {
            { "nvim-lua/plenary.nvim" },
            -- Fuzzy Finder Algorithm which requires local dependencies to be built.
            -- Only load if `make` is available. Make sure you have the system
            -- requirements installed.
            {
                "nvim-telescope/telescope-fzf-native.nvim",
                build = "make",
                cond = function()
                    return vim.fn.executable("make") == 1
                end,
            },
        },

        config = function(_, opts)
            pcall(require("telescope").load_extension, "fzf")
            require("telescope").setup(opts)
            require("alwaldend/nvim-lib/keymaps").setup_telescope()
        end,
        opts = {
            defaults = {
                disable_devicons = true,
                layout_strategy = "vertical",
                layout_config = {
                    width = 0.9,
                    height = 0.9,
                    preview_cutoff = 1,
                    mirror = true,
                },
                mappings = {
                    i = {
                        ["<C-u>"] = false,
                        ["<C-d>"] = false,
                    },
                },
            },
        },
    },
    {
        --- nvim tree
        "nvim-tree/nvim-tree.lua",
        version = config["nvim-tree/nvim-tree.lua"].version,
        lazy = false,
        enabled = false,
        dependencies = {
            --{ "nvim-tree/nvim-web-devicons" },
        },
        config = function()
            local tree = require("nvim-tree")
            local api = require("nvim-tree.api")
            tree.setup({
                renderer = {
                    icons = {
                        show = {
                            file = false,
                            folder = false,
                            folder_arrow = false,
                            git = true,
                            modified = true,
                            diagnostics = true,
                            bookmarks = true,
                        },
                    },
                },
            })
            api.tree.open()
        end,
    },
    {
        --- treesitter
        "nvim-treesitter/nvim-treesitter",
        version = config["nvim-treesitter/nvim-treesitter"].version,
        dependencies = {
            "nvim-treesitter/nvim-treesitter-textobjects",
        },
        build = ":TSUpdate",
        opts = {
            ensure_installed = {
                "c",
                "cpp",
                "go",
                "lua",
                "python",
                "rust",
                "tsx",
                "javascript",
                "typescript",
                "vimdoc",
                "vim",
                "bash",
                "html",
                "yaml",
                "json",
                "toml",
                "css",
                "bash",
                "dockerfile",
                "gomod",
                "markdown",
                "gowork",
                "groovy",
                "kotlin",
            },
            auto_install = true,
            highlight = {
                enable = true,
                additional_vim_regex_highlighting = false,
            },
            indent = {
                enable = true,
            },
        },
        config = function(_, opts)
            local keymaps =
                require("alwaldend/nvim-lib/keymaps").config_treesitter
            opts["incremental_selection"] = keymaps.incremental_selection
            opts["textobjects"] = keymaps.textobjects
            require("nvim-treesitter.configs").setup(opts)
        end,
    },
    {
        --- theme
        "sainnhe/sonokai",
        version = config["sainnhe/sonokai"].version,
        lazy = false,
        priority = 1000,
        enabled = true,
        config = function()
            vim.g.sonokai_style = "default"
            vim.g.sonokai_colors_override = {
                bg0 = { "#333333", "235" },
            }
            vim.cmd.colorscheme("sonokai")
        end,
    },
    {
        --- autoformatting
        "stevearc/conform.nvim",
        version = config["stevearc/conform.nvim"].version,
        dependencies = {
            "neovim/nvim-lspconfig",
            "nvim-lua/plenary.nvim",
            "mason-org/mason.nvim",
        },
        event = { "BufWritePre" },
        cmd = { "ConformInfo" },
        opts = {
            formatters_by_ft = {
                lua = { "stylua" },
                bzl = { "buildifier" },
                python = { "isort", "black" },
                bash = { "shfmt", "shellcheck" },
                sh = { "shfmt", "shellcheck" },
                go = { "goimports", "gofumpt", "goimports-reviser" },
                javascript = {
                    "prettierd",
                    "prettier",
                    stop_after_first = true,
                },
                typescript = {
                    "prettierd",
                    "prettier",
                    stop_after_first = true,
                },
                vue = { "prettierd", "prettier", stop_after_first = true },
                css = { "prettierd", "prettier", stop_after_first = true },
                scss = { "prettierd", "prettier", stop_after_first = true },
                less = { "prettierd", "prettier", stop_after_first = true },
                html = { "prettierd", "prettier", stop_after_first = true },
                json = { "prettierd", "prettier", stop_after_first = true },
                jsonc = { "prettierd", "prettier", stop_after_first = true },
                markdown = { "prettierd", "prettier", stop_after_first = true },
                yaml = { "prettierd", "prettier", stop_after_first = true },
                ["*"] = { "trim_whitespace", "trim_newlines" },
            },
            format_on_save = {
                timeout_ms = 2000,
                lsp_fallback = true,
            },
            formatters = {
                shfmt = {
                    prepend_args = { "-i", "4" },
                },
            },
        },
    },
    {
        --- Git wrapper
        "tpope/vim-fugitive",
        version = config["tpope/vim-fugitive"].version,
        dependencies = {
            "tpope/vim-rhubarb",
        },
        config = function() end,
    },
    {
        -- Detect tabstop and shiftwidth automatically
        "tpope/vim-sleuth",
        version = config["tpope/vim-sleuth"].version,
        config = function() end,
    },
}
