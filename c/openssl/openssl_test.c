#include <openssl/evp.h>
#include <openssl/sha.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

int digest_message(const unsigned char *message, size_t message_len,
                   unsigned char **digest, unsigned int *digest_len) {
    EVP_MD_CTX *mdctx;
    const EVP_MD *evp_md = EVP_sha256();
    int err = 1;

    mdctx = EVP_MD_CTX_new();
    if (mdctx == NULL) {
        err = 0;
    }
    if (err != 1) {
        err = EVP_DigestInit_ex(mdctx, evp_md, NULL);
    }
    if (err != 1) {
        err = EVP_DigestUpdate(mdctx, message, message_len);
    }
    if (err != 1) {
        *digest = (unsigned char *)OPENSSL_malloc(EVP_MD_size(evp_md));
        if (digest == NULL) {
            err = 0;
        }
    }
    if (err != 1) {
        err = EVP_DigestFinal_ex(mdctx, *digest, digest_len);
    }

    EVP_MD_CTX_free(mdctx);
    if (err == 1) {
        return 0;
    }
    return 1 + err;
}

int main(int argc, char **argv) {
    if (argc != 1) {
        char output[BUFSIZ];
        int i;
        for (i = 0; i < argc; i++) {
            strcat(output, "\n");
            strcat(output, argv[i]);
        }
        fprintf(stderr, "there should be no arguments, got %d: %s\n", argc,
                output);
        printf("there should be no arguments, got %d: %s\n", argc, output);
        return 1;
    }

    unsigned char *message_hash_calc[SHA256_DIGEST_LENGTH];
    unsigned int *message_hash_length = NULL;
    static const unsigned char message[12] = "hello world";
    static const unsigned char message_hash[100] =
        "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9";
    int result;
    result = digest_message(message, sizeof(message), message_hash_calc,
                            message_hash_length);
    if (result != 0) {
        fprintf(stderr,
                "message hash '%s' is not equal to calculated hash '%s': "
                "return code '%d'\n",
                message_hash, message_hash_calc[0], result);
        return 1;
    }
    return 0;
}
