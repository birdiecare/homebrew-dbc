
FROM ubuntu:20.04 AS sessionmanagerplugin
# session manager is not available in alpine, so we grab it from Ubuntu
ADD https://s3.amazonaws.com/session-manager-downloads/plugin/latest/ubuntu_arm64/session-manager-plugin.deb .
RUN dpkg -i "session-manager-plugin.deb"

FROM golang:alpine AS homebrew-dbc


COPY --from=sessionmanagerplugin /usr/local/sessionmanagerplugin/bin/session-manager-plugin /usr/local/bin/
RUN apk update
RUN apk add aws-cli
RUN apk add socat

WORKDIR /homebrew-dbc
COPY . .

# This 👇🏼 should be done in the repo probably.
RUN go get -u github.com/aws/aws-sdk-go-v2/...
RUN go build

EXPOSE 5432
# socat routes requests to port 5432 to port 5555 where the ssm tunnel is running.
# this is needed because the ssm tunnel only accepts connections from inside the same container.
CMD [ "sh", "-c", "socat tcp-listen:5432,reuseaddr,fork tcp:localhost:5555 & ./dbc connect --host=${HOST} --localport=5555" ]

