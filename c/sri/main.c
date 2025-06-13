#include <ctype.h>
#include <err.h>
#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#include "c/sri/errors.h"
#include "c/sri/sri.h"

int main(int argc, char **argv) {
    char *file_path = NULL;
    char *digest_type = NULL;
    char flag;
    FILE *file;

    while ((flag = getopt(argc, argv, "d:f:")) != -1) {
        switch (flag) {
            case 'd':
                digest_type = optarg;
                break;
            case 'f':
                file_path = optarg;
                break;
            case '?':
                if (isprint(optopt)) {
                    fprintf(stderr, "Unknown option '-%c'.\n", optopt);
                } else {
                    fprintf(stderr, "Unknown option character '\\x%x'.\n",
                            optopt);
                }
                return 1;
            default:
                abort();
        }
    }
    file = fopen(file_path, "r");
    if (file == NULL) {
        fprintf(stderr, "could not open file %s: %s", file_path,
                strerror(errno));
        return 1;
    }
    if (!sri_write(digest_type, file, stdout)) {
        fclose(file);
        errors_print(stderr);
        return 1;
    }
    fclose(file);
    return 0;
}
