//! Shared protobuf dispatch for transport servers and direct native callers.

pub use anki::backend::{init_backend, Backend};

include!(concat!(env!("OUT_DIR"), "/commands.rs"));

/// CommandError distinguishes missing routes from Anki's serialized native errors.
pub enum CommandError {
    /// Unknown indicates that the generated upstream route does not exist.
    Unknown,
    /// Native preserves Anki's encoded BackendError without reinterpretation.
    Native(Vec<u8>),
}

/// command executes an upstream route using its generated service and method indices.
pub fn command(backend: &Backend, path: &str, request: &[u8]) -> Result<Vec<u8>, CommandError> {
    let (service, method) = command_ids(path).ok_or(CommandError::Unknown)?;
    backend
        .run_service_method(service, method, request)
        .map_err(CommandError::Native)
}
