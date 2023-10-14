FROM golang:1.20

WORKDIR /app

COPY . /app

RUN apt -y update
RUN apt -y upgrade
RUN apt install -y python3-pip
RUN apt install -y python3-dev
RUN apt install -y ffmpeg
RUN pip install pillow --break-system-packages


RUN go mod download
RUN go build main.go
ENV GIN_MODE=release

EXPOSE 4232
CMD ["./main"]