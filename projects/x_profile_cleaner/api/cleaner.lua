---@meta
-- Editor-only declarations for the globals bound by internal/policy/policy.go.
-- Runtime behavior is owned by that implementation; do not execute this file.
-- Direct selection functions accept one ID or a dense ID array; empty arrays
-- do nothing. Queries take one ID. Unsupported arguments and options are errors.
-- Predicates register immediately and run sequentially after configuration ends.
-- Each array entry with callbacks uses its own preservation seed context.
--# selene: allow(unused_variable)

---@class CleanerReference
---@field id string Referenced identifier; its foreign content is not loaded.
---@field kind string quoted, replied_to, or retweeted.

---@class CleanerAttachments
---@field media_keys? string[] Attached media identifiers; no media expansion.
---@field poll_ids? string[] Attached poll identifiers.
---@field media_source_ids? string[] Source post IDs for reused media.

---@class CleanerTextEntity
---@field start number Starting offset in the exposed text.
---@field end number Exclusive ending offset.
---@field tag? string Hashtag or cashtag text.
---@field id? string Mentioned user identifier, when returned.
---@field username? string Mentioned username, when returned.
---@field url? string Shortened link, when returned.
---@field expanded_url? string Expanded link, when returned.
---@field display_url? string Display form of the link, when returned.

---@class CleanerEntities
---@field hashtags? CleanerTextEntity[] Returned hashtag spans.
---@field cashtags? CleanerTextEntity[] Returned cashtag spans.
---@field mentions? CleanerTextEntity[] Returned mention spans.
---@field urls? CleanerTextEntity[] Returned URL spans.

---@class CleanerPost
---@field id string Own post or repost activity ID, never a Lua number.
---@field author_id string Authenticated account ID.
---@field kind string post or repost.
---@field source_id? string Reposted original ID used by the unrepost endpoint.
---@field text string Full available text; no foreign source expansion.
---@field language? string Returned language code; absent metadata remains nil.
---@field attachments? CleanerAttachments Returned media and poll identifiers.
---@field entities? CleanerEntities Entities corresponding to text, including long posts.
---@field created_at number Activity time as Unix seconds.
---@field age_days number Fractional days at the planning run's fixed time.
---@field is_reply boolean Includes own replies to foreign posts.
---@field is_quote boolean Whether this post quotes another post.
---@field metrics table<string, number> Returned public counters; absent values are nil.
---@field references CleanerReference[] Typed references from this activity.

---@class CleanerEdge
---@field from_id string Owned post from which traversal proceeds.
---@field to_id string Owned destination; may be a reply parent during thread traversal.
---@field kind "quoted"|"replied_to" Relationship being traversed.
---@field direction "child"|"parent" Parent identifies a reply-parent step.

---@alias CleanerPredicate fun(post: CleanerPost): boolean
---@class CleanerPreservationContext
---@field seed_id string Canonical seed for this traversal; each call gets a fresh context.

---@alias CleanerInclusionPredicate fun(post: CleanerPost, context: CleanerPreservationContext): boolean
---@alias CleanerEdgePredicate fun(edge: CleanerEdge, context: CleanerPreservationContext): boolean

---@class CleanerRuleOptions
---@field name? string Nonempty rule label for plan explanations.

---@class CleanerRepostOptions: CleanerRuleOptions
---@field include_post? CleanerInclusionPredicate Filters the selected repost activity.

---@class CleanerPreservationOptions: CleanerRepostOptions
-- Inclusion skips marking but never prunes descendants.
---@field follow_edge? CleanerEdgePredicate False prunes that edge, not alternate paths.

--- Return a whole-second UTC RFC3339 timestamp using the fixed evaluation time.
--- Accepts finite nonnegative fractional days in the supported discovery range.
--- Performs no inventory read; suitable for either set_time_window bound.
---@param days number Elapsed 24-hour days before this planning run.
---@return string
function days_ago(days) end

--- Limit API discovery, Lua queries, and graph traversal to an activity-time window.
--- Bounds are RFC3339 timestamps at whole seconds, on or after 2010-11-06.
--- Start is inclusive and end exclusive; nil omits either bound. Both nil means
--- unbounded discovery. Set before any inventory read; callbacks cannot change it.
--- Thread seeds or connecting posts outside the window remain unavailable.
---@param start_time? string Inclusive start, or nil for no lower bound.
---@param end_time? string Exclusive end, or nil for no upper bound.
function set_time_window(start_time, end_time) end

--- Request ordinary deletion of one owned post; preservation can prevent it.
---@param id string|string[]
---@param options? CleanerRuleOptions
function delete_post(id, options) end

--- Request deletion of one owned post even if preserved; never cascades.
---@param id string|string[]
---@param options? CleanerRuleOptions
function force_delete_post(id, options) end

--- Request removal of one own repost using its activity ID, not its source ID.
---@param id string|string[]
---@param options? CleanerRuleOptions
function delete_repost(id, options) end

--- Request removal of an own repost even if its activity was preserved.
---@param id string|string[]
---@param options? CleanerRuleOptions
function force_delete_repost(id, options) end

--- Preserve a post and recursively its own reply descendants, never quotes.
--- A filter applies to every candidate, including replies reached indirectly.
---@param id string|string[]
---@param options? CleanerInclusionPredicate|CleanerPreservationOptions
function preserve_post(id, options) end

--- Preserve connected own reply chains and every later quote branch.
--- Replies are followed to both parents and descendants, including a reply that
--- quotes the seed. Quoted ancestors and foreign posts are never traversed.
---@param id string|string[]
---@param options? CleanerInclusionPredicate|CleanerPreservationOptions
function preserve_thread(id, options) end

--- Preserve one repost activity; the foreign source post is not a graph node.
---@param id string|string[]
---@param options? CleanerInclusionPredicate|CleanerRepostOptions
function preserve_repost(id, options) end

--- Select ordinary post deletions, including own quotes and all own replies.
---@param predicate CleanerPredicate
---@param options? CleanerRuleOptions
function delete_post_when(predicate, options) end

--- Select overriding post deletions, including recent or preserved posts.
---@param predicate CleanerPredicate
---@param options? CleanerRuleOptions
function force_delete_post_when(predicate, options) end

--- Select ordinary own repost removals; age refers to the repost action.
---@param predicate CleanerPredicate
---@param options? CleanerRuleOptions
function delete_repost_when(predicate, options) end

--- Select repost removals even when their activities are preserved.
---@param predicate CleanerPredicate
---@param options? CleanerRuleOptions
function force_delete_repost_when(predicate, options) end

--- Preserve each matching post together with its filtered reply descendants.
---@param predicate CleanerPredicate
---@param options? CleanerInclusionPredicate|CleanerPreservationOptions
function preserve_post_when(predicate, options) end

--- Preserve the reply-connected quote thread containing each matching post.
---@param predicate CleanerPredicate
---@param options? CleanerInclusionPredicate|CleanerPreservationOptions
function preserve_thread_when(predicate, options) end

--- Preserve each matching repost activity.
---@param predicate CleanerPredicate
---@param options? CleanerInclusionPredicate|CleanerRepostOptions
function preserve_repost_when(predicate, options) end

--- Read a fresh copy of one collected owned post, or nil if absent or foreign.
---@param id string
---@return CleanerPost?
function get_post(id) end

--- Read a fresh copy of one own repost activity, or nil if absent.
---@param id string
---@return CleanerPost?
function get_repost(id) end

--- Read collected owned posts in stable ID order, including replies and quotes.
---@return CleanerPost[]
function get_posts() end

--- Read own repost activities in stable ID order.
---@return CleanerPost[]
function get_reposts() end

--- Read direct own replies to an own graph node; does not traverse descendants.
---@param id string
---@return CleanerPost[]
function get_replies(id) end

--- Read all direct own quotes of an own graph node, including every branch.
---@param id string
---@return CleanerPost[]
function get_quotes(id) end
