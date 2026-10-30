## Ownership

The original Go store-sync implementation is superseded by the user-directed
Ansible workflow in infra/download. Its historical evidence remains in this
change; maintained synchronization requirements and implementation belong to
the download owner. This tool retains release publication and extracted-site
retention only. No parallel sync implementation or compatibility command remains.

## Local publication collision policy

Local publication rejects an existing artifact or site version. A preflight
rejects known collisions; promotion repeats the check under the project lock
so two concurrent uploads cannot both publish. Upload the complete artifact
directory privately before renaming it into place. Failed transfers expose
no partial release. Extraction also rejects an already extracted version;
link-only selection remains available for rollback. A completed upload can
be followed by a separate extract and link operation.

Normal publication never overwrites versions. The public workflow is sync,
which continues to overwrite same-version public content from local and remove
public-only content. Cover first upload, unchanged and changed duplicate
attempts, retained site-only collision, duplicate extraction, concurrent upload,
and current/content preservation through the existing real SSH fixture before
changing publication code.

## Extracted website retention

The user explicitly authorized deleting old extracted websites after publishing
a new one, while retaining all downloadable archives. Linking switches current
atomically before removing other extracted versions of that same project,
within the existing project lock. Upload-only and extract-only operations do
not remove sites. A missing new index or unsafe directory ancestry fails before
selection or deletion. Recursive removal must not follow symlinks or cross into
another filesystem. Archives, other projects, and private state remain outside
the cleanup root. Rollback re-extracts a retained archive and links it.

Before implementation, extend the real SSH fixture to verify failed selection
preserves the old site; extraction alone retains it; successful selection
removes older extracted versions and preserves archive hashes and other
projects; nested symlinks cannot delete their targets; repeated selection is
safe; a cleanup failure reports failure after selection and can be retried;
and symlinked release ancestry is refused before deletion. Confirm the existing
sync operation propagates this reduced tree while preserving archives.

## Live publication corrections

The 2026.44.0 rollout exposed two differences from the initial loopback fixture.
The content account cannot traverse the SSH login account's private home, so
remote publication commands start from `/`. A private-home SSH fixture first
reproduced the cleanup failure and then passed with this correction.

A restrictive remote umask also affects an rsync destination root when
`--files-from` omits its attributes. Explicitly set the staged archive directory
to 0755 before promotion. A repeat `link` validates the artifact ancestry and
repairs only that selected archive directory's mode; bytes and absent artifact
trees remain untouched. The fixture now starts sshd with umask 077, checks
public readability, and exercises repair of an existing 0750 directory.
