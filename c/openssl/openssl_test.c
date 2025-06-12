#include <openssl/evp.h>
#include <stdio.h>
#include <string.h>

int main(int argc, char *argv[]) {
    EVP_MD_CTX *mdctx;
    const EVP_MD *md;
    char message[] = "Hello World\n";
    char message_digest_base64[] =
        "0qhPS4tlCTfsj3PNi+LHSt1akRumTfJ0WO2CKdqASiY=";
    unsigned char calc_digest_value[EVP_MAX_MD_SIZE];
    unsigned int calc_digest_len, i;
    char calc_digest_base64[EVP_MAX_MD_SIZE];
    if (argc != 1) {
        fprintf(stderr, "did not expect args: ");
        for (i = 0; i < (unsigned int)argc; i++) {
            fprintf(stderr, "%s", argv[i]);
        }
        fprintf(stderr, "\n");
        return 1;
    }

    OpenSSL_add_all_digests();

    md = EVP_get_digestbyname("sha256");

    if (!md) {
        fprintf(stderr, "Unknown message digest\n");
        return 1;
    }

    mdctx = EVP_MD_CTX_create();
    EVP_DigestInit_ex(mdctx, md, NULL);
    EVP_DigestUpdate(mdctx, message, strlen(message));
    EVP_DigestFinal_ex(mdctx, calc_digest_value, &calc_digest_len);
    EVP_MD_CTX_destroy(mdctx);

    fprintf(stderr, "Digest is:\n");
    for (i = 0; i < calc_digest_len; i++) {
        fprintf(stderr, "%02x", calc_digest_value[i]);
    }
    fprintf(stderr, "\n");

    BIO *bio, *b64;
    FILE *md_value_base64_file =
        fmemopen(calc_digest_base64, EVP_MAX_MD_SIZE, "w");

    b64 = BIO_new(BIO_f_base64());
    bio = BIO_new_fp(md_value_base64_file, BIO_CLOSE);
    BIO_push(b64, bio);
    BIO_write(b64, calc_digest_value, calc_digest_len);
    BIO_flush(b64);
    BIO_free_all(b64);

    fprintf(stderr, "base64 digest: %s\n", calc_digest_base64);
    if (!strncmp(message_digest_base64, calc_digest_base64, EVP_MAX_MD_SIZE)) {
        fprintf(stderr, "calculated base64 digest '%s' is different from '%s'",
                calc_digest_base64, message_digest_base64);
        return 1;
    }

    EVP_cleanup();
    return 0;
}
