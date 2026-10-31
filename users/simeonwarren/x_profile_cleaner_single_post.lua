-- X profile cleaner runtime globals; see projects/x_profile_cleaner/api/cleaner.lua.
--# selene: allow(undefined_variable)

-- Fetch only the target day; this policy does not traverse threads.
set_time_window("2026-10-09T00:00:00Z", "2026-10-10T00:00:00Z")

-- Single-post lifecycle demonstration; no age rules or thread expansion.
-- https://x.com/difangtianyuan/status/2108673614795198522
delete_post(
    "2108673614795198522",
    { name = "Single-post lifecycle demonstration" }
)
