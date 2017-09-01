#!/bin/bash

# ***************************************************************************************************************************************************************************************************************
# 1. generates root-ca
openssl genrsa -out "root-ca.key" 4096
openssl req -new -key "root-ca.key" -out "root-ca.csr" -sha256 -subj '/C=US/ST=CA/L=San Francisco/O=Docker/CN=Notary Testing CA'

cat > "root-ca.cnf" <<EOL
[root_ca]
basicConstraints = critical,CA:TRUE,pathlen:1
keyUsage = critical, nonRepudiation, cRLSign, keyCertSign
subjectKeyIdentifier=hash
EOL

openssl x509 -req -days 3650 -in "root-ca.csr" -signkey "root-ca.key" -sha256 -out "root-ca.crt" -extfile "root-ca.cnf" -extensions root_ca
echo "generate root-ca.crt ok"

# ***************************************************************************************************************************************************************************************************************
# 2. generate intermediate-ca
openssl genrsa -out "intermediate-ca.key" 4096
openssl req -new -key "intermediate-ca.key" -out "intermediate-ca.csr" -sha256 -subj '/C=US/ST=CA/L=San Francisco/O=Docker/CN=Notary Intermediate Testing CA'

cat > "intermediate-ca.cnf" <<EOL
[intermediate_ca]
authorityKeyIdentifier=keyid,issuer
basicConstraints = critical,CA:TRUE,pathlen:0
extendedKeyUsage=serverAuth,clientAuth
keyUsage = critical, nonRepudiation, cRLSign, keyCertSign
subjectKeyIdentifier=hash
EOL

openssl x509 -req -days 3650 -in "intermediate-ca.csr" -sha256 -CA "root-ca.crt" -CAkey "root-ca.key"  -CAcreateserial -out "intermediate-ca.crt" -extfile "intermediate-ca.cnf" -extensions intermediate_ca
echo "generate intermediate-ca.crt ok"

# ***************************************************************************************************************************************************************************************************************
# 3. generate notary-server
openssl genrsa -out "notary-server.key" 4096
openssl req -new -key "notary-server.key" -out "notary-server.csr" -sha256 -subj '/C=US/ST=CA/L=San Francisco/O=Docker/CN=notary-server'

cat > "notary-server.cnf" <<EOL
[notary_server]
authorityKeyIdentifier=keyid,issuer
basicConstraints = critical,CA:FALSE
extendedKeyUsage=serverAuth,clientAuth
keyUsage = critical, digitalSignature, keyEncipherment
subjectAltName = DNS:notary-server, DNS:notaryserver, DNS:localhost, IP:192.168.1.50
subjectKeyIdentifier=hash
EOL

openssl x509 -req -days 750 -in "notary-server.csr" -sha256 -CA "intermediate-ca.crt" -CAkey "intermediate-ca.key"  -CAcreateserial -out "notary-server.crt" -extfile "notary-server.cnf" -extensions notary_server

# append the intermediate cert to this one to make it a proper bundle
cat "intermediate-ca.crt" >> "notary-server.crt"
echo "generate notary-server.crt ok"

# ***************************************************************************************************************************************************************************************************************
# 4. generate notary-signer
openssl genrsa -out "notary-signer.key" 4096
openssl req -new -key "notary-signer.key" -out "notary-signer.csr" -sha256 -subj '/C=US/ST=CA/L=San Francisco/O=Docker/CN=notary-signer'

cat > "notary-signer.cnf" <<EOL
[notary_signer]
authorityKeyIdentifier=keyid,issuer
basicConstraints = critical,CA:FALSE
extendedKeyUsage=serverAuth,clientAuth
keyUsage = critical, digitalSignature, keyEncipherment
subjectAltName = DNS:notary-signer, DNS:notarysigner, DNS:localhost, IP:192.168.1.50
subjectKeyIdentifier=hash
EOL

openssl x509 -req -days 750 -in "notary-signer.csr" -sha256 -CA "intermediate-ca.crt" -CAkey "intermediate-ca.key"  -CAcreateserial -out "notary-signer.crt" -extfile "notary-signer.cnf" -extensions notary_signer

# append the intermediate cert to this one to make it a proper bundle
cat "intermediate-ca.crt" >> "notary-signer.crt"
echo "generate notary-signer.crt ok"

# ***************************************************************************************************************************************************************************************************************
# Clean workspace
rm "root-ca.cnf" "root-ca.csr"
rm "root-ca.key" "root-ca.srl"
rm "intermediate-ca.cnf" "intermediate-ca.csr"
rm "intermediate-ca.key" "intermediate-ca.srl"
rm "notary-server.cnf" "notary-server.csr"
rm "notary-signer.cnf" "notary-signer.csr"

