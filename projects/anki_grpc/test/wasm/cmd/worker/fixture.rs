//! Executed full Rust CLI fixture for protocol, lifetime and timezone behavior.

use anki_grpc_server::Dispatcher;
use std::sync::{
    atomic::{AtomicU32, Ordering},
    Arc,
};
use tonic::{
    metadata::{MetadataMap, MetadataValue},
    Code, Status,
};

/// COUNT records shared backend state across independent client connections.
static COUNT: AtomicU32 = AtomicU32::new(0);

/// check_timezone exercises Chrono's local and UTC conversions through the host.
fn check_timezone(path: &[u8]) {
    use chrono::{Local, MappedLocalTime, TimeZone, Utc};

    if path == b"/fixture/TimezoneApia" {
        assert!(matches!(
            Local.with_ymd_and_hms(2011, 12, 30, 12, 0, 0),
            MappedLocalTime::None
        ));
        assert_eq!(
            Local
                .with_ymd_and_hms(2011, 12, 29, 12, 0, 0)
                .unwrap()
                .offset()
                .local_minus_utc(),
            -36_000
        );
        assert_eq!(
            Local
                .with_ymd_and_hms(2011, 12, 31, 12, 0, 0)
                .unwrap()
                .offset()
                .local_minus_utc(),
            50_400
        );
        return;
    }
    let (january, july, gap, overlap) = if path == b"/fixture/TimezoneNewYork" {
        (-18_000, -14_400, (2024, 3, 10, 2, 30), (2024, 11, 3, 1, 30))
    } else {
        (39_600, 37_800, (2024, 10, 6, 2, 15), (2024, 4, 7, 1, 45))
    };
    for (month, offset) in [(1, january), (7, july)] {
        let utc = Utc.with_ymd_and_hms(2024, month, 15, 12, 0, 0).unwrap();
        assert_eq!(
            Local
                .timestamp_opt(utc.timestamp(), 0)
                .unwrap()
                .offset()
                .local_minus_utc(),
            offset
        );
        assert_eq!(
            Local
                .with_ymd_and_hms(2024, month, 15, 12, 0, 0)
                .unwrap()
                .offset()
                .local_minus_utc(),
            offset
        );
    }
    assert!(matches!(
        Local.with_ymd_and_hms(gap.0, gap.1, gap.2, gap.3, gap.4, 0),
        MappedLocalTime::None
    ));
    let MappedLocalTime::Ambiguous(earliest, latest) =
        Local.with_ymd_and_hms(overlap.0, overlap.1, overlap.2, overlap.3, overlap.4, 0)
    else {
        panic!("daylight-saving overlap was not ambiguous");
    };
    assert_eq!(earliest.offset().local_minus_utc(), january.max(july));
    assert_eq!(latest.offset().local_minus_utc(), january.min(july));
}

/// Fixture supplies deterministic protobuf behavior through the actual Rust server.
struct Fixture;

impl Dispatcher for Fixture {
    fn dispatch(&self, path: &str, request: &[u8]) -> Result<Vec<u8>, Status> {
        match path {
            "/fixture/Echo" => Ok(request.to_vec()),
            "/fixture/NativeError" => {
                let mut metadata = MetadataMap::new();
                metadata.insert_bin("anki-error-bin", MetadataValue::from_bytes(request));
                Err(Status::with_metadata(
                    Code::Unknown,
                    "Anki command failed",
                    metadata,
                ))
            }
            "/fixture/Count" => {
                let value = (COUNT.fetch_add(1, Ordering::Relaxed) + 1).to_string();
                let mut bytes = vec![0x0a, value.len() as u8];
                bytes.extend(value.as_bytes());
                Ok(bytes)
            }
            "/fixture/TimezoneNewYork" | "/fixture/TimezoneLordHowe" | "/fixture/TimezoneApia" => {
                check_timezone(path.as_bytes());
                Ok(request.to_vec())
            }
            "/fixture/Trap" => core::arch::wasm32::unreachable(),
            "/fixture/Loop" => loop {
                core::hint::spin_loop();
            },
            _ => Err(Status::unimplemented(format!(
                "unknown Anki command: {path}"
            ))),
        }
    }
}

/// main executes the complete shared CLI with one initialized fixture dispatcher.
fn main() {
    std::process::exit(anki_grpc_server::run_with_backend(
        std::env::args().collect(),
        || Ok(Arc::new(Fixture)),
    ));
}
