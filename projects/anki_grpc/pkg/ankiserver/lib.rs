//! Shared Rust CLI, opaque protobuf gRPC server and Anki backend dispatch.
//! Platform I/O stays below this server; callers own worker process lifetime.

mod cli;
use anki_backend_dispatch::{command, init_backend, Backend, CommandError};
use bytes::{Buf, BufMut, Bytes};
use hyper::{body::Incoming, rt::Executor as _, server::conn::http2, service::service_fn};
use hyper_util::rt::TokioIo;
use std::{convert::Infallible, future::Future, pin::Pin, sync::Arc};
use tokio::io::{AsyncRead, AsyncWrite};
use tonic::{codec, server::UnaryService, Request, Response, Status};

#[cfg(not(target_os = "wasi"))]
use hyper_util::rt::TokioExecutor as Executor;
#[cfg(not(target_os = "wasi"))]
use tokio::net::{TcpListener, UnixListener};
#[cfg(target_os = "wasi")]
mod wasi;
#[cfg(target_os = "wasi")]
use wasi::Executor;

/// Dispatcher supplies opaque protobuf calls to the shared CLI server.
pub trait Dispatcher: Send + Sync {
    /// dispatch executes one method and preserves its gRPC status and trailers.
    fn dispatch(&self, path: &str, request: &[u8]) -> Result<Vec<u8>, Status>;
}

/// BackendInitializer constructs the one dispatcher shared by the server's clients.
type BackendInitializer = fn() -> Result<Arc<dyn Dispatcher>, Box<dyn std::error::Error>>;

/// AnkiDispatcher owns the single upstream backend shared by all connections.
struct AnkiDispatcher {
    /// backend owns Anki collection state for this server process.
    backend: Backend,
}

impl Dispatcher for AnkiDispatcher {
    fn dispatch(&self, path: &str, request: &[u8]) -> Result<Vec<u8>, Status> {
        command(&self.backend, path, request).map_err(|error| match error {
            CommandError::Unknown => Status::unimplemented(format!("unknown Anki command: {path}")),
            CommandError::Native(details) => {
                let mut metadata = tonic::metadata::MetadataMap::new();
                metadata.insert_bin(
                    "anki-error-bin",
                    tonic::metadata::MetadataValue::from_bytes(&details),
                );
                Status::with_metadata(tonic::Code::Unknown, "Anki command failed", metadata)
            }
        })
    }
}

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

/// Command dispatches one request against the process-owned backend.
struct Command {
    /// backend shares collection state between clients and RPCs.
    backend: Arc<dyn Dispatcher>,
    /// path identifies the upstream service and method.
    path: String,
}

impl UnaryService<Bytes> for Command {
    type Response = Bytes;
    type Future = Pin<Box<dyn Future<Output = Result<Response<Bytes>, Status>> + Send>>;
    fn call(&mut self, request: Request<Bytes>) -> Self::Future {
        let backend = self.backend.clone();
        let path = self.path.clone();
        Box::pin(async move {
            #[cfg(not(target_os = "wasi"))]
            let result =
                tokio::task::spawn_blocking(move || backend.dispatch(&path, &request.into_inner()))
                    .await
                    .expect("Anki command task");
            #[cfg(target_os = "wasi")]
            let result = backend.dispatch(&path, &request.into_inner());
            result.map(|bytes| Response::new(Bytes::from(bytes)))
        })
    }
}

/// run executes the complete CLI using Anki's default backend initialization.
pub fn run(arguments: Vec<String>) -> i32 {
    run_with_backend(arguments, || {
        let backend = init_backend(&[]).map_err(std::io::Error::other)?;
        Ok(Arc::new(AnkiDispatcher { backend }))
    })
}

/// run_with_backend serves a caller-supplied dispatcher through the same CLI and transports.
/// Initialization occurs only after argument validation; every connection shares its result.
pub fn run_with_backend(arguments: Vec<String>, initialize: BackendInitializer) -> i32 {
    cli::run(arguments, |transport| serve(transport, initialize))
}

/// serve retains synchronous backend ownership until its serving runtime has stopped.
#[cfg(not(target_os = "wasi"))]
fn serve(
    transport: cli::Transport,
    initialize: BackendInitializer,
) -> Result<(), Box<dyn std::error::Error>> {
    let runtime = tokio::runtime::Builder::new_multi_thread()
        .enable_all()
        .build()?;
    let backend = runtime.block_on(async { initialize() })?;
    let result = runtime.block_on(serve_native(transport, backend.clone()));
    // Finish every serving task before releasing a backend that can own its own runtime.
    drop(runtime);
    drop(backend);
    result
}

/// serve_native serves all native transports while its caller retains the backend.
#[cfg(not(target_os = "wasi"))]
async fn serve_native(
    transport: cli::Transport,
    backend: Arc<dyn Dispatcher>,
) -> Result<(), Box<dyn std::error::Error>> {
    match transport {
        cli::Transport::Stdio => {
            serve_connection(
                tokio::io::join(tokio::io::stdin(), tokio::io::stdout()),
                backend,
            )
            .await?
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

/// serve executes the same server on WASI's single-thread event loop.
#[cfg(target_os = "wasi")]
fn serve(
    transport: cli::Transport,
    initialize: BackendInitializer,
) -> Result<(), Box<dyn std::error::Error>> {
    wasi::serve(transport, initialize()?)
}

/// spawn_connection serves one client and reports its transport failure.
fn spawn_connection(
    stream: impl AsyncRead + AsyncWrite + Unpin + Send + 'static,
    backend: Arc<dyn Dispatcher>,
) {
    Executor::new().execute(async move {
        if let Err(error) = serve_connection(stream, backend).await {
            eprintln!("gRPC connection failed: {error}");
        }
    });
}

/// serve_connection shares gRPC framing and dispatch between native and WASI builds.
async fn serve_connection(
    stream: impl AsyncRead + AsyncWrite + Unpin + Send + 'static,
    backend: Arc<dyn Dispatcher>,
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
    http2::Builder::new(Executor::new())
        .serve_connection(TokioIo::new(stream), service)
        .await
}
