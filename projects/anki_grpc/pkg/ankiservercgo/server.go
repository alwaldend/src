package ankiservercgo

/*
#include <stdlib.h>
int anki_grpc_run(int argc, const char *const *argv);
*/
import "C"

import (
	"fmt"
	"strings"
	"unsafe"
)

// Run invokes the native CLI with a full argument vector including its program name.
// Its serve subcommand blocks until the server exits; Run reports nonzero native exit statuses.
// Arguments are copied into native-owned strings and retained only for this call.
func Run(arguments []string) error {
	if len(arguments) == 0 || len(arguments) > 1<<31-1 {
		return fmt.Errorf("Anki CLI requires a supported argument count including its program name")
	}
	argv := make([]*C.char, len(arguments))
	for index, argument := range arguments {
		if strings.ContainsRune(argument, 0) {
			return fmt.Errorf("Anki CLI argument %d contains a null byte", index)
		}
		argv[index] = C.CString(argument)
		defer C.free(unsafe.Pointer(argv[index]))
	}
	status := C.anki_grpc_run(C.int(len(argv)), (**C.char)(unsafe.Pointer(&argv[0])))
	if status != 0 {
		return fmt.Errorf("Anki CLI exited with status %d", int(status))
	}
	return nil
}
