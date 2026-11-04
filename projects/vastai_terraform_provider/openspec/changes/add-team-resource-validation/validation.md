# Validation

The existing Terraform CLI acceptance suite exercises non-compute CRUD, import,
unchanged plans, failure diagnostics, sensitive values, generated team/key
lifecycles, and multiple teams using the merged reusable module. Fixtures use
synthetic credentials and literal loopback endpoints without live API access.

Ownership fixtures reproduced duplicate-name deletion and abandoned management
keys before the resource fix. The corrected provider retains private creation
confirmation and captured account identity, including after external deletion,
malformed owner metadata, and failed revocation followed by retry. Repeatable
team-ownership-*.json reports retain those observations.

The role fixture uses test_engineering/read through create, refresh, import,
permission updates, and deletion. team-resources.json records encoded GET and
DELETE paths. Module coverage verifies two teams and three project keys under
one synthetic bootstrap account; team-module.json records the lifecycle.

This candidate contains no checked-in Terraform configuration, AppRole,
Vault setup, injector, backend, or live credential. PR #249 retains that
integration. The provider landing date follows trunk's frontmatter convention;
the README is undated. AL receipts and local review evidence identify the exact
candidate, selected checks, publication, and review outcomes.

The launch-mode fixture reproduces unwanted replacement before normalization
for nine compound, alias, and reordered API representations. With normalization,
unchanged configurations and canonical imports plan no change, while a real
SSH/Jupyter family change still plans replacement. Missing and unsupported
API modes fail refresh without abandoning state. launch-mode-*.json reports
record these cases. The full CLI suite also exercises the shared client
configuration helper across six resources and the offer data source.

Mode representations follow the pinned schema's image_runtype example and the
[official client](https://github.com/vast-ai/vast-python/blob/master/vast.py#L2205-L2222).

The complete suite also passes after separating registration, resources,
data sources, and provider data. Registration imports factories; concrete
resource types remain private. Client extraction is shared without importing
the registration package from resources or data sources.

Consumers now retain the shared providerdata.Data payload. The bootstrap fixture
also replaces project keys after name and permission changes, enforcing revocation
before recreation and returning a distinct ID and secret for each generation.
team-bootstrap.json records three project-key generations and complete cleanup.
This verifies that valid plans use replacement rather than the defensive Update
guard. The full suite passes after the payload and deletion-clarity changes.

Seven additional CLI cases reproduced successful refreshes overwriting image or
disk state when those fields were missing, null, blank, zero, or negative. Refresh
now rejects these values before assigning any remote attributes. The read-failure
fixture verifies that the original image and disk allocation remain in state,
along with resource ownership. instance-read-failure-*.json records each result.
