-- Remove repost relationships older than 180 days, using original post IDs.
-- Unresolved sources remain untouched; repeated sources produce one operation.
--# selene: allow(undefined_variable)
local cutoff = context.evaluated_at - 180 * 24 * 60 * 60
local selected = {}

for _, activity in ipairs(archive.tweets()) do
    local source = activity.source_id
    if
        activity.kind == "REPOST"
        and activity.created_at.seconds < cutoff
        and source ~= nil
        and not selected[source]
    then
        plan.unrepost(source)
        selected[source] = true
    end
end
