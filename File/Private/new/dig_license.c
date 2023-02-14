#include <string.h>
#include <unistd.h>

#include <openssl/rsa.h>
#include <openssl/pem.h>
#include <openssl/err.h>

#include "dig_license.h"

#define MAX_BUF_LEN 65535
#define SHA512_LEN  64
#define SHA256_LEN  32

#define ENCRYPTO_INC    16
#define MAX_ENCRYPTO_LEN    65535
 
static int
encrypt_buf_with_public_key(const uint8_t *buf, const uint8_t *public_key_buf, int public_key_size, uint8_t *output_buf, int *output_len);

// generate public key & private key file
// public key: 1、sha256 machine info 2、encrypt digest info with public key
int
gen_key_into_file(const char *public_key_file, const char *private_key_file, const char *machine_info_file)
{
    // 1. generate public && private key file
    BIO *public_out = BIO_new_file(public_key_file, "w");
    if (!public_out) {
        fprintf(stderr, "fail bio new public_key_file  %s\n", public_key_file);
        return -1;
    }

    BIO *private_out = BIO_new_file(private_key_file, "w");
    if (!private_out) {
        fprintf(stderr, "fail bio new private_key_file  %s\n", private_key_file);
        return -1;
    }

    int bits = 1024;
    unsigned long e = RSA_3;
    RSA *rsa_key = RSA_generate_key(bits, e, NULL, NULL);

    int ret = PEM_write_bio_RSAPrivateKey(private_out, rsa_key, NULL, NULL, 0, NULL, NULL);
    if (1 != ret) {
        fprintf(stderr, "fail write bio private key, ret %d\n", ret);
        return -2;
    }

    ret = PEM_write_bio_RSAPublicKey(public_out, rsa_key);
    if (1 != ret) {
        fprintf(stderr, "fail write bio public key, ret %d\n", ret);
        return -2;
    }

    RSA_free(rsa_key);
    BIO_flush(private_out);
    BIO_free(private_out);
    BIO_flush(public_out);
    BIO_free(public_out);


    // 2. update public key
    // a) read and sha256 machine info 
    char machine_info_buf[MAX_INFO_LEN] = {0};
    FILE *fp = fopen(machine_info_file, "r");
    if (!fp) {
        fprintf(stderr, "fail open %s\n", machine_info_file);
        return -1;
    }

    int info_len = fread(machine_info_buf, sizeof(char), MAX_INFO_LEN, fp);
    fclose(fp);
    if (info_len <= 0) {
        fprintf(stderr, "fail read %s, len %d\n", machine_info_file, info_len);
        return -1;
    }

    uint8_t digest_buf[SHA256_LEN] = {0};
    SHA256((const unsigned char *)machine_info_buf, info_len, digest_buf);

    // b) read public key
    static uint8_t public_key_buf[MAX_BUF_LEN] = {0};
    fp = fopen(public_key_file, "r");
    if (!fp) {
        fprintf(stderr, "open orignal pub ke failed\n");
        return -1;
    }
    int public_key_size = fread(public_key_buf, sizeof(uint8_t), MAX_BUF_LEN, fp);
    fclose(fp);
    if (public_key_size <= 0) {
        fprintf(stderr, "read orig pub file failed\n");
        return -1;
    }
    
    // c) encrypt info with public key
    uint8_t encryp_buf[MAX_BUF_LEN] = {0};
    int encryp_len = -1;
    ret = encrypt_buf_with_public_key(digest_buf, public_key_buf, public_key_size, encryp_buf, &encryp_len);
    if (ret) {
        fprintf(stderr, "encrypt pub key failed\n");
        return -1;
    }

    // d) write encrypt info into public key file
    fp = fopen(public_key_file, "w");
    if (!fp) {
        fprintf(stderr, "open2 pubkey failed\n");
        return -1;
    }
    size_t write_len = fwrite(encryp_buf, sizeof(uint8_t), encryp_len, fp);
    if (write_len != encryp_len) {
        fprintf(stderr, "encrypt info write failed\n");
        fclose(fp);
        return -1;
    }

    fclose(fp);
    return 0;
}

// read private key
EVP_PKEY *
read_private_key_from_file(const char *private_key_file)
{
    if (!private_key_file || access(private_key_file, F_OK))
        return NULL;

    BIO *private_in = BIO_new_file(private_key_file, "rb");
    if (!private_in)
        return NULL;

    RSA *private_key = NULL;
    private_key = PEM_read_bio_RSAPrivateKey(private_in, &private_key, NULL, NULL);
    if (!private_key)
        return NULL;

    EVP_PKEY *pri_key = EVP_PKEY_new();
    if (!pri_key)
    {
        RSA_free(private_key);
        return NULL;
    }

    EVP_PKEY_assign_RSA(pri_key, private_key);
    return pri_key;
}

// encrypt buf with public key
static int
encrypt_buf_with_public_key(const uint8_t *buf, const uint8_t *public_key_buf, int public_key_size, uint8_t *output_buf, int *output_len)
{
    EVP_CIPHER_CTX ctx;
    EVP_CIPHER_CTX_init(&ctx);
    const EVP_CIPHER *cipher = EVP_des_ede3_cbc();
    int ret = EVP_EncryptInit_ex(&ctx, cipher, NULL, buf, NULL);
    if (1 != ret)
        return -1;

    int once_out = -1;
    EVP_EncryptUpdate(&ctx, output_buf, &once_out, public_key_buf, public_key_size);
    *output_len = once_out;
    EVP_EncryptFinal_ex(&ctx, output_buf + once_out, &once_out);
    *output_len = *output_len + once_out;

    EVP_CIPHER_CTX_cleanup(&ctx);

    return 0;
}

int
sign_buf_with_private_key(const uint8_t *input_buf, int input_len, uint8_t *output_buf, int *output_len, EVP_PKEY *pri_key)
{
    if (!input_buf || (input_len <= 0) || !output_buf || !output_len || !pri_key)
        return -1;

    EVP_MD_CTX md_ctx_sign;
    EVP_SignInit(&md_ctx_sign, EVP_sha1());

    EVP_SignUpdate(&md_ctx_sign, input_buf, input_len);
    int err = EVP_SignFinal(&md_ctx_sign, output_buf, (unsigned int *)output_len, pri_key);
    if (1 != err)
    {
        ERR_print_errors_fp(stderr);
        return -1;
    }

    return 0;
}
