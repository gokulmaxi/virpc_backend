FROM golang:1.18.2-bullseye
# add curl for healthcheck
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
    curl 

WORKDIR /app
ENV GOPATH=/
COPY ./app/go.mod .
RUN go install github.com/cosmtrek/air@latest
# RUN go get

CMD ["air"]
