//! Exercises the real Anki HTTP sync-server entry points in a guest with no capabilities.

fn main() {
    use anki::sync::http_server::SimpleServer;

    assert!(!SimpleServer::is_running());
    let Err(error) = SimpleServer::new(std::path::Path::new("/must-not-be-created")) else {
        panic!("HTTP sync-server construction must be unsupported");
    };
    assert!(error.to_string().contains("unsupported"));
    let error = futures::executor::block_on(SimpleServer::run())
        .expect_err("HTTP sync-server startup must be unsupported");
    println!("{error}");
}
