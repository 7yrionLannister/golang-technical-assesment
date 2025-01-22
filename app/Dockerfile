# Use the official Golang image
FROM golang:1.23.2

# Uncomment the following line if you need access to the timezone data
# RUN apt-get update && apt-get install -y tzdata

WORKDIR /app

COPY ./app/src .

RUN go mod download

# Use ```RUN go build -o main -tags timetzdata .``` to include the timezone data. This will increase the size of the image.
RUN go build -o main .

EXPOSE 8181

VOLUME ["/app/data"]

CMD ["./main"]