//! Standalone argument and exit-status adapter for the shared native Anki CLI.

/// main runs the shared native CLI with the process argument vector.
fn main() {
    std::process::exit(anki_grpc_server::run(std::env::args().collect()));
}
