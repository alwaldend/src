//! Shared Anki server implementation for standalone and linked CLI entry points.
//! Process lifetime and failure recovery belong to the caller.

mod cli;

use std::{convert::Infallible, future::Future, pin::Pin};

use anki::backend::{init_backend, Backend};
use bytes::{Buf, BufMut, Bytes};
use hyper::{body::Incoming, server::conn::http2, service::service_fn};
use hyper_util::rt::{TokioExecutor, TokioIo};
use tokio::{
    io::{AsyncRead, AsyncWrite},
    net::{TcpListener, UnixListener},
};
use tonic::{codec, server::UnaryService, Request, Response, Status};

include!(concat!(env!("OUT_DIR"), "/commands.rs"));

/// ProtobufCodec forwards upstream message bytes without another schema binding.
#[derive(Default)]
struct ProtobufCodec;

impl codec::Codec for ProtobufCodec {
    type Encode = Bytes;
    type Decode = Bytes;
    type Encoder = Self;
    type Decoder = Self;

    fn encoder(&mut self) -> Self {
        Self
    }
    fn decoder(&mut self) -> Self {
        Self
    }
}

impl codec::Encoder for ProtobufCodec {
    type Item = Bytes;
    type Error = Status;

    fn encode(&mut self, item: Bytes, buffer: &mut codec::EncodeBuf<'_>) -> Result<(), Status> {
        buffer.put(item);
        Ok(())
    }
}

impl codec::Decoder for ProtobufCodec {
    type Item = Bytes;
    type Error = Status;

    fn decode(&mut self, buffer: &mut codec::DecodeBuf<'_>) -> Result<Option<Bytes>, Status> {
        Ok(Some(buffer.copy_to_bytes(buffer.remaining())))
    }
}

/// Command dispatches one upstream call against the process-owned backend.
struct Command {
    /// backend shares Anki's synchronized state with all client connections.
    backend: Backend,
    /// path is the gRPC service and method selected by Hyper.
    path: String,
}

impl UnaryService<Bytes> for Command {
    type Response = Bytes;
    type Future = Pin<Box<dyn Future<Output = Result<Response<Bytes>, Status>> + Send>>;

    fn call(&mut self, request: Request<Bytes>) -> Self::Future {
        let backend = self.backend.clone();
        let path = self.path.clone();
        Box::pin(async move {
            tokio::task::spawn_blocking(move || {
                let (service, method) = command_ids(&path).ok_or_else(|| {
                    Status::unimplemented(format!("unknown Anki command: {path}"))
                })?;
                let result = backend
                    .run_service_method(service, method, &request.into_inner())
                    .map_err(|details| {
                        let mut metadata = tonic::metadata::MetadataMap::new();
                        metadata.insert_bin(
                            "anki-error-bin",
                            tonic::metadata::MetadataValue::from_bytes(&details),
                        );
                        Status::with_metadata(tonic::Code::Unknown, "Anki command failed", metadata)
                    })?;
                Ok(Response::new(Bytes::from(result)))
            })
            .await
            .expect("Anki command task")
        })
    }
}

/// run executes the native CLI using a complete argument vector and returns its exit status.
pub fn run(arguments: Vec<String>) -> i32 {
    let mut arguments = arguments.into_iter();
    let program = arguments.next().unwrap_or_else(|| "anki_grpc".into());
    let program = std::path::Path::new(&program)
        .file_name()
        .and_then(|name| name.to_str())
        .unwrap_or("anki_grpc");
    let transport = match cli::Command::parse(arguments) {
        Ok(cli::Command::Serve(transport)) => transport,
        Ok(cli::Command::Help) => {
            print!("{}", cli::HELP.replacen("anki_grpc", program, 1));
            return 0;
        }
        Ok(cli::Command::ServeHelp) => {
            print!("{}", cli::SERVE_HELP.replacen("anki_grpc", program, 1));
            return 0;
        }
        Err(error) => {
            eprintln!("{error}\nSee --help for usage.");
            return 2;
        }
    };
    match serve(transport) {
        Ok(()) => 0,
        Err(error) => {
            eprintln!("Error: {error:?}");
            1
        }
    }
}

/// serve owns the native backend and Tokio runtime for the selected transport.
#[tokio::main]
async fn serve(transport: cli::Transport) -> Result<(), Box<dyn std::error::Error>> {
    // An empty BackendInit uses Anki's default initialization options.
    let backend = init_backend(&[]).map_err(std::io::Error::other)?;
    match transport {
        cli::Transport::Stdio => {
            serve_connection(
                tokio::io::join(tokio::io::stdin(), tokio::io::stdout()),
                backend,
            )
            .await?;
        }
        cli::Transport::Tcp(address) => {
            let listener = TcpListener::bind(&address).await?;
            eprintln!("Listening on tcp://{}", listener.local_addr()?);
            loop {
                let (stream, _) = listener.accept().await?;
                spawn_connection(stream, backend.clone());
            }
        }
        cli::Transport::Unix(path) => {
            let listener = UnixListener::bind(&path)?;
            eprintln!("Listening on unix://{}", path.display());
            loop {
                let (stream, _) = listener.accept().await?;
                spawn_connection(stream, backend.clone());
            }
        }
    }
    Ok(())
}

/// spawn_connection serves a client against the backend and reports transport failures.
fn spawn_connection(
    stream: impl AsyncRead + AsyncWrite + Unpin + Send + 'static,
    backend: Backend,
) {
    tokio::spawn(async move {
        if let Err(error) = serve_connection(stream, backend).await {
            eprintln!("gRPC connection failed: {error}");
        }
    });
}

/// serve_connection dispatches a client's RPCs against the initialized backend.
async fn serve_connection(
    stream: impl AsyncRead + AsyncWrite + Unpin + Send + 'static,
    backend: Backend,
) -> Result<(), hyper::Error> {
    let service = service_fn(move |request: hyper::Request<Incoming>| {
        let command = Command {
            backend: backend.clone(),
            path: request.uri().path().to_owned(),
        };
        async move {
            let response = tonic::server::Grpc::new(ProtobufCodec)
                .max_decoding_message_size(i32::MAX as usize)
                .max_encoding_message_size(i32::MAX as usize)
                .unary(command, request)
                .await;
            Ok::<_, Infallible>(response)
        }
    });
    http2::Builder::new(TokioExecutor::new())
        .serve_connection(TokioIo::new(stream), service)
        .await
}
