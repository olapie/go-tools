openssl pkcs12 -in cert.p12 -out cert.pem -clcerts -nokeys
openssl pkcs12 -in cert.p12 -out key.pem -nocerts -nodes