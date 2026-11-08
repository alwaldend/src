"""infra flux configuration targets."""

load("//infra:al_config.bzl", "infra_kubernetes_login")

def infra_flux_kubernetes_login(name, labels = {}, config_name = None, **kwargs):
    """Declare the kubernetes login configuration target."""
    infra_kubernetes_login(
        name = name,
        config_name = config_name,
        labels = labels,
        cluster_ca = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUJkakNDQVIyZ0F3SUJBZ0lCQURBS0JnZ3Foa2pPUFFRREFqQWpNU0V3SHdZRFZRUUREQmhyTTNNdGMyVnkKZG1WeUxXTmhRREUzT0RNeE56QTNNell3SGhjTk1qWXdOekEwTVRJeE1qRTJXaGNOTXpZd056QXhNVEl4TWpFMgpXakFqTVNFd0h3WURWUVFEREJock0zTXRjMlZ5ZG1WeUxXTmhRREUzT0RNeE56QTNNell3V1RBVEJnY3Foa2pPClBRSUJCZ2dxaGtqT1BRTUJCd05DQUFTM3F3WkZ4bFFXSFQ0ZVdTa3hNUGNkd2k1eDVsZ0JwVTVXdk0zMndwS2YKSC9BcFhKODNFVU1QZGZrZXpKMWw1KzdlQ3pndjRGdXk0ck10K2dtRFZsampvMEl3UURBT0JnTlZIUThCQWY4RQpCQU1DQXFRd0R3WURWUjBUQVFIL0JBVXdBd0VCL3pBZEJnTlZIUTRFRmdRVTg5OXJ4ZklpMk1XYUMzWmhZWUdiClNsNnFCRGt3Q2dZSUtvWkl6ajBFQXdJRFJ3QXdSQUlnQS80a0FHS0JGbjV3dk43L0Q4V2EwakI2RHNjaHFuR3oKSkp4Q3o1dE1CVWdDSUZVbGgzc3lPTnFGb2NORmphUHRydGRGaWpLMVRoaTNrTHFKWHdab2tSZ08KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=",
        oidc = {
            "client_id": "51Bw8vtyWuoZIm45O7ug83SghbLCxwPE",
            "name": "src_infra_flux_cluster_provider",
            "redirect_uri": "https://unused",
            "scope": "openid user email groups",
        },
        **kwargs
    )
