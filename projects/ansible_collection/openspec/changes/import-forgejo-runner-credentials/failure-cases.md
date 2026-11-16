# Failure cases

- Missing id, uuid, token, instance or name must fail without writing a file.
- Name or instance mismatch must fail before writing credentials.
- Existing worker files must not be overwritten or require bootstrap input.
- Imported token values must not appear in task output or diff output.
- Existing registration-token behavior must continue to work.
- A failed prerequisite between install/service phases must not start a fresh
  worker through forced restart handlers.
