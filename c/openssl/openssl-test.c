#include <openssl/sha.h>
#include <stdio.h>

int main(int argc, char **argv) {
    unsigned char buffer[BUFSIZ];
    FILE *f;
    SHA_CTX ctx;
    size_t len;
    if (argc < 2) {
        fprintf(stderr, "usage: %s <file>\n", argv[0]);
        return 1;
    }

    f = fopen(argv[1], "r");
    if (!f) {
        fprintf(stderr, "couldn't open %s\n", argv[1]);
        return 1;
    }

    SHA1_Init(&ctx);

    do {
        len = fread(buffer, 1, BUFSIZ, f);
        SHA1_Update(&ctx, buffer, len);
    } while (len == BUFSIZ);

    SHA1_Final(buffer, &ctx);

    fclose(f);

    for (len = 0; len < SHA_DIGEST_LENGTH; ++len) {
        printf("%02x", buffer[len]);
    }
    putchar('\n');
    return 0;
}
