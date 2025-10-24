#include <stdbool.h>
#include <stdio.h>
#include <string.h>

_Thread_local static char _ERRORS[100][BUFSIZ];
_Thread_local static unsigned int _ERRORS_INDEX = 0;

/*
 * Append an error message to thread local error messages
 *
 * Args:
 *     err: error message
 * */
void errors_append(const char* err) {
    strncpy(_ERRORS[_ERRORS_INDEX], err, BUFSIZ);
    _ERRORS_INDEX += 1;
}

/*
 * Print thread local errors to stream
 *
 * Args:
 *     stream: stream to write errors to
 * */
void errors_print(FILE* stream) {
    bool start = true;
    for (int i = _ERRORS_INDEX - 1; i >= 0; i--) {
        if (start) {
            start = false;
        } else {
            fprintf(stream, ": ");
        }
        fprintf(stream, "%s", _ERRORS[i]);
    }
    fprintf(stream, "\n");
}
