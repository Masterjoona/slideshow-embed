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

# Modify below line to copy the compiled .so file or comment it out if dont want to use fpng
COPY ./fpng_py/build/lib.linux-aarch64-3.10/fpng_py /app/fpng_py

WORKDIR /app
RUN rm -rf /tmp/*
ENV GIN_MODE=release
EXPOSE 4232
CMD ["./meow"]
