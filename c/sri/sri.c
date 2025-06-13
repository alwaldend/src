#include <errno.h>
#include <openssl/err.h>
#include <openssl/evp.h>
#include <stdbool.h>
#include <stdio.h>
#include <string.h>

#include "c/sri/errors.h"

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
 *     true if success, false if failure
 * */
bool sri_calculate_base64(const unsigned char *in, const unsigned int in_length,
                          unsigned char *out, unsigned int *out_length) {
    BIO *bio_fp, *bio_base64;
    char err[BUFSIZ];
    *out_length = ((4 * in_length / 3) + 3) & ~3;
    FILE *in_file = fmemopen(out, *out_length, "w");
    if (in_file == NULL) {
        sprintf(err, "could not open fmemopen: %s", strerror(errno));
        errors_append(err);
        return false;
    }
    bio_base64 = BIO_new(BIO_f_base64());
    if (bio_base64 == NULL) {
        sprintf(err, "could not create base64 BIO: %s", strerror(errno));
        errors_append(err);
        return false;
    }
    bio_fp = BIO_new_fp(in_file, BIO_CLOSE);
    if (bio_fp == NULL) {
        sprintf(err, "could not create fp BIO: %s", strerror(errno));
        errors_append(err);
        return false;
    }
    BIO_push(bio_base64, bio_fp);
    if (BIO_write(bio_base64, in, in_length) != 1) {
        sprintf(err, "could not write to base64 bio: %s", strerror(errno));
        errors_append(err);
        return false;
    };
    if (BIO_flush(bio_base64) != 1) {
        sprintf(err, "could not flush base64 bio: %s", strerror(errno));
        errors_append(err);
        return false;
    };
    BIO_free_all(bio_base64);
    return true;
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
 *     true if success, false if failure
 * */
bool sri_calculate_digest(char *digest_type, char *in, unsigned int in_length,
                          unsigned char *out, unsigned int *out_length) {
    EVP_MD_CTX *evp_md_ctx;
    const EVP_MD *evp_md_st;
    char err[BUFSIZ];

    OpenSSL_add_all_digests();

    evp_md_st = EVP_get_digestbyname(digest_type);
    if (evp_md_st == NULL) {
        sprintf(err, "unknown digest type: %s", digest_type);
        errors_append(err);
        return false;
    }

    evp_md_ctx = EVP_MD_CTX_new();
    if (evp_md_ctx == NULL) {
        sprintf(err, "could not create evp md ctx: %s",
                ERR_error_string(ERR_get_error(), NULL));
        errors_append(err);
        return false;
    }

    if (EVP_DigestInit_ex2(evp_md_ctx, evp_md_st, NULL) != 1) {
        sprintf(err, "message digest initialization failed: %s\n",
                ERR_error_string(ERR_get_error(), NULL));
        errors_append(err);
        EVP_MD_CTX_free(evp_md_ctx);
        return false;
    };

    if (EVP_DigestUpdate(evp_md_ctx, in, in_length) != 1) {
        sprintf(err, "message digest update failed: %s\n",
                ERR_error_string(ERR_get_error(), NULL));
        errors_append(err);
        EVP_MD_CTX_free(evp_md_ctx);
        return false;
    };

    if (EVP_DigestFinal_ex(evp_md_ctx, out, out_length) != 1) {
        sprintf(err, "message digest finalization failed: %s\n",
                ERR_error_string(ERR_get_error(), NULL));
        errors_append(err);
        EVP_MD_CTX_free(evp_md_ctx);
        return false;
    };

    EVP_MD_CTX_free(evp_md_ctx);
    return true;
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
 *      true if success, false if failure
 * */
bool sri_calculate(char *digest_type, char *in, unsigned int in_length,
                   unsigned char *out, unsigned int *out_length) {
    unsigned int digest_type_length = strlen(digest_type);
    unsigned char digest[BUFSIZ], digest_base64[BUFSIZ];
    unsigned int digest_length, digest_base64_length;
    char err[BUFSIZ];

    if (!sri_calculate_digest(digest_type, in, in_length, digest,
                              &digest_length)) {
        sprintf(err, "could not calculate '%s' sri: could not calculate digest",
                digest_type);
        errors_append(err);
        return false;
    };

    if (!sri_calculate_base64(digest, digest_length, digest_base64,
                              &digest_base64_length)) {
        sprintf(
            err,
            "could not calculate '%s' sri: could not convert digest to base64",
            digest_type);
        errors_append(err);
        return false;
    };
    strncat((char *)out, digest_type, digest_type_length);
    strncat((char *)out, "-", 1);
    strncat((char *)out, (char *)digest_base64, digest_base64_length);
    *out_length = digest_type_length + 1 + digest_base64_length;
    return true;
}

/*
 * Write sri to stream
 *
 * Args:
 *     digest_type: openssl digest type (sha256, for example)
 *     in: input string
 *     in_length: input string length
 *     out: output string (openssl digest)
 *     out_length: length of the output string
 *
 * Returns:
 *      true if success, false if failure
 * */
bool sri_write(char *digest_type, FILE *stream, char *in) {
    unsigned char sri[BUFSIZ];
    unsigned int sri_length;
    char err[BUFSIZ];
    if (!sri_calculate(digest_type, in, strlen(in), sri, &sri_length)) {
        sprintf(err, "could not write '%s' sri", digest_type);
        errors_append(err);
        return false;
    }
    if (fwrite(sri, sizeof sri[0], sri_length, stream) != sri_length) {
        sprintf(err, "could not write '%s' sri: %s", digest_type,
                strerror(errno));
        errors_append(err);
        return false;
    };
    return true;
}
