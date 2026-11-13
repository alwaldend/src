-- Plan one explicit read for each reply/quote target in the export.
-- This does not fetch data during evaluation or follow further references.
--# selene: allow(undefined_variable)
local selected = {}

for _, post in ipairs(archive.tweets()) do
    for _, reference in ipairs(post.references) do
        if
            (reference.kind == "REPLY" or reference.kind == "QUOTE")
            and not selected[reference.id]
        then
            plan.read_post(reference.id)
            selected[reference.id] = true
        end
    end
end
