#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <string.h>
#include <unistd.h>
#include <time.h> 
#include "openssl/md5.h" 

bool file_exist_or_not(const char *filename)
{
    if(!access(filename, F_OK))
        return true;
    return false;
}

int gen_client_id(const char* machine_file, char *client_id)
{
    int i;
    time_t t;
    struct tm* ti;
    const size_t date_len = 9;
    char today[date_len];
    const size_t machine_data_len = 65;
    char machine_data[machine_data_len];
    const size_t origin_data_len = date_len + machine_data_len -1;
    char origin_data[origin_data_len];
    const size_t md5_len = 16;
    unsigned char md5[md5_len];                                                                    
       
    int pt_index = 0; 
    int password_table[100] = {0};
    MD5_CTX ctx;
    const int root_password_len = 32;
    char new_root_password[root_password_len + 1];
    const int origin_password_len = 32;
    char origin_password[origin_password_len + 1];

    if (file_exist_or_not(machine_file) == false) {
        printf("Error: not exist machine file:%s\n", machine_file);
        return -1;
    }

    FILE *fp = fopen(machine_file, "r");
    if (!fp) {
        printf("Error: open machine file error\n");
        return -1;
    }
        
    memset(machine_data, 0, machine_data_len);
    size_t read_size = fread(machine_data, sizeof(char), 64, fp);
    if (read_size < machine_data_len - 1) {
        printf("Error: read machine data error\n");
        fclose(fp);
        return -1;
    }
    fclose(fp);
    time(&t);
    ti = localtime(&t);
    memset(today, 0, date_len);
    sprintf(today, "%02d%02d%02d", ti->tm_year + 1900, ti->tm_mon + 1, ti->tm_mday);

    memset(origin_data, 0, origin_data_len);
    if (ti->tm_mday % 2 == 0) {
        memcpy(origin_data, machine_data, machine_data_len);
        memcpy(origin_data + machine_data_len - 1, today, date_len);
    } else {
        memcpy(origin_data, today, date_len);
        memcpy(origin_data + date_len - 1, machine_data, machine_data_len);
    }

    {
       //add hour to origin_data
       origin_data[0] += ti->tm_hour;
       //add minute to origin_data
       origin_data[1] += ti->tm_min;
       //add second to origin_data
       origin_data[2] += ti->tm_sec;
    }

    { 
        char command[100] = {0};
        memset(md5, 0, md5_len);

        MD5_CTX tctx;
        MD5_Init(&tctx);
        MD5_Update(&tctx, origin_data, origin_data_len);
        MD5_Final(md5, &tctx);
        memset(origin_password, 0, origin_password_len + 1);
        for (i = 0; i < md5_len; ++i)
            sprintf(origin_password + i * 2, "%02x", md5[i]);
    }

    {
        for (i = 0; i < origin_password_len - 1 ; ++i) {
           origin_password[i] = origin_password[i] + (int)origin_password[i] / 2; 
        }
    }
    memset(md5, 0, md5_len);
    MD5_Init(&ctx);
    MD5_Update(&ctx, origin_password, origin_password_len);
    MD5_Final(md5, &ctx);
    memset(origin_password, 0, origin_password_len + 1);
    for (i = 0; i < md5_len; ++i)
        sprintf(origin_password + i * 2, "%02x", md5[i]);

    //init password table
    pt_index = 0;
    for(i = (int)'a'; i <= (int)'z'; ++i) {
        password_table[pt_index++] = i;
    }
    for (i = (int)'A'; i <= (int)'Z'; ++i) {
        password_table[pt_index++] = i;
    }
    for (i = 0; i <= 9; ++i) {
        password_table[pt_index++] = i + 48;
    }

    memset(client_id, 0, root_password_len + 1);
    for (i = 0; i < origin_password_len; ++i) {
        int t = (int)origin_password[i];
        t = t > 0 ? t : -t;
        client_id[i] = password_table[t % pt_index];
    }
    return 0;
}
