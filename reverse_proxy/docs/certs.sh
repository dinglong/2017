#!/bin/bash

openssl req -newkey rsa:4096 -nodes -sha256 -keyout root-ca.key -x509 -days 365 -out root-ca.crt

# wrapper
openssl req -newkey rsa:4096 -nodes -sha256 -keyout wrapper.key -out wrapper.csr

echo subjectAltName = IP:192.168.1.50 > wrapper.cnf
openssl x509 -req -days 365 -in wrapper.csr -CA root-ca.crt -CAkey root-ca.key -CAcreateserial -extfile wrapper.cnf -out wrapper.crt

# proxy
openssl req -newkey rsa:4096 -nodes -sha256 -keyout proxy.key -out proxy.csr

echo subjectAltName = IP:192.168.1.50 > proxy.cnf
openssl x509 -req -days 365 -in proxy.csr -CA root-ca.crt -CAkey root-ca.key -CAcreateserial -extfile proxy.cnf -out proxy.crt
