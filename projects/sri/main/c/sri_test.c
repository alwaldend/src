#include <assert.h>
#include <errno.h>
#include <stdbool.h>
#include <stdio.h>
#include <string.h>

#include "projects/sri/main/c/cmd.h"

struct TestItem {
    char** argv;
    int argc;
    char* test_string;
    char* test_string_digest;
    int result;
};

struct TestItem ARGS[] = {
    {
        .argv = (char*[]){"sri", "sfsf"},
        .argc = 2,
        .result = 1,
    },
    {
        .argv = (char*[]){"sri", "--digest", "sha256", "--file",
                          "test-file.txt", "--output", "test-output.txt"},
        .argc = 7,
        .test_string = "Hello World",
        .test_string_digest =
            "sha256-pZGm1Av0IEBKARczz7exkNYsZb8LzaMrV7J32a2fFG4=\n",
        .result = 0,
    },
    {
        .argv =
            (char*[]){"sri", "--digest", "sha256", "--file", "test-file.txt"},
        .argc = 5,
        .test_string = "sdfsdfsfsfsdf",
        .result = 0,
    },
};

int main() {
    int result;
    for (int i = 0; i < (*(&ARGS + 1) - ARGS); i++) {
        struct TestItem args = ARGS[i];
        if (args.test_string != NULL) {
            FILE* test_file = fopen("test-file.txt", "w");
            if (test_file == NULL) {
                fprintf(stderr, "could not open file %s: %s\n", "test-file.txt",
                        strerror(errno));
                return 1;
            }
            fprintf(test_file, "%s", args.test_string);
            fclose(test_file);
        }
        result = cmd_run(args.argc, args.argv);
        if (result != args.result) {
            fprintf(stderr, "failed test case %d: result %d, should be %d\n", i,
                    result, ARGS[i].result);
            return 1;
        };
        if (args.test_string_digest != NULL) {
            FILE* test_output = fopen("test-output.txt", "r");
            if (test_output == NULL) {
                fprintf(stderr, "could not open output file %s: %s\n",
                        "test-output.txt", strerror(errno));
                return 1;
            }
            char test_output_string[BUFSIZ];
            fgets(test_output_string, BUFSIZ, test_output);
            fclose(test_output);
            if (0 != strcmp(args.test_string_digest, test_output_string)) {
                fprintf(stderr, "invalid output: should be '%s', got '%s'\n",
                        args.test_string_digest, test_output_string);
                return 1;
            }
        }
    }
    return 0;
}
