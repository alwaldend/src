-- This script produces a plan; it makes no X requests.
-- A backend must validate these operation names and implement their execution.
-- API declarations are in api/archive.lua and api/types.lua.
--# selene: allow(undefined_variable)
for _, row in ipairs(archive.likes()) do
    plan.unlike(row.post_id)
end

for _, row in ipairs(archive.following()) do
    plan.unfollow(row.account_id)
end
