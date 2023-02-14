#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <getopt.h>
#include <string.h>

#include "dig_license.h"
#include "gen_clientId.h"

#define MAX_FILE_PATH 1024
#define MAX_BUF_LEN   65535
#define MAX_QPS_LEN   10
#define QPS_SHIFT     33
#define MAX_DATETIME_LEN  20
#define MAX_CLIENT_ID_LEN 33
#define MAX_SIGN_LEN 512
#define MAX_EMPTY_SIZE 512

static void
print_help()
{
    fprintf(stderr, "Usage: keygen [OPTION]...\n");
    fprintf(stderr, "Supported options:\n"
            "  -p public key file path.\n"
            "  -q private key file path.\n"
            "  -i machine information file path.\n"
            "  -l license file path.\n"
            "  -r restart parameters.\n"
            "  -e other parameters.\n"
            "  -h print help info\n");
}

int write_data(char *data_buf, int data_len, FILE* fp) {

    //write data length
    static char data_length[MAX_BUF_LEN] = {0};
    memset(data_length, 0, sizeof data_length);

    sprintf(data_length, "%d", data_len);
    int i = 0;
    for (; i < strlen(data_length); ++i) {
        data_length[i] += (i + QPS_SHIFT);
    }
    size_t write_data_len = fwrite(data_length, sizeof(char), MAX_QPS_LEN, fp);
    if (write_data_len < 0) {
        fprintf(stderr, "write data_length error\n");
        return -1;
    }

    // write data
    i = 0;
    for (; i < data_len; ++i) {
        data_buf[i] += (i + QPS_SHIFT);
    }
    uint8_t *write_pos = (uint8_t *)data_buf;
    size_t write_len = fwrite(write_pos, sizeof(char), data_len, fp);
    if (write_len < 0) {
        fprintf(stderr, "write data failed\n");
        return -1;
    }
    return 0;
}

int
main(int argc, char *argv[])
{
    if (2 != argc && argc !=13)
    {
        print_help();
        return 1;
    }

    char pub_file[MAX_FILE_PATH] = {0};
    char pri_file[MAX_FILE_PATH] = {0};
    char info_file[MAX_FILE_PATH] = {0};
    char license_file[MAX_FILE_PATH] = {0};

    char extend_license_config[MAX_BUF_LEN] = {0};
    char extend_license_length[MAX_FILE_PATH] = {0};
    char restart_param_config[MAX_BUF_LEN] = {0};

    int opt = -1;
    while((opt = getopt(argc, argv, "hp:q:i:l:e:r:")) != -1)
    {
        switch (opt)
        {
            case 'p' :
                strncpy(pub_file, optarg, MAX_FILE_PATH);
                break;
            case 'q' :
                strncpy(pri_file, optarg, MAX_FILE_PATH);
                break;
            case 'i' :
                strncpy(info_file, optarg, MAX_FILE_PATH);
                break;
            case 'l' :
                strncpy(license_file, optarg, MAX_FILE_PATH);
                break;
            case 'r' :
                strncpy(restart_param_config, optarg, MAX_BUF_LEN);
                break;
            case 'e' :
                strncpy(extend_license_config, optarg, MAX_BUF_LEN);
                break;
            default :
                print_help();
                exit(1);
        }
    }

    // generate public and private key file
    if (gen_key_into_file(pub_file, pri_file, info_file)) {
        fprintf(stderr, "generate key failed\n");
        return 1;
    }

    EVP_PKEY *pri_key = read_private_key_from_file(pri_file);
    if (!pri_key) {
        fprintf(stderr, "read private key file failed\n");
        return -1;
    }

    // read machine info
    char info_buf[MAX_INFO_LEN] = {0};
    FILE *fp = fopen(info_file, "r");
    if (!fp) {
        fprintf(stderr, "open infomation file failed\n");
        return -1;
    }
    size_t info_len = fread(info_buf, sizeof(char), MAX_INFO_LEN, fp);
    fclose(fp);
    if (info_len <= 0) {
        fprintf(stderr, "read infomation file failed\n");
        return -1;
    }

    // sign machine info with private key
    uint8_t sign_buf[MAX_SIGN_LEN] = {0};
    int sign_len = -1;
    int ret = sign_buf_with_private_key((const uint8_t *)info_buf, info_len, 
            sign_buf, &sign_len, pri_key);
    if (ret) {
        fprintf(stderr, "sign infomation failed\n");
        return -1;
    }

    // write license file
    fp = fopen(license_file, "w");
    if (!fp) {
        fprintf(stderr, "open digest file failed\n");
        return -1;
    }

    //write restart param
    write_data(restart_param_config, strlen(restart_param_config), fp);

    //write sign length
    write_data(sign_buf, sign_len, fp);

    //add client id
    char client_id[MAX_CLIENT_ID_LEN] = {0};
    if (gen_client_id(info_file, client_id)) {
        fprintf(stderr, "gen client id error\n");
        return -1;
    }
    char* client_name = (char*)"client_id:#:";
    size_t extend_license_config_length = strlen(extend_license_config);
    size_t new_extend_license_config_length = strlen(client_id) + extend_license_config_length + strlen(client_name) + 10;

    char *new_extend_license_config = (char*)malloc(new_extend_license_config_length * sizeof(char));
    if (new_extend_license_config == NULL) {
        fprintf(stderr, "malloc error\n");
        return -1;
    }
    memset(new_extend_license_config, 0, new_extend_license_config_length);
    sprintf(new_extend_license_config, "%s%s;%s", client_name, client_id, extend_license_config);

    //write extend config
    write_data(new_extend_license_config, new_extend_license_config_length, fp);
    fclose(fp);
    return 0;
}
