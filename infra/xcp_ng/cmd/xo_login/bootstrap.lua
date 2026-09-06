-- Loaded by AL at runtime so the selected AppRole never enters build outputs.
local approle = os.getenv("XO_BOOTSTRAP_APPROLE")
assert(
    approle
        and (
            approle:match("^src_[a-z0-9_]+$")
            or approle == "user_simeonwarren"
        ),
    "Set XO_BOOTSTRAP_APPROLE to an approved AppRole name"
)
local issuer = os.getenv("XO_BOOTSTRAP_ISSUER_APPROLE")
if issuer == "" then
    issuer = nil
end
assert(
    not issuer or issuer:match("^src_[a-z0-9_]+$"),
    "Invalid issuer AppRole name"
)

config({
    vault_auth = {
        { name = "xo_bootstrap", approle = { name = issuer or approle } },
    },
    plugins = {
        {
            name = "xo_login",
            bin = "com_alwaldend_src/infra/xcp_ng/cmd/xo_login/xo_login_/xo_login",
            labels = { xoa_login = "1" },
            data = to_pb_json({
                xoa_url = "https://xoa.xcp-ng.alwaldend.com",
                discovery_url = "https://vault.alwaldend.com:8200/v1/identity/oidc/provider/src_infra_xcp_ng_provider/.well-known/openid-configuration",
                vault_auth = "xo_bootstrap",
                issuer_auth = issuer and "xo_bootstrap" or nil,
                target_approle = issuer and approle or nil,
            }),
        },
    },
})
