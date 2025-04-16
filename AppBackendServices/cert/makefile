CA_KEY=ca.key
CA_CERT=ca.crt
CA_SUBJECT="/C=US/ST=State/L=City/O=MyOrg/CN=MyRootCA"

SERVICE_NAME=auth
SERVICE_KEY=$(SERVICE_NAME).key
SERVICE_CSR=$(SERVICE_NAME).csr
SERVICE_CERT=$(SERVICE_NAME).crt
SERVICE_SUBJECT="/C=US/ST=State/L=City/O=MyOrg/CN=$(SERVICE_NAME)"

all: clean ca service bundle

ca:
	@echo "Generating Root CA..."
	openssl genrsa -out $(CA_KEY) 4096
	openssl req -x509 -new -nodes -key $(CA_KEY) -sha256 -days 3650 -out $(CA_CERT) -subj "$(CA_SUBJECT)"

service:
	@echo "Generating certificate for $(SERVICE_NAME)..."
	openssl genrsa -out $(SERVICE_KEY) 2048
	openssl req -new -key $(SERVICE_KEY) -out $(SERVICE_CSR) -subj "$(SERVICE_SUBJECT)" -config ssl.cnf -extensions v3_req
	openssl x509 -req -in $(SERVICE_CSR) -CA $(CA_CERT) -CAkey $(CA_KEY) -CAcreateserial -out $(SERVICE_CERT) -days 365 -sha256 -extfile ssl.cnf -extensions v3_req

bundle:
	@echo "Bundling certificates..."
	cat $(SERVICE_CERT) $(CA_CERT) > $(SERVICE_NAME)-bundle.crt

clean:
	@echo "Cleaning up old certificates..."
	rm -f *.key *.crt *.csr *.srl *.pem