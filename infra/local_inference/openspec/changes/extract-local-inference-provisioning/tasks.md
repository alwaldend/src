## 1. Independent Terraform source

- [x] 1.1 Extract Terraform and stage-local AL configuration without changing
      provisioning resources, inputs, provider pins or state identity.
- [x] 1.2 Verify offline Terraform init/validate, packaged wiring and owner checks.

## 2. Host-only delivery

- [x] 2.1 Remove provisioning files and credentials from PR #233, preserving
      host configuration and signed SSH injection.
- [x] 2.2 Validate and prepare both candidates; delivery receipts own publication.

## 3. Authorized provisioning continuation

- [x] 3.1 Verify provider start behavior and required deployment prerequisites.
- [x] 3.2 Obtain scope approval for the reviewed folder/account prerequisite
      and resolve the required image ID and administrator network inputs.
- [ ] 3.3 Apply the reviewed provisioning plan, allowing only the authorized
      initial boot; stop that instance immediately and verify STOPPED.
      Attempted on 2026-11-03: six supporting resources created, VM creation
      blocked by quotas and allocation failures, no VM created or booted.
      The A100 and Gen2 alternatives also failed. See evidence.md.
- [x] 3.4 Record actual results and remove private plan/state artifacts.

## 4. Authorized availability-zone change

- [x] 4.1 Select ru-central1-b and review replacement of the unused zonal
      resources under the operator's explicit replacement authorization.
- [x] 4.2 Attempt the reviewed zone-b plan and record its result. The operator
      subsequently superseded further zone-b provisioning with section 5.
      Attempted on 2026-11-03: subnet replaced and old unused disk deleted;
      Yandex rejected the new disk with an allocation restriction. No VM exists.
- [x] 4.3 Restore disk destruction protection, verify the resulting state,
      record actual results, and remove private plan artifacts.

## 5. Smaller V100 in zone a

- [x] 5.1 Review a fresh ru-central1-a plan using gpu-standard-v2, one V100,
      eight vCPUs and 48 GiB host RAM.
- [x] 5.2 Attempt the reviewed smaller plans and record their outcomes.
      Further 48 GiB attempts are superseded by the operator's original-size
      choice in section 6; no VM was created or booted by these attempts.
      Attempted on 2026-11-03: disk created and subnet replaced in zone a;
      48 GiB VM creation still exceeded memory quota plus one unidentified
      quota. No VM exists; do not retry against unchanged limits.
- [x] 5.3 Verify actual state, preserve guards, record results and clean up
      private plan artifacts.

## 6. Original V100 after the quota increase

- [x] 6.1 Review and validate a fresh gpu-standard-v1 plan with eight vCPUs
      and 96 GiB RAM in ru-central1-a, preserving the six supporting resources.
- [x] 6.2 Attempt the reviewed regular plan and record its result. Further
      attempts are superseded by the operator's preemptible choice in section 7.
      Attempted on 2026-11-03: create request accepted, but allocation operation
      fhmk1pqi4vth1k6f4oj6 failed with Not enough resources; no VM exists.
- [x] 6.3 Verify actual state, record the result and remove private artifacts.

## 7. Capacity alternatives

- [x] 7.1 Inspect current zone status, supported GPU/zone combinations, quotas
      and public GPU prices; prepare concrete alternatives.
- [x] 7.2 Obtain the operator's choice of interruption behavior or larger
      hardware before changing the current regular V100 configuration.
      The operator selected the preemptible V100 after its behavior was explained.
- [ ] 7.3 Review and apply the selected alternative, stop any new VM after
      its initial boot, and record the actual outcome.
      Both preemptible variants were attempted: 96 GiB operation
      fhmp4c79geskjtismueq and 48 GiB operation fhmljpjc6622s8uo7jlk failed
      allocation with Not enough resources. The control-plane reads found no VM.

## 8. Larger GPU and compatibility attempts

- [x] 8.1 Attempt A100 80 GB in zone a and Gen2 80 GB in zone d through the
      owning Terraform apply, preserving the authorized initial-stop scope.
- [x] 8.2 After the reported zone-a CPU increase, attempt A100 in both regular
      and preemptible modes, then with explicit legacy PCI topology V1.
      All four allocation operations failed with Not enough resources;
      no VM was created or booted.
- [x] 8.3 Restore all destruction guards and reconcile source, documentation
      and failure evidence with the last attempted A100 configuration.
