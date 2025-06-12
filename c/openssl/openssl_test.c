#include <openssl/evp.h>
#include <stdio.h>
#include <string.h>

int main(int argc, char *argv[]) {
    EVP_MD_CTX *mdctx;
    const EVP_MD *md;
    char message[] = "Hello World\n";
    unsigned char md_value[EVP_MAX_MD_SIZE];
    unsigned int md_len, i;
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
    EVP_DigestFinal_ex(mdctx, md_value, &md_len);
    EVP_MD_CTX_destroy(mdctx);

    fprintf(stderr, "Digest is:\n");
    for (i = 0; i < md_len; i++) {
        fprintf(stderr, "%02x", md_value[i]);
    }
    fprintf(stderr, "\n");

    BIO *bio, *b64;

    b64 = BIO_new(BIO_f_base64());
    bio = BIO_new_fp(stderr, BIO_NOCLOSE);
    BIO_push(b64, bio);
    BIO_write(b64, md_value, md_len);
    BIO_flush(b64);

    BIO_free_all(b64);

    EVP_cleanup();
    return 0;
}
