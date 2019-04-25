FROM golang:1.11.2
WORKDIR /app
COPY . /app
RUN go build -o plague_doctor
ENV GIN_MODE release
CMD ["./plague_doctor"]