//! WASI command entry point for the same Rust CLI and server as native builds.

/// main forwards the complete argument vector and preserves the Rust exit status.
fn main() {
    std::process::exit(anki_grpc_server::run(std::env::args().collect()));
}
