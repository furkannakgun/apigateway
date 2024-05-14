FROM harbor.taas.vodafone.local/baseimage/golang:1.20.2

# RUN groupadd -g 1001 -r naapgroup && \
#    useradd -u 1001 -r naapuser -G naapgroup && \
#    chown -R naapuser:naapgroup /var/log

WORKDIR /app/src

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
ARG username
ARG password 
RUN echo machine github.vodafone.com login $username password $password > ~/.netrc
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /app/nebula

WORKDIR /app

RUN rm -rf src

EXPOSE 8080

# USER naapuser

ENTRYPOINT ["/app/nebula"]