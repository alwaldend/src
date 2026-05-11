local infra = require("infra.al_lib")
local lib = require("al_lib")

lib.auth({
    name = "default",
    approle = {
        name = "src_infra_dc1_pve1",
    },
})
