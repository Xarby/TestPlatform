gcc -g -o license -lcrypto license.c dig_license.c gen_clientId.c
mv license /usr/local/bin/license_old
cp publisher_old /usr/bin/publisher_old
chmod 755 /usr/bin/publisher_old
