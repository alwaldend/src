// Package luaapi owns descriptor-derived operation argument shapes shared by Lua
// runtime bindings and generated editor declarations. Protobuf owns record fields,
// enums and operation variants; generated Lua is editor-only code. Archive
// records remain worker-local even though their schema uses protobuf.
package luaapi
