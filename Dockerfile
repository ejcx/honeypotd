FROM ubuntu

COPY honeypotd honeypotd
COPY authorized_keys authorized_keys
COPY id_rsa id_rsa
COPY id_rsa.pub id_rsa.pub

ENTRYPOINT ["./honeypotd", "env"]
