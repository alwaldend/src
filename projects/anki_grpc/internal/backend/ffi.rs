//! Shared C ABI for owned Anki backend handles and command buffers.

use std::ffi::c_int;
use std::ptr;

use anki_backend_dispatch::{command, init_backend, Backend, CommandError};

/// AnkiBuffer owns a Rust byte allocation until anki_backend_free_buffer releases it.
#[repr(C)]
pub struct AnkiBuffer {
    /// data points to len bytes allocated as a boxed slice.
    pub(crate) data: *mut u8,
    /// len records the complete allocation length.
    pub(crate) len: usize,
}

impl AnkiBuffer {
    /// new transfers the byte allocation into the C ABI.
    pub(crate) fn new(bytes: Vec<u8>) -> Self {
        let bytes = bytes.into_boxed_slice();
        let len = bytes.len();
        Self {
            data: Box::into_raw(bytes) as *mut u8,
            len,
        }
    }
}

/// anki_backend_new initializes a Rust-owned backend or returns a UTF-8 initialization error.
///
/// # Safety
/// error must point to writable storage for one AnkiBuffer.
#[no_mangle]
pub unsafe extern "C" fn anki_backend_new(error: *mut AnkiBuffer) -> *mut Backend {
    if error.is_null() {
        return ptr::null_mut();
    }
    match init_backend(&[]) {
        Ok(backend) => Box::into_raw(Box::new(backend)),
        Err(message) => {
            unsafe {
                *error = AnkiBuffer::new(message.into_bytes());
            }
            ptr::null_mut()
        }
    }
}

/// anki_backend_command returns 0 for success, 1 for native error, 2 for unknown route, or 3 for invalid input.
///
/// # Safety
/// backend must remain live for the call. Nonempty byte ranges must be readable;
/// output must point to writable storage for one AnkiBuffer. Borrowed bytes are not retained.
#[no_mangle]
pub unsafe extern "C" fn anki_backend_command(
    backend: *const Backend,
    path: *const u8,
    path_len: usize,
    input: *const u8,
    input_len: usize,
    output: *mut AnkiBuffer,
) -> c_int {
    if backend.is_null()
        || output.is_null()
        || path.is_null()
        || (input.is_null() && input_len != 0)
    {
        return 3;
    }
    let path = unsafe { std::slice::from_raw_parts(path, path_len) };
    let Ok(path) = std::str::from_utf8(path) else {
        return 3;
    };
    let input = if input_len == 0 {
        &[]
    } else {
        unsafe { std::slice::from_raw_parts(input, input_len) }
    };
    match command(unsafe { &*backend }, path, input) {
        Ok(bytes) => {
            unsafe {
                *output = AnkiBuffer::new(bytes);
            }
            0
        }
        Err(CommandError::Native(bytes)) => {
            unsafe {
                *output = AnkiBuffer::new(bytes);
            }
            1
        }
        Err(CommandError::Unknown) => 2,
    }
}

/// anki_backend_free_buffer returns an owned result allocation to Rust.
///
/// # Safety
/// buffer must come from this ABI and must be released exactly once.
#[no_mangle]
pub unsafe extern "C" fn anki_backend_free_buffer(buffer: AnkiBuffer) {
    if !buffer.data.is_null() {
        unsafe {
            drop(Box::from_raw(ptr::slice_from_raw_parts_mut(
                buffer.data,
                buffer.len,
            )));
        }
    }
}

/// anki_backend_free releases backend state after all native calls have finished.
///
/// # Safety
/// backend must be null or an unreleased handle returned by anki_backend_new.
/// No call may use the handle concurrently with or after its release.
#[no_mangle]
pub unsafe extern "C" fn anki_backend_free(backend: *mut Backend) {
    if !backend.is_null() {
        unsafe {
            drop(Box::from_raw(backend));
        }
    }
}
