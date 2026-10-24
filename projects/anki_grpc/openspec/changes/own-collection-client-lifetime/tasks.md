## 1. Transport contract

- [x] 1.1 Adapt real transport and application behavior coverage before removing
      the generic collection-client API; retain successful-media shutdown coverage.
- [x] 1.2 Restore transport lifetime to Start, Connection and Stop, return factory results
      directly, and delete the generic wrapper and collection storage state.

## 2. Application ownership

- [x] 2.1 Migrate collection ownership to the application, simplify prepared
      handles, and update consumers, contracts and documentation.
- [x] 2.2 Verify native, WASM and MCP workflows, independent review, strict
      OpenSpec validation and final delivery checks.

The prior unified-storage candidate is superseded by the user's three-method
client direction. Git history and delivery receipts retain its implementation
and validation evidence. Native server destruction outside asynchronous serving
remains part of the final change.

## 3. Shared client capabilities

- [x] 3.1 Define Client as Connection and Stop, embed it in StartableClient with
      Start, and replace the application's separate Collection interface with Client.
- [x] 3.2 Validate affected callers and existing lifecycle workflows, update
      contract documentation, and complete review and delivery checks.

## 4. Constructor and completion boundaries

- [x] 4.1 Replace the one-method ClientFactory interface with a constructor
      function and ToClientFactory, preserving default and nil handling.
- [x] 4.2 Remove public AsyncClient while retaining concrete Done methods and
      application-owned completion observation; verify existing lifecycle behavior.

## 5. Handwritten public symbols

- [x] 5.1 Audit all handwritten public types, functions, methods, fields and
      errors against their real consumers, preserving required ABI and framework exports.
- [x] 5.2 Hide package-only native C bridge declarations and module-only Rust
      CLI parsing/help declarations, plus the private WASM host cleanup method;
      verify existing native and CLI behavior.

## 6. Constructor conversion naming

- [x] 6.1 Rename the constructor adapter to ToClientFactory throughout callers
      and maintained documentation, preserving construction and nil behavior.
