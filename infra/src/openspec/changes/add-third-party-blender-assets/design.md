## Decisions

- Use the existing manual Bazel download targets and SHA-256 integrity pins.
  Keep asset archives intact when Blender needs sibling files or catalogs.
- Extract the individual CMU BVH clips through `http_archive` so consumers
  receive importable motion files directly.
- Pin the supplied Darlin CDN URLs for Pippa's mesh and diffuse texture. Check
  the mesh against the original attachment and keep the binaries at the source.
- Pin the 41 Mixamo motion sources from the final album to a Hugging Face
  revision and derive their checksums from the production FBX files.
