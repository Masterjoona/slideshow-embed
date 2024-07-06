"""
The MIT License (MIT)

Copyright (c) 2014-2015 vingtcinq.io

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
https://github.com/VingtCinq/python-resize-image/blob/9c9a1f6d61abf3f5072ca0934963fcd75ed24c08/LICENSE

Kanged from https://github.com/VingtCinq/python-resize-image/blob/9c9a1f6d61abf3f5072ca0934963fcd75ed24c08/resizeimage/resizeimage.py#L98-L114
"""

from PIL.Image import Image as ImageType
from PIL.Image import Resampling
from PIL import Image

import math
from typing import List
import time
import io
import os
from maths import ensure_even
from config import save_image, find_largest_image_dimensions, temporary_dir


def resize_contain(
    image_buffer: bytes,
    size: tuple[int, int],
) -> ImageType:
    """
    Resize image according to size.
    image:      a Pillow image instance
    size:       a list of two integers [width, height]
    """
    image = Image.open(io.BytesIO(image_buffer))
    img = image.copy()
    img.thumbnail((size[0], size[1]), Resampling.LANCZOS)
    background = Image.new("RGBA", (size[0], size[1]), (0, 0, 0))
    img_position = (
        int(math.ceil((size[0] - img.size[0]) / 2)),
        int(math.ceil((size[1] - img.size[1]) / 2)),
    )
    background.paste(img, img_position)
    return background.convert("RGBA")


def resize_images(images: List[bytes], output: str) -> float:
    start = time.time()
    os.makedirs(f"{temporary_dir}/{output}", exist_ok=True)
    to_width, to_heigth = find_largest_image_dimensions(images)

    to_width = ensure_even(to_width)
    to_heigth = ensure_even(to_heigth)

    for index, image in enumerate(images):
        cover = resize_contain(image, (to_width, to_heigth))
        save_image(
            f"{temporary_dir}/{output}/{index}.png", cover, to_width, to_heigth, False
        )

    return time.time() - start
