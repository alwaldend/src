#include <stdio.h>
#include <string.h>

#include "c/subresource-integrity-calculator/calculator.h"

int main(int argc, char** argv) {
    unsigned char out[BUFSIZ];
    unsigned int out_length;

    if (argc != 3) {
        fprintf(stderr, "invalid argument count\n");
        return 1;
    }

    if (calculate_sri(argv[1], argv[2], strlen(argv[2]), out, &out_length) !=
        0) {
        fprintf(stderr, "could not calculate sri\n");
        return 1;
    }

    fprintf(stderr, "output:\n");
    printf("%s\n", out);
    return 1;
}
