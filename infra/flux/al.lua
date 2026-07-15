local lib = require("al_lib")
local infra = require("infra.al_lib")

lib.vault_auth({
    name = "default",
    approle = {
        name = "src_infra_flux",
    },
})

infra.kubernetes_login({
    name = "kubernetes_login",
    labels = { k8s = "1" },
    cluster_ca = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUJkakNDQVIyZ0F3SUJBZ0lCQURBS0JnZ3Foa2pPUFFRREFqQWpNU0V3SHdZRFZRUUREQmhyTTNNdGMyVnkKZG1WeUxXTmhRREUzT0RNeE56QTNNell3SGhjTk1qWXdOekEwTVRJeE1qRTJXaGNOTXpZd056QXhNVEl4TWpFMgpXakFqTVNFd0h3WURWUVFEREJock0zTXRjMlZ5ZG1WeUxXTmhRREUzT0RNeE56QTNNell3V1RBVEJnY3Foa2pPClBRSUJCZ2dxaGtqT1BRTUJCd05DQUFTM3F3WkZ4bFFXSFQ0ZVdTa3hNUGNkd2k1eDVsZ0JwVTVXdk0zMndwS2YKSC9BcFhKODNFVU1QZGZrZXpKMWw1KzdlQ3pndjRGdXk0ck10K2dtRFZsampvMEl3UURBT0JnTlZIUThCQWY4RQpCQU1DQXFRd0R3WURWUjBUQVFIL0JBVXdBd0VCL3pBZEJnTlZIUTRFRmdRVTg5OXJ4ZklpMk1XYUMzWmhZWUdiClNsNnFCRGt3Q2dZSUtvWkl6ajBFQXdJRFJ3QXdSQUlnQS80a0FHS0JGbjV3dk43L0Q4V2EwakI2RHNjaHFuR3oKSkp4Q3o1dE1CVWdDSUZVbGgzc3lPTnFGb2NORmphUHRydGRGaWpLMVRoaTNrTHFKWHdab2tSZ08KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=",
    domain = "flux.alwaldend.com",
    oidc = {
        name = "src_infra_flux_cluster_provider";
        scope = "openid user email groups";
        client_id = "51Bw8vtyWuoZIm45O7ug83SghbLCxwPE";
        redirect_uri = "https://unused";
    },
})

infra.ansible_keys({
    name = "ansible_keys",
    labels = { ansible = "1" },
    vault_ssh = {
        backend = "ssh/clients/sign/admins",
        ttl = 60 * 60 * 2
    }
})

lib.plugin_call({
    name = "tf_backend_tf_setup",
    plugin = "tf_backend",
    labels = { tf = "setup" },
    data = {
        vault_secret = "alwaldend.com/vault1/approles/src_infra_flux/tf_backend/tf_setup",
        vault_secret_mount = "secrets"
    },
})

lib.plugin_call({
    name = "tf_backend_tf",
    plugin = "tf_backend",
    labels = { tf = "main" },
    data = {
        vault_secret = "alwaldend.com/vault1/approles/src_infra_flux/tf_backend/tf",
        vault_secret_mount = "secrets"
    },
})

lib.plugin_call({
    name = "pve_login",
    plugin = "pve_login",
    labels = { tf = "setup" },
})

infra.k3s_token({
    name = "k3s_token",
    path = "alwaldend.com/vault1/approles/src_infra_flux/config",
    labels = { ansible = "1" },
})
