"""al injector fixture configuration targets."""

load("//:defs.bzl", "al_inject")

def al_injector_fixture_selected(name, labels = {}, config_name = None, **kwargs):
    """Declare the selected configuration target."""
    al_inject(
        name = name,
        config_name = config_name,
        labels = labels,
        res_vault_env = {"vault_environment": {"conn": "default", "auth": "default"}},
        res_env = {
            "VAULT_TOKEN": {"value": "synthetic-child-token", "deps": ["vault_environment"]},
            "FIXTURE_API_KEY": {"value": "{{ .Last.Data.api_key }}", "deps": ["credentials"]},
            "FIXTURE_LAST": {"value": "{{ .Last.Data.value }}", "deps": ["omega", "alpha"]},
            "FIXTURE_OPERATION": {"value": "{{ .Last.Data.value }}", "deps": ["operation"]},
        },
        res_op = {"operation": {"method": "write", "path": "fixture/input", "data": {"value": "synthetic input"}}},
        res_kv = {
            "credentials": {"mount": "secrets", "path": "fixture/credentials"},
            "alpha": {"mount": "secrets", "path": "fixture/alpha"},
            "omega": {"mount": "secrets", "path": "fixture/omega"},
        },
        **kwargs
    )
