FROM hub.fif.tech/base/golang:1.16-latest

WORKDIR /app

# Pending
COPY server .
COPY vault ./vault

ENTRYPOINT envconsul -config="vault/config.${ENVIRONMENT}.hcl" ./server