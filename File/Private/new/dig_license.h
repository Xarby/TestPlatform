#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>

#include <openssl/evp.h>

#define MAX_KEY_LEN 10240
#define MAX_INFO_LEN 2048

int 
gen_key_into_file(const char *public_key_file, const char *private_key_file, const char *info_file);

EVP_PKEY *
read_private_key_from_file(const char *private_key_file);

int
sign_buf_with_private_key(const uint8_t *input_buf, int input_len, 
        uint8_t *output_buf, int *output_len, EVP_PKEY *pkey);
