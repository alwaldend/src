-- Derive a profile description from the export, then propose an announcement.
-- Operations remain ordered; neither call changes the account during evaluation.
--# selene: allow(undefined_variable)
local account = archive.account()
local owner = account.username
if owner == nil or owner == "" then
    owner = account.id
else
    owner = "@" .. owner
end

local count = #archive.tweets()
plan.update_profile({
    description = string.format(
        "Archive reviewed: %d activities for %s",
        count,
        owner
    ),
})
plan.create_post("Archive review complete.")
