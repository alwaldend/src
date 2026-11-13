-- Unlike matching exported text and unfollow accounts outside an explicit list.
-- Review/edit the keywords and retained account IDs before submitting this script.
--# selene: allow(undefined_variable)
local keywords = { "giveaway", "sponsored" }
local keep_following = { ["200"] = true }

---@param text string
---@return boolean
local function matches(text)
    text = string.lower(text)
    for _, keyword in ipairs(keywords) do
        if string.find(text, keyword, 1, true) then
            return true
        end
    end
    return false
end

for _, like in ipairs(archive.likes()) do
    if matches(like.text or "") then
        plan.unlike(like.post_id)
    end
end

for _, account in ipairs(archive.following()) do
    if not keep_following[account.account_id] then
        plan.unfollow(account.account_id)
    end
end
