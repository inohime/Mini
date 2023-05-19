FROM golang:1.19

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN if [ ! -d "./bin/" ]; then mkdir bin/ && cd bin; else cd bin; fi; go build -o /bin/mini -v ../src/main.go

CMD ["/bin/mini"]