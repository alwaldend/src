-- X profile cleaner runtime globals; see projects/x_profile_cleaner/api/cleaner.lua.
--# selene: allow(undefined_variable)

-- Preserve each thread and its connected own quote/reply branches, regardless of age.
preserve_thread({
    "2070942172661133373", -- Main
    "2104477520162853247", -- MLOps
    "2107074222065586585", -- Prompt factory
    "2108303760120484139", -- Anki
    "2112275092659224723", -- Literally me
    "2081468014138958077", -- Books
})

--- Select activity older than 30 elapsed days at the fixed planning time.
---@param item CleanerPost
---@return boolean
local function older_than_30_days(item)
    return item.age_days > 30
end

-- Ordinary deletion respects preservation; posts include all own replies.
delete_post_when(
    older_than_30_days,
    { name = "Posts and replies older than 30 days" }
)
delete_repost_when(older_than_30_days, { name = "Reposts older than 30 days" })
