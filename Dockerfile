FROM golang:1.22

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY data/ ./data/
COPY exceptions/ ./exceptions/
COPY routes/ ./routes/
COPY security/ ./security/
COPY *.go ./
RUN go build -o ./export.out

EXPOSE 8080/tcp
CMD ["/app/export.out"]
