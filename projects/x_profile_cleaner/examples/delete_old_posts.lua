-- Retain configured IDs, media, recent posts and unknown/high engagement.
-- Generic worker script; see api/archive.lua and api/types.lua.
--# selene: allow(undefined_variable)
local keep_ids = { ["1001"] = true }
local cutoff = context.evaluated_at - 90 * 24 * 60 * 60

---@param post ArchiveTweet
---@return boolean
local function should_delete(post)
    local likes = post.metrics.like_count
    -- Exact counters outside Lua's numeric range are strings; retain them.
    return post.kind == "POST"
        and not keep_ids[post.id]
        and post.created_at.seconds < cutoff
        and #post.attachments.media_keys == 0
        and type(likes) == "number"
        and likes < 10
end

for _, post in ipairs(archive.tweets()) do
    if should_delete(post) then
        plan.delete_post(post.id)
    end
end
