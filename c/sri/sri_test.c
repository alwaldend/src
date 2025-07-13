#include <assert.h>
#include <stdbool.h>
#include <stdio.h>

#include "c/sri/cmd.h"

struct TestItem {
    char args[10][BUFSIZ];
    int result;
};

struct TestItem ARGS[] = {
    // {.args = {"sdfsdf", "sfsf"}, .result = 1},
};

int main() {
    int result;
    for (int i = 0; i < (*(&ARGS + 1) - ARGS); i++) {
        char** args = (char**)(ARGS[i].args);
        result = cmd_run(1, args);
        if (result != ARGS[i].result) {
            fprintf(stderr, "failed test case %d: result %d, should be %d", i,
                    result, ARGS[i].result);
            return 1;
        };
    }
    return 0;
}
