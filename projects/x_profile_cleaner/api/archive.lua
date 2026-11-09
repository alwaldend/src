---@meta

-- Editor declarations for archive scripts; do not execute this file.
-- Load this file together with api/types.lua.
--# selene: allow(unscoped_variables, unused_variable)

---@alias ArchiveNull userdata
---@alias ArchiveJSONValue boolean|number|string|ArchiveNull|ArchiveObject|ArchiveArray
---@class ArchiveObject
---@field [string] ArchiveJSONValue
local ArchiveObject = {}
---@class ArchiveArray
---@field [integer] ArchiveJSONValue
local ArchiveArray = {}

---@class ArchiveContext
---@field account_id string Validated archive owner; never X authorization.
---@field evaluated_at integer Fixed evaluation time, Unix seconds.
local ArchiveContext = {}

---@type ArchiveContext
context = {}

archive = {}

---@return string[] names Available top-level YTD dataset names, sorted.
function archive.datasets() end

--- Return the validated account, without its YTD wrapper.
---@return ArchiveAccount account
function archive.account() end

--- Return validated followed-account records, without YTD wrappers.
---@return ArchiveFollowing[] accounts
function archive.following() end

--- Return validated liked-post records, without YTD wrappers.
---@return ArchiveLike[] likes
function archive.likes() end

--- Return normalized protobuf tweet records with typed timestamps and metadata.
---@return ArchiveTweet[] tweets
function archive.tweets() end

--- Return raw JSON rows in their original export shape, without a dataset schema.
---@param name string Dataset filename without data/ or .js.
---@return ArchiveJSONValue[] rows
function archive.read_raw(name) end

json = {}

---@type ArchiveNull Explicit JSON null in archive data, distinct from an absent field.
json.null = nil

cleaner = {}

---@param source string Cleanup configuration using api/cleaner.lua declarations.
--- Append the validated cleanup plan and its explanations to this evaluation.
--- The helper returns no value; failed calls abort publication even across pcall.
function cleaner.plan(source) end
