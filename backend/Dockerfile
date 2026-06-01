FROM golang:1.25-alpine
 
WORKDIR /app
 
COPY go.mod ./
RUN go mod download && go mod tidy
 
COPY . .
RUN go build -o server .
 
EXPOSE 8000
CMD ["./server"]