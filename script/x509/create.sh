#!/bin/bash

# compact win MSYS
export MSYS_NO_PATHCONV=1

# script dir
SCRIPT_DIR=$(cd $(dirname ${BASH_SOURCE[0]}); pwd)

pushd ${SCRIPT_DIR}/../../
rm -rf certs
mkdir certs

cd certs

CONFIG_FILE=../script/x509/openssl.cnf

# Create the server CA certs.
openssl req -x509                                     \
  -newkey rsa:4096                                    \
  -nodes                                              \
  -days 3650                                          \
  -keyout server_ca_key.pem                           \
  -out server_ca_cert.pem                             \
  -subj /C=US/ST=CA/L=SVL/O=gRPC/CN=ut-server_ca/     \
  -config ${CONFIG_FILE}                              \
  -extensions test_ca                                 \
  -sha256

# Create the client CA certs.
openssl req -x509                                     \
  -newkey rsa:4096                                    \
  -nodes                                              \
  -days 3650                                          \
  -keyout client_ca_key.pem                           \
  -out client_ca_cert.pem                             \
  -subj /C=US/ST=CA/L=SVL/O=gRPC/CN=ut-client_ca/     \
  -config ${CONFIG_FILE}                              \
  -extensions test_ca                                 \
  -sha256

# Generate a server cert.
openssl genrsa -out server_key.pem 4096
openssl req -new                                    \
  -key server_key.pem                               \
  -days 3650                                        \
  -out server_csr.pem                               \
  -subj /C=US/ST=CA/L=SVL/O=gRPC/CN=ut-server1/     \
  -config ${CONFIG_FILE}                            \
  -reqexts test_server
openssl x509 -req           \
  -in server_csr.pem        \
  -CAkey server_ca_key.pem  \
  -CA server_ca_cert.pem    \
  -days 3650                \
  -set_serial 1000          \
  -out server_cert.pem      \
  -extfile ${CONFIG_FILE}   \
  -extensions test_server   \
  -sha256
openssl verify -verbose -CAfile server_ca_cert.pem  server_cert.pem

# Generate a client cert.
openssl genrsa -out client_key.pem 4096
openssl req -new                                    \
  -key client_key.pem                               \
  -days 3650                                        \
  -out client_csr.pem                               \
  -subj /C=US/ST=CA/L=SVL/O=gRPC/CN=ut-client1/    \
  -config ${CONFIG_FILE}                            \
  -reqexts test_client
openssl x509 -req           \
  -in client_csr.pem        \
  -CAkey client_ca_key.pem  \
  -CA client_ca_cert.pem    \
  -days 3650                \
  -set_serial 1000          \
  -out client_cert.pem      \
  -extfile ${CONFIG_FILE}   \
  -extensions test_client   \
  -sha256
openssl verify -verbose -CAfile client_ca_cert.pem  client_cert.pem

rm *_csr.pem

mkdir certs_cli
cp -rf server_ca_cert.pem ./certs_cli/
cp -rf client_cert.pem ./certs_cli/
cp -rf client_key.pem ./certs_cli/

mkdir certs_svc
cp -rf client_ca_cert.pem ./certs_svc/
cp -rf server_cert.pem ./certs_svc/
cp -rf server_key.pem ./certs_svc/

popd