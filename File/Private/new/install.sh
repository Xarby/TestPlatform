gcc -g -o license -lcrypto license.c dig_license.c gen_clientId.c
mv license /usr/local/bin/license_new
cp publisher_new /usr/bin/publisher_new
chmod 755 /usr/bin/publisher_new