# -*- coding: utf-8 -*-
"""
MIT License

Copyright (c) 2020 Tim Wilson

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

https://github.com/twilsonco/PyPhotoCollage/blob/068551c93e8aed5d24f9846d9ea0f01f6a12efcb/LICENSE


Created on May 24, 2020

@author: Tim Wilson
"""
import math
from PIL import ImageOps
from PIL import Image
from PIL.Image import Image as ImageType
from typing import List
from maths import linear_partition, clamp, ensure_even
from config import save_image, width_arg, height_arg
from io import BytesIO
import time

# got idea from https://medium.com/@jtreitz/the-algorithm-for-a-perfectly-balanced-photo-gallery-914c94a5d8af


# takes list of PIL image objects and returns the collage as a PIL image object
def create_collage(img_list: List[ImageType]) -> ImageType:
    max_height = max([img.height for img in img_list])
    img_list = [
        (
            img.resize((int(img.width / img.height * max_height), max_height))
            if img.height < max_height
            else img
        )
        for img in img_list
    ]

    total_width = sum([img.width for img in img_list])
    avg_width = total_width / len(img_list)
    target_width = avg_width * math.sqrt(len(img_list))

    num_rows = clamp(int(round(total_width / target_width)), len(img_list))
    if num_rows == 1:
        img_rows = [img_list]
    elif num_rows == len(img_list):
        img_rows = [[img] for img in img_list]
    else:

        aspect_ratios = [int(img.width / img.height * 100) for img in img_list]

        img_rows = linear_partition(aspect_ratios, num_rows, img_list)

        row_widths = [sum([img.width for img in row]) for row in img_rows]
        min_row_width = min(row_widths)
        row_width_ratios = [min_row_width / w for w in row_widths]
        img_rows = [
            [
                img.resize(
                    (int(img.width * width_ratio), int(img.height * width_ratio))
                )
                for img in row
            ]
            for row, width_ratio in zip(img_rows, row_width_ratios)
        ]

    row_widths = [sum([img.width for img in row]) for row in img_rows]
    row_heights = [max([img.height for img in row]) for row in img_rows]

    w, h = min(row_widths), sum(row_heights)

    w = ensure_even(w)
    h = ensure_even(h)

    result_image = Image.new("RGBA", (w, h))
    xPos, yPos = (0, 0)

    for row in img_rows:
        for img in row:
            result_image.paste(img, (xPos, yPos))
            xPos += img.width
            continue
        yPos += max([img.height for img in row])
        xPos = 0
        continue

    return result_image


def make_collage(images: List[bytes], output: str) -> float:
    print(output)
    start = time.time()
    if len(images) == 1:
        save_image(output, Image.open(images[0]), 0, 0)
        return time.time() - start

    pilImages = []

    for image in images:
        img = Image.open(BytesIO(image))
        exif = img.getexif()
        for k in exif.keys():
            if k != 0x0112:
                exif[k] = (
                    None  # If I don't set it to None first (or print it) the del fails for some reason.
                )
            del exif[k]

        img.info["exif"] = exif.tobytes()
        # Rotate the image based on EXIF orientation data
        ImageOps.exif_transpose(img)
        # uncomment to save space and little time
        # commented makes bigger images
        # if img.height > init_height:
        #    new_width = int(img.width / img.height * init_height)
        #    pilImages.append(img.resize((new_width, init_height), Image.LANCZOS))
        # else:
        pilImages.append(img)

    collage = create_collage(pilImages)

    if collage.width > width_arg:
        collage = collage.resize(
            (width_arg, int(collage.height / collage.width * width_arg)),
            Image.LANCZOS,
        )
    elif collage.height > height_arg:
        collage = collage.resize(
            (int(collage.width / collage.height * height_arg), height_arg),
            Image.LANCZOS,
        )
    save_image("collages/" + output, collage, collage.width, collage.height)

    return time.time() - start
