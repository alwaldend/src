//! Single-thread platform I/O for the shared Rust CLI server.
//! These imports supply host byte streams that WASI Preview 1 cannot open.
//! Rust owns Anki request/response allocations and gRPC dispatch behind these streams.

use super::{cli::Transport, Dispatcher};
use std::{
    cell::RefCell,
    future::{poll_fn, Future},
    io,
    pin::{pin, Pin},
    sync::Arc,
    task::{Context, Poll, Wake, Waker},
};
use tokio::io::{AsyncRead, AsyncWrite, ReadBuf};

/// Task is one pinned protocol future retained until completion.
type Task = Pin<Box<dyn Future<Output = ()>>>;

#[link(wasm_import_module = "anki_io")]
extern "C" {
    fn open(kind: u32, pointer: u32, length: u32) -> i32;
    fn accept(handle: u32) -> i32;
    fn read(handle: u32, pointer: u32, length: u32) -> i32;
    fn write(handle: u32, pointer: u32, length: u32) -> i32;
    fn flush(handle: u32) -> i32;
    fn close(handle: u32);
    fn notify();
    fn wait();
}

thread_local! {
    /// TASKS queues Hyper tasks without entering Anki's nested futures executor.
    static TASKS: RefCell<Vec<Task>> = const { RefCell::new(Vec::new()) };
}

/// result decodes host progress without spinning while I/O is pending.
fn result(value: i32) -> Poll<io::Result<usize>> {
    match value {
        -2 => Poll::Pending,
        value if value < 0 => Poll::Ready(Err(io::Error::other("host stream operation failed"))),
        value => Poll::Ready(Ok(value as usize)),
    }
}

/// Executor queues protocol tasks without creating WASI threads.
#[derive(Clone)]
pub(super) struct Executor;

impl Executor {
    /// new selects the current command's executor.
    pub(super) fn new() -> Self {
        Self
    }
}

impl<F: Future<Output = ()> + 'static> hyper::rt::Executor<F> for Executor {
    fn execute(&self, future: F) {
        TASKS.with(|tasks| tasks.borrow_mut().push(Box::pin(future)));
        unsafe { notify() };
    }
}

/// Stream owns one host byte stream until its connection finishes.
struct Stream(u32);

impl Drop for Stream {
    fn drop(&mut self) {
        unsafe { close(self.0) };
    }
}

impl AsyncRead for Stream {
    fn poll_read(
        self: Pin<&mut Self>,
        _context: &mut Context<'_>,
        buffer: &mut ReadBuf<'_>,
    ) -> Poll<io::Result<()>> {
        let available = buffer.initialize_unfilled();
        match result(unsafe {
            read(
                self.0,
                available.as_mut_ptr() as u32,
                available.len() as u32,
            )
        }) {
            Poll::Ready(Ok(length)) => {
                buffer.advance(length);
                Poll::Ready(Ok(()))
            }
            Poll::Ready(Err(error)) => Poll::Ready(Err(error)),
            Poll::Pending => Poll::Pending,
        }
    }
}

impl AsyncWrite for Stream {
    fn poll_write(
        self: Pin<&mut Self>,
        _context: &mut Context<'_>,
        bytes: &[u8],
    ) -> Poll<io::Result<usize>> {
        result(unsafe { write(self.0, bytes.as_ptr() as u32, bytes.len() as u32) })
    }
    fn poll_flush(self: Pin<&mut Self>, _context: &mut Context<'_>) -> Poll<io::Result<()>> {
        result(unsafe { flush(self.0) }).map(|value| value.map(|_| ()))
    }
    fn poll_shutdown(self: Pin<&mut Self>, context: &mut Context<'_>) -> Poll<io::Result<()>> {
        self.poll_flush(context)
    }
}

/// Listener owns a host listener and yields independently managed streams.
struct Listener(u32);
impl Drop for Listener {
    fn drop(&mut self) {
        unsafe { close(self.0) };
    }
}

/// HostWake uses the same retained notification as host stream readiness.
struct HostWake;
impl Wake for HostWake {
    fn wake(self: Arc<Self>) {
        self.wake_by_ref();
    }
    fn wake_by_ref(self: &Arc<Self>) {
        unsafe { notify() };
    }
}

/// serve polls protocol tasks and blocks in the host only when all tasks are stalled.
/// It deliberately avoids futures::executor so Anki may use its own block_on.
pub(super) fn serve(
    transport: Transport,
    backend: Arc<dyn Dispatcher>,
) -> Result<(), Box<dyn std::error::Error>> {
    let mut tasks = Vec::new();
    let mut server = pin!(serve_transport(transport, backend));
    let waker = Waker::from(Arc::new(HostWake));
    let mut context = Context::from_waker(&waker);
    loop {
        TASKS.with(|pending| tasks.extend(pending.take()));
        // Host readiness is coalesced across streams, so poll every task with one
        // shared waker. This small worker needs no per-stream waiter registry.
        tasks.retain_mut(|task| task.as_mut().poll(&mut context).is_pending());
        if let Poll::Ready(outcome) = server.as_mut().poll(&mut context) {
            return outcome;
        }
        unsafe { wait() };
    }
}

/// serve_transport initializes the validated transport and preserves shared backend state.
async fn serve_transport(
    transport: Transport,
    backend: Arc<dyn Dispatcher>,
) -> Result<(), Box<dyn std::error::Error>> {
    let (kind, address) = match transport {
        Transport::Stdio => (0, Vec::new()),
        Transport::Tcp(address) => (1, address.into_bytes()),
        Transport::Unix(path) => (2, path.as_os_str().as_encoded_bytes().to_vec()),
    };
    let handle = unsafe { open(kind, address.as_ptr() as u32, address.len() as u32) };
    if handle < 0 {
        return Err(io::Error::other("open host transport").into());
    }
    if kind == 0 {
        super::serve_connection(Stream(handle as u32), backend).await?;
    } else {
        let listener = Listener(handle as u32);
        loop {
            let handle = poll_fn(|_context| result(unsafe { accept(listener.0) })).await?;
            super::spawn_connection(Stream(handle as u32), backend.clone());
        }
    }
    Ok(())
}
