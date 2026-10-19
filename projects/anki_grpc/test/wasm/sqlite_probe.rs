//! Cross-platform SQLite command for WASI-to-native storage acceptance.

use rusqlite::Connection;

fn verify(connection: &Connection) -> rusqlite::Result<()> {
    let data: Vec<u8> =
        connection.query_row("SELECT payload FROM proof WHERE id = 1", [], |row| {
            row.get(0)
        })?;
    assert_eq!(data.len(), 11 * 1024 * 1024);
    assert!(data.starts_with(b"wasi-sqlite"));
    assert!(data[b"wasi-sqlite".len()..]
        .iter()
        .all(|byte| *byte == 0xa5));
    let integrity: String = connection.query_row("PRAGMA integrity_check", [], |row| row.get(0))?;
    assert_eq!(integrity, "ok");
    Ok(())
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let arguments: Vec<String> = std::env::args().collect();
    let action = arguments.get(1).ok_or("missing write/verify action")?;
    let path = arguments.get(2).ok_or("missing database path")?;
    let mut connection = Connection::open(path)?;
    // Anki selects exclusive locking before WAL. WASI has no shared-memory VFS.
    connection.execute_batch(
        "PRAGMA locking_mode = exclusive;
         PRAGMA journal_mode = wal;
         PRAGMA temp_store = memory;",
    )?;
    match action.as_str() {
        "write" => {
            let transaction = connection.transaction()?;
            transaction
                .execute_batch("CREATE TABLE proof (id INTEGER PRIMARY KEY, payload BLOB)")?;
            let mut payload = vec![0xa5; 11 * 1024 * 1024];
            payload[..b"wasi-sqlite".len()].copy_from_slice(b"wasi-sqlite");
            transaction.execute("INSERT INTO proof VALUES (1, ?1)", [&payload])?;
            transaction.commit()?;
            let transaction = connection.transaction()?;
            transaction.execute("UPDATE proof SET payload = X'00'", [])?;
            transaction.rollback()?;
        }
        "verify" => {}
        _ => return Err("expected write or verify".into()),
    }
    verify(&connection)?;
    connection.execute_batch("PRAGMA wal_checkpoint(TRUNCATE)")?;
    connection.close().map_err(|(_, error)| error)?;
    println!("sqlite-roundtrip-ok");
    Ok(())
}
