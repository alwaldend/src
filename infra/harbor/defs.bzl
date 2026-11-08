"""infra harbor configuration targets."""

load("//infra:al_config.bzl", "infra_kubernetes_login")

def infra_harbor_kubernetes_login(name, labels = {}, config_name = None, **kwargs):
    """Declare the kubernetes login configuration target."""
    infra_kubernetes_login(
        name = name,
        config_name = config_name,
        labels = labels,
        cluster_ca = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUJkakNDQVIyZ0F3SUJBZ0lCQURBS0JnZ3Foa2pPUFFRREFqQWpNU0V3SHdZRFZRUUREQmhyTTNNdGMyVnkKZG1WeUxXTmhRREUzT0RNM01EWXhNVGt3SGhjTk1qWXdOekV3TVRZMU5URTVXaGNOTXpZd056QTNNVFkxTlRFNQpXakFqTVNFd0h3WURWUVFEREJock0zTXRjMlZ5ZG1WeUxXTmhRREUzT0RNM01EWXhNVGt3V1RBVEJnY3Foa2pPClBRSUJCZ2dxaGtqT1BRTUJCd05DQUFSQnlDb2VTM2diS3NiVTIxR2tFOWxta0RMcDlUWm1wYVR0Y1JidlhrRXMKNFBzVU5wTFZjS0RKWWwzOXp2U3NsSWNKM3JqQUUvdmRhZUdOc09EeXNja01vMEl3UURBT0JnTlZIUThCQWY4RQpCQU1DQXFRd0R3WURWUjBUQVFIL0JBVXdBd0VCL3pBZEJnTlZIUTRFRmdRVWVWaTd1MVhPZE1TV0FsNHRsdTQzCmtsaklIcVF3Q2dZSUtvWkl6ajBFQXdJRFJ3QXdSQUlnYzB3RWhtVmRVeVU1MUlhTjRUNElNTENsbDJnbVF4eW0KaU0xSU5KcEVkUklDSUFJbi9xZzdoNW4zTUFLR3NGRm1tSXJ3UnVkWVkrTmRLd2VGS1hudXk3aDQKLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=",
        oidc = {
            "client_id": "d8LmH5d7CHKwOQ925CGllVDkytJ1mYgn",
            "name": "src_infra_harbor_cluster_provider",
            "redirect_uri": "https://unused",
            "scope": "openid user email groups",
        },
        **kwargs
    )
