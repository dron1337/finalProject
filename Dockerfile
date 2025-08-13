FROM golang:1.24.2
WORKDIR /cmd
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /scheduler
CMD ["/scheduler"]