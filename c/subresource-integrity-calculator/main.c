#include <ctype.h>
#include <getopt.h>
#include <openssl/sha.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>

struct Config {
    char* output_file;
};

int main(int argc, char** argv) {
    int option_character;
    while ((option_character = getopt(argc, argv, "o:")) != -1) {
        switch (option_character) {
            case 'o':
                break;
            case '?':
                if (optopt == 'c') {
                    fprintf(stderr, "Option -%c requires an argument.\n",
                            optopt);
                } else if (isprint(optopt)) {
                    fprintf(stderr, "Unknown option '-%c'.\n", optopt);
                } else {
                    fprintf(stderr, "Unknown option character `\\x%x'.\n",
                            optopt);
                }
                return 1;
            default:
                abort();
        }
    }
    return 0;
}
