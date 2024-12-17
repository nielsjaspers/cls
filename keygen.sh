#!/bin/bash

# Check if the ./server-crt directory exists, and create it if it doesn't
if [ ! -d "./server-crt" ]; then
    mkdir ./server-crt
    if [ $? -ne 0 ]; then
        echo "Error creating directory ./server-crt"
        exit 1
    fi
fi

# Remove any existing files from the server-crt directory
rm -f ./server-crt/ca.crt ./server-crt/ca.key ./server-crt/ca.srl ./server-crt/server.crt ./server-crt/server.csr ./server-crt/server.key

# Generate the root CA private key
openssl genrsa -out ./server-crt/ca.key 2048
if [ $? -ne 0 ]; then
    echo "Error generating CA key"
    exit 1
fi

# Generate the root CA certificate
openssl req -x509 -new -nodes -key ./server-crt/ca.key -sha256 -days 365 -out ./server-crt/ca.crt -subj "/C=US/ST=State/L=City/O=MyOrg/OU=MyUnit/CN=MyCA"
if [ $? -ne 0 ]; then
    echo "Error generating CA certificate"
    exit 1
fi

# Generate the server private key
openssl genrsa -out ./server-crt/server.key 2048
if [ $? -ne 0 ]; then
    echo "Error generating server key"
    exit 1
fi

# Generate the server certificate signing request (CSR)
openssl req -new -key ./server-crt/server.key -out ./server-crt/server.csr -config ./server-crt/server.cnf
if [ $? -ne 0 ]; then
    echo "Error generating server CSR"
    exit 1
fi

# Sign the server certificate with the CA key
openssl x509 -req -in ./server-crt/server.csr -CA ./server-crt/ca.crt -CAkey ./server-crt/ca.key -CAcreateserial -out ./server-crt/server.crt -days 365 -extfile ./server-crt/server.cnf -extensions v3_ext
if [ $? -ne 0 ]; then
    echo "Error signing server certificate"
    exit 1
fi

# Verify the server certificate
openssl verify -CAfile ./server-crt/ca.crt ./server-crt/server.crt
if [ $? -ne 0 ]; then
    echo "Error verifying server certificate"
    exit 1
fi

echo "CA and server certificates generated and verified successfully!"

