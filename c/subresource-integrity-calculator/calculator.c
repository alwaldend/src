#include <openssl/err.h>
#include <openssl/evp.h>
#include <stdio.h>
#include <string.h>

/*
 * Calculate base64
 *
 * Args:
 *     in: input string
 *     in_length: length of the input string
 *     out: output string
 *     out_length: length of the output string
 *
 * Returns:
 *     0: if success
 *     1: if failure
 * */
int calculate_base64(unsigned char *in, unsigned int in_length,
                     unsigned char *out, unsigned int *out_length) {
    BIO *bio_fp, *bio_base64;
    *out_length = ((4 * in_length / 3) + 3) & ~3;
    FILE *in_file = fmemopen(&out, *out_length, "w");
    if (in_file == NULL) {
        fprintf(stderr, "could not open fmemopen\n");
        return 1;
    }
    bio_base64 = BIO_new(BIO_f_base64());
    if (bio_base64 == NULL) {
        fprintf(stderr, "could not create base64 BIO\n");
        return 1;
    }
    bio_fp = BIO_new_fp(in_file, BIO_CLOSE);
    if (bio_fp == NULL) {
        fprintf(stderr, "could not create fp BIO");
        return 1;
    }
    BIO_push(bio_base64, bio_fp);
    if (BIO_write(bio_base64, in, in_length) != 1) {
        fprintf(stderr, "could not write to base64 bio: %s\n",
                ERR_error_string(ERR_get_error(), NULL));
        return 1;
    };
    if (BIO_flush(bio_base64) != 1) {
        fprintf(stderr, "could not flush base64 bio: %s\n",
                ERR_error_string(ERR_get_error(), NULL));
        return 1;
    };
    BIO_free_all(bio_base64);
    return 0;
}

/*
 * Calculate openssl digest
 *
 * Args:
 *     digest_type: openssl digest type (sha256, for example)
 *     in: input string
 *     in_length: input string length
 *     out: output string (openssl digest)
 *     out_length: length of the output string
 *
 * Returns:
 *     0: if success
 *     1: if failure
 * */
int calculate_digest(char *digest_type, char *in, unsigned int in_length,
                     unsigned char *out, unsigned int *out_length) {
    EVP_MD_CTX *evp_md_ctx;
    const EVP_MD *evp_md_st;

    OpenSSL_add_all_digests();

    evp_md_st = EVP_get_digestbyname(digest_type);
    if (evp_md_st == NULL) {
        fprintf(stderr, "unknown message digest: %s\n", digest_type);
        return 1;
    }

    evp_md_ctx = EVP_MD_CTX_new();
    if (evp_md_ctx == NULL) {
        fprintf(stderr, "could not create evp md ctx: %s\n",
                ERR_error_string(ERR_get_error(), NULL));
        return 1;
    }

    if (EVP_DigestInit_ex2(evp_md_ctx, evp_md_st, NULL) != 1) {
        fprintf(stderr, "message digest initialization failed: %s\n",
                ERR_error_string(ERR_get_error(), NULL));
        EVP_MD_CTX_free(evp_md_ctx);
        return 1;
    };

    if (EVP_DigestUpdate(evp_md_ctx, in, in_length) != 1) {
        fprintf(stderr, "message digest update failed: %s\n",
                ERR_error_string(ERR_get_error(), NULL));
        EVP_MD_CTX_free(evp_md_ctx);
        return 1;
    };

    if (EVP_DigestFinal_ex(evp_md_ctx, out, out_length) != 1) {
        fprintf(stderr, "message digest finalization failed: %s\n",
                ERR_error_string(ERR_get_error(), NULL));
        EVP_MD_CTX_free(evp_md_ctx);
        return 1;
    };
    EVP_MD_CTX_free(evp_md_ctx);
    return 0;
}

/*
 * Calculate sri
 *
 * Args:
 *     digest_type: openssl digest type (sha256, for example)
 *     in: input string
 *     in_length: input string length
 *     out: output string (openssl digest)
 *     out_length: length of the output string
 *
 * Returns:
 *     0: if success
 *     1: if failure
 * */
int calculate_sri(char *digest_type, char *in, unsigned int in_length,
                  unsigned char *out, unsigned int *out_length) {
    unsigned int digest_type_length = strlen(digest_type);
    unsigned char digest[BUFSIZ], digest_base64[BUFSIZ];
    unsigned int digest_length, digest_base64_length;

    if (calculate_digest(digest_type, in, in_length, digest, &digest_length) !=
        0) {
        fprintf(stderr, "could not calculate digest\n");
        return 1;
    };

    if (calculate_base64(digest, digest_length, digest_base64,
                         &digest_base64_length) != 0) {
        fprintf(stderr, "could not calculate base64\n");
        return 1;
    };
    strncat((char *)out, digest_type, digest_type_length);
    strncat((char *)out, "-", 1);
    strncat((char *)out, (char *)digest_base64, digest_base64_length);
    *out_length = digest_type_length + 1 + digest_base64_length;
    return 0;
}
