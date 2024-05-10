import http.server
import cgi  # deprecated do something in the future
from typing import List
from collage_maker import make_collage
from resizing import resize_images
from config import collages_path


def parse_post(
    handler: http.server.BaseHTTPRequestHandler,
) -> tuple[List[bytes], str]:
    content_length = int(handler.headers.get("Content-Length", 0))
    if content_length == 0:
        return [], ""

    environ = {
        "REQUEST_METHOD": "POST",
        "CONTENT_TYPE": handler.headers.get("Content-Type"),
        "CONTENT_LENGTH": content_length,
    }

    form = cgi.FieldStorage(fp=handler.rfile, headers=handler.headers, environ=environ)
    video_id = form.getvalue("video_id")

    return form.getlist("images"), str(video_id).replace("b", "").replace("'", "")


class ImageHTTPRequestHandler(http.server.BaseHTTPRequestHandler):
    def do_POST(self):
        if self.path == "/collage":
            self.make_collage_server()
        elif self.path == "/resize":
            self.make_resize_server()
        else:
            self.send_error(404, "Not Found")

    def make_collage_server(self):
        images, video_id = parse_post(self)
        if not images:
            self.send_error(400, "No images provided")
            return
        time_taken = make_collage(images, f"{collages_path}/collage-{video_id}.png")
        self.send_response(200)
        self.end_headers()
        response_message = f"Collage created in {time_taken:.4f} seconds".encode()
        print(response_message)
        self.wfile.write(bytes(response_message))

    def make_resize_server(self):
        images, video_id = parse_post(self)
        if not images:
            self.send_error(400, "No images provided")
            return
        time_taken = resize_images(images, video_id)
        self.send_response(200)
        self.end_headers()
        response_message = f"Images resized in {time_taken:.4f} seconds".encode()
        print(response_message)
        self.wfile.write(bytes(response_message))


if __name__ == "__main__":
    server_address = ("", 9700)
    httpd = http.server.HTTPServer(server_address, ImageHTTPRequestHandler)
    print("Starting server at http://localhost:9700")
    httpd.serve_forever()
