FROM arm64v8/alpine

RUN apk --no-cache add ca-certificates curl bash xz-libs git gcompat python3 py3-pip

WORKDIR /tmp
RUN curl -L -O https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-arm64-static.tar.xz
RUN tar -xf ffmpeg-release-arm64-static.tar.xz && cd ff* && mv ff* /usr/local/bin

WORKDIR /app

COPY ./meow /app
COPY templates /app/templates
COPY collage_maker.py /app/collage_maker.py
RUN pip install pillow --break-system-packages

ENV GIN_MODE=release
EXPOSE 4232
CMD ["./meow"]
