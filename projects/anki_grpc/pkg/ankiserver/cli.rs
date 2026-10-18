//! Command-line transport selection, validated before starting the server.

use std::path::PathBuf;

/// HELP lists native CLI commands without starting a server.
pub(crate) const HELP: &str = "Usage: anki_grpc <COMMAND>

Commands:
  serve  Run the Anki gRPC server.

Options:
  -h, --help  Show this help and exit.
";

/// SERVE_HELP describes the serve command's transports and listener parameters.
pub(crate) const SERVE_HELP: &str =
    "Usage: anki_grpc serve [--transport stdio|tcp|unix] [--listen ADDRESS_OR_PATH]

Transports:
  stdio  One gRPC connection over stdin/stdout (default); no --listen.
  tcp    Listen at --listen HOST:PORT, for example 127.0.0.1:50051.
  unix   Listen at --listen SOCKET_PATH, for example ./anki.sock.

Options:
  --transport TRANSPORT  Select stdio, tcp, or unix.
  --listen ADDRESS_OR_PATH  Required for tcp and unix.
  -h, --help  Show this help and exit.
";

/// Command selects native CLI behavior before runtime initialization.
pub(crate) enum Command {
    /// Help displays available native commands.
    Help,
    /// ServeHelp displays transport options without starting the server.
    ServeHelp,
    /// Serve runs the server on the selected transport.
    Serve(Transport),
}

impl Command {
    /// parse validates the native subcommand and delegates its transport arguments.
    pub(crate) fn parse(args: impl IntoIterator<Item = String>) -> Result<Self, String> {
        let mut args = args.into_iter();
        match args.next().as_deref() {
            None | Some("-h" | "--help") => Ok(Self::Help),
            Some("serve") => match Transport::parse(args)? {
                Some(transport) => Ok(Self::Serve(transport)),
                None => Ok(Self::ServeHelp),
            },
            Some(command) => Err(format!("unknown command: {command}")),
        }
    }
}

/// Transport identifies the stream source and its explicit listener endpoint.
pub(crate) enum Transport {
    /// Stdio serves a single connection on the process pipes.
    Stdio,
    /// Tcp accepts connections at the specified host and port.
    Tcp(String),
    /// Unix accepts connections at the specified filesystem socket path.
    Unix(PathBuf),
}

impl Transport {
    /// parse validates arguments, returning None when help was requested.
    pub(crate) fn parse(args: impl IntoIterator<Item = String>) -> Result<Option<Self>, String> {
        let mut args = args.into_iter();
        let mut transport = None;
        let mut listen = None;
        while let Some(arg) = args.next() {
            match arg.as_str() {
                "-h" | "--help" => return Ok(None),
                "--transport" => {
                    if transport.is_some() {
                        return Err("--transport may only be specified once".into());
                    }
                    transport = Some(value(&mut args, "--transport")?);
                }
                "--listen" => {
                    if listen.is_some() {
                        return Err("--listen may only be specified once".into());
                    }
                    listen = Some(value(&mut args, "--listen")?);
                }
                _ => return Err(format!("unknown argument: {arg}")),
            }
        }
        match transport.as_deref().unwrap_or("stdio") {
            "stdio" if listen.is_none() => Ok(Some(Self::Stdio)),
            "stdio" => Err("--listen is not supported with stdio".into()),
            "tcp" => Ok(Some(Self::Tcp(
                listen.ok_or("tcp requires --listen HOST:PORT")?,
            ))),
            "unix" => Ok(Some(Self::Unix(PathBuf::from(
                listen.ok_or("unix requires --listen SOCKET_PATH")?,
            )))),
            unknown => Err(format!(
                "unknown transport: {unknown}; choose stdio, tcp, or unix"
            )),
        }
    }
}

/// value reads a required nonempty option value without consuming another flag.
fn value(args: &mut impl Iterator<Item = String>, option: &str) -> Result<String, String> {
    args.next()
        .filter(|value| !value.is_empty() && !value.starts_with('-'))
        .ok_or_else(|| format!("{option} requires a value"))
}
