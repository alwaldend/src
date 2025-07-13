#include <argp.h>
#include <err.h>
#include <errno.h>
#include <stdio.h>
#include <string.h>

#include "c/sri/errors.h"
#include "c/sri/sri.h"

static struct argp_option options[] = {
    {.name = "digest",
     .key = 'd',
     .arg = "String",
     .doc = "Digest type (sha256, for example)"},
    {.name = "file",
     .key = 'f',
     .arg = "Path",
     .doc = "Path to the file to parse"},
};

struct arguments {
    char *digest;
    char *file;
};

static error_t parse_opt(int key, char *arg, struct argp_state *state) {
    struct arguments *arguments = state->input;
    switch (key) {
        case 'd':
            arguments->digest = arg;
            break;
        case 'f':
            arguments->file = arg;
            break;
        case ARGP_KEY_END:
            if (arguments->digest == NULL) {
                fprintf(stderr, "missing digest\n");
                return ARGP_ERR_UNKNOWN;
            }
            if (arguments->file == NULL) {
                fprintf(stderr, "missing file\n");
                return ARGP_ERR_UNKNOWN;
            }
            break;
        default:
            return ARGP_ERR_UNKNOWN;
    }
    return 0;
}

static struct argp argp = {
    .options = options,
    .parser = parse_opt,
    .doc =
        "Generate sri of a file \n\
\n\
Example:\n\
    bazel run //c/sri:bin -- --digest sha256 --file ${PWD}/README.md\n\
",
};

int cmd_run(int argc, char **argv) {
    FILE *file;
    struct arguments arguments = {.digest = NULL, .file = NULL};
    if (!argp_parse(&argp, argc, argv, 0, 0, &arguments)) {
        fprintf(stderr, "could not parse arguments\n");
        return 1;
    };
    file = fopen(arguments.file, "r");
    if (file == NULL) {
        fprintf(stderr, "could not open file %s: %s", arguments.file,
                strerror(errno));
        return 1;
    }
    if (!sri_write(arguments.digest, file, stdout)) {
        fclose(file);
        errors_print(stderr);
        return 1;
    }
    fclose(file);
    return 0;
}
