# Official docker container image
FROM --platform=linux/x86_64 golang:1.15.10
#FROM golang:1.15.10

# development channel. Either dev, preprod or prod.
# docker-compose will pass in dev channel.
# wil default to prod when deplying to aws default to prod.
# docker-compose won't be on aws
# preprod is for using aws databases and S3 but local docker
# netwroking for inter container communicatation
# prod will have AWS databse and S3 and AWS vpc networking.
ARG channel=preprod
# prod

# the arg will be the value of the channel_var environment variable.
ENV channel_env_var=$channel

# turn on modules with the env variable
ENV GO111MODULE=on

WORKDIR /app

# copy go mod files first
COPY go/go.mod go/go.sum ./

# download go mod dependencies
RUN go mod download

# This was needed to fix a bug for some reason
COPY .config .

# copy the rest of the code
COPY go .

# go build and name output file
RUN go build -o go-createmusic .

#expore port to outside the container
EXPOSE 8000

# run go output file
CMD ["./go-createmusic"]
