## Design

Keep one authoritative Mermaid source beside the lifecycle page. Build and Run
are the two top-level phase subgraphs. Build contains source AL configs and a
nested Runfiles subgraph containing Merged AL config, Plugin runfiles and the
Command executable. The merged configuration points to its configuration and
plugin runfile dependencies, labeled Composed from and Runtime dependencies
respectively. Source configs are build inputs outside Runfiles.
Route both dependency edges through a shared junction so their labels clear
the command execution edge.
Run contains AL and Plugins. AL reads the merged configuration, starts plugin
subprocesses, and runs the command executable with a prepared environment.
A separate bidirectional edge represents gRPC requests and responses over stdio.
A dashed Plugins-to-AL edge shows resources carried by that connection, labeled
Resources (env, files, secrets, etc). Its ResourceResult payload carries
environment values, file paths and opaque data. The AL node keeps its short
label; the built-in injector remains implicit within it.
Plugins point to External services outside both phases. Keep lifecycle detail
in the surrounding prose. Use stacked documents for source AL configs and Plugin runfiles, and stacked
processes for Plugins and External services to show multiplicity.

AL remains a standalone module. Its repository integration in tools/al owns the
site-themed render and publishes both SVG variants into the lifecycle page's
asset directory. Use the regular Mermaid flowchart layout. Generated images
are build outputs. A manual dark WebP target provides repeatable visual
inspection without checked-in binary projections or image postprocessing.
