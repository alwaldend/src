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
    {.name = "output",
     .key = 'o',
     .arg = "Path",
     .doc = "Output file path (default: stdout)"},
    {0},
};

struct arguments {
    char *digest;
    char *file;
    char *output_file;
};

static error_t argp_parser(int key, char *arg, struct argp_state *state) {
    struct arguments *arguments = state->input;
    switch (key) {
        case 'd':
            arguments->digest = arg;
            break;
        case 'f':
            arguments->file = arg;
            break;
        case 'o':
            arguments->output_file = arg;
            break;
        case ARGP_KEY_ARG:
            if (state->arg_num != 0) {
                fprintf(stderr, "arguments are not supported\n");
                return 1;
            }
        case ARGP_KEY_END:
            if (arguments->digest == NULL) {
                fprintf(stderr, "missing digest\n");
                return 1;
            }
            if (arguments->file == NULL) {
                fprintf(stderr, "missing file\n");
                return 1;
            }
            break;
    }
    return 0;
}

static struct argp argp = {
    .options = options,
    .parser = argp_parser,
    .doc =
        "Generate sri of a file \n\
\n\
Example:\n\
    bazel run //c/sri:bin -- --digest sha256 --file ${PWD}/README.md\n\
",
};

int cmd_run(int argc, char **argv) {
    FILE *input_file, *output_file;
    struct arguments arguments = {
        .digest = NULL,
        .file = NULL,
        .output_file = NULL,
    };
    if (0 != argp_parse(&argp, argc, argv, 0, 0, &arguments)) {
        fprintf(stderr, "could not parse arguments\n");
        return 1;
    };
    input_file = fopen(arguments.file, "r");
    if (input_file == NULL) {
        fprintf(stderr, "could not open input file %s: %s\n", arguments.file,
                strerror(errno));
        return 1;
    }
    if (arguments.output_file == NULL) {
        output_file = stdout;
    } else {
        output_file = fopen(arguments.output_file, "w");
        if (output_file == NULL) {
            fprintf(stderr, "could not open output file %s: %s\n",
                    arguments.output_file, strerror(errno));
            return 1;
        }
    }
    if (!sri_write(arguments.digest, input_file, output_file)) {
        fclose(input_file);
        errors_print(stderr);
        return 1;
    }
    fclose(input_file);
    if (arguments.output_file != NULL) {
        fclose(output_file);
    }
    return 0;
}
