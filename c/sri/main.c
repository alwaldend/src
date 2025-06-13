#include <err.h>
#include <stdio.h>

#include "c/sri/errors.h"
#include "c/sri/sri.h"

int main(int argc, char** argv) {
    if (argc != 3) {
        fprintf(stderr, "invalid argument count\n");
        return 1;
    }
    if (!sri_write(argv[1], stdout, argv[2])) {
        errors_print(stderr);
        return 1;
    }
    return 0;
}
