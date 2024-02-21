FROM golang:1.22-rc-alpine
# alpine is not necessary. it can be 1.20 or other tags

WORKDIR /app

# technically you dont need to copy all the files. ONLY the go stuff
COPY . .

RUN go get
# or
# RUN go mod download

RUN go build -o bin .

EXPOSE 8080

ENTRYPOINT [ "/app/bin" ]