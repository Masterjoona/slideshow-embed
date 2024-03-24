FROM python:3.12-alpine AS builder

WORKDIR /app

RUN pip install --no-cache-dir pillow
COPY ../../fpng_py/build/lib.linux-aarch64-3.10/fpng_py/ /app/fpng_py

FROM python:3.12-alpine

WORKDIR /app

COPY --from=builder /usr/local/lib/python3.12/site-packages /usr/local/lib/python3.12/site-packages

COPY ../../python_side/ /app/

CMD ["python", "server_collage_maker.py"]