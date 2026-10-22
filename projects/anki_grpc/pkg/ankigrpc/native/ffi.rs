//! C ABI adapters for direct Anki backend calls and shared CLI startup.

use std::ffi::{c_char, c_int, CStr};

#[path = "../../../internal/backend/ffi.rs"]
mod backend_ffi;

/// anki_grpc_run copies a C argument vector and returns the native CLI exit status.
///
/// # Safety
/// When argc is positive, argv must point to argc readable pointers to terminated C strings.
/// The argument vector and strings must remain valid for the duration of this call.
#[no_mangle]
pub unsafe extern "C" fn anki_grpc_run(argc: c_int, argv: *const *const c_char) -> c_int {
    if argc <= 0 || argv.is_null() {
        eprintln!("Anki CLI requires an argument vector including its program name");
        return 2;
    }
    let pointers = unsafe { std::slice::from_raw_parts(argv, argc as usize) };
    let mut arguments = Vec::with_capacity(pointers.len());
    for &pointer in pointers {
        if pointer.is_null() {
            eprintln!("Anki CLI argument must not be null");
            return 2;
        }
        let value = unsafe { CStr::from_ptr(pointer) };
        let Ok(value) = value.to_str() else {
            eprintln!("Anki CLI arguments must be UTF-8");
            return 2;
        };
        arguments.push(value.to_owned());
    }
    anki_grpc_server::run(arguments)
}
