//! Browser reactor exposing the shared backend ABI and scalar buffer accessors.

#[path = "../../../internal/backend/ffi.rs"]
mod backend_ffi;

use backend_ffi::{anki_backend_free_buffer, AnkiBuffer};

/// Allocate a zeroed byte buffer and its correctly aligned, Rust-owned descriptor.
/// Use a zero-length buffer as an empty output slot for the shared backend ABI.
#[no_mangle]
pub extern "C" fn anki_browser_buffer_new(len: usize) -> *mut AnkiBuffer {
    Box::into_raw(Box::new(AnkiBuffer::new(vec![0; len])))
}

/// Return the buffer's byte address without exposing a struct-by-value ABI to JavaScript.
///
/// # Safety
/// buffer must be null or a live descriptor returned by anki_browser_buffer_new.
#[no_mangle]
pub unsafe extern "C" fn anki_browser_buffer_data(buffer: *const AnkiBuffer) -> *mut u8 {
    if buffer.is_null() {
        return std::ptr::null_mut();
    }
    unsafe { (*buffer).data }
}

/// Return the buffer's current byte length, including results written by backend commands.
///
/// # Safety
/// buffer must be null or a live descriptor returned by anki_browser_buffer_new.
#[no_mangle]
pub unsafe extern "C" fn anki_browser_buffer_len(buffer: *const AnkiBuffer) -> usize {
    if buffer.is_null() {
        return 0;
    }
    unsafe { (*buffer).len }
}

/// Release both the byte allocation and its descriptor through the shared ownership ABI.
///
/// # Safety
/// buffer must be null or an unreleased descriptor returned by anki_browser_buffer_new.
/// No call may access the descriptor or its bytes concurrently with or after release.
#[no_mangle]
pub unsafe extern "C" fn anki_browser_buffer_free(buffer: *mut AnkiBuffer) {
    if !buffer.is_null() {
        let buffer = unsafe { Box::from_raw(buffer) };
        unsafe { anki_backend_free_buffer(*buffer) };
    }
}
