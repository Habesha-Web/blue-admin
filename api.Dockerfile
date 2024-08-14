
FROM golang:latest

RUN apt install -y libc6 libc-bin

RUN apt -y update && apt -y upgrade

RUN apt -y install build-essential pkg-config g++ git cmake yasm

RUN apt install build-essential pkg-config git

WORKDIR /playground/

COPY docs /playground/

COPY main /playground/

COPY server.pem  /playground/

COPY server-key.pem  /playground/

COPY configs /playground/

RUN chmod +x main

CMD ["./main","prod"]
 