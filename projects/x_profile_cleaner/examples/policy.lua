-- The cleaner supplies these documented runtime globals; see api/cleaner.lua.
--# selene: allow(undefined_variable)

-- Optional discovery bounds; nil leaves that side open. Threads can be truncated.
-- Set before queries or nonempty direct selections. Adjacent chunks share an end/start.
-- set_time_window("2024-01-01T00:00:00Z", "2024-02-01T00:00:00Z")
-- set_time_window(days_ago(90), days_ago(30))

-- Add collected IDs here. Empty arrays do nothing; unknown IDs fail planning.
local preserved_post_ids = {
    -- "1111111111111111111",
}
local preserved_thread_ids = {
    -- "2222222222222222222",
}
local forced_post_ids = {
    -- "3333333333333333333",
}

--- Select activity older than 30 elapsed days, using the fixed planning time.
---@param item CleanerPost
---@return boolean
local function older_than_30_days(item)
    return item.age_days > 30
end

-- Preserve selected posts with their replies, or connected quote/reply threads.
preserve_post(preserved_post_ids)
preserve_thread(preserved_thread_ids)

delete_post_when(
    older_than_30_days,
    { name = "Posts and replies older than 30 days" }
)
delete_repost_when(older_than_30_days, { name = "Reposts older than 30 days" })

-- Force deletion of the selected posts even when preserved or recent.
force_delete_post(forced_post_ids)
