#include <stdbool.h>
#include <stdio.h>

bool sri_calculate(char *digest_type, char *in, unsigned int in_length,
                   unsigned char *out, unsigned int *out_length);

bool sri_write(char *digest_type, FILE *stream, char *in);
