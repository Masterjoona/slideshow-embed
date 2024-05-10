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

from PIL import Image
from PIL.Image import Image as ImageType
from typing import List
from maths import linear_partition, clamp, ensure_even
from config import save_image, width_arg, height_arg, init_height
from io import BytesIO
import time
import concurrent.futures

# got idea from https://medium.com/@jtreitz/the-algorithm-for-a-perfectly-balanced-photo-gallery-914c94a5d8af




def process_image(image: bytes, index: int):
    img = Image.open(BytesIO(image))
    height = img.height
    if height > init_height:
        new_width = int(img.width / height * init_height)
        return img.resize((new_width, init_height), Image.LANCZOS), index
    return img, index


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

    len_img_list = len(img_list)
    total_width = sum([img.width for img in img_list])
    avg_width = total_width / len_img_list
    target_width = avg_width * (len_img_list**0.5)
    num_rows = clamp(int(round(total_width / target_width)), len_img_list)

    if num_rows == 1:
        img_rows = [img_list]
    elif num_rows == len_img_list:
        img_rows = [[img] for img in img_list]
    else:
        aspect_ratios = [int(img.width / img.height * 100) for img in img_list]
        img_rows = linear_partition(aspect_ratios, num_rows, img_list)
        row_widths = [sum([img.width for img in row]) for row in img_rows]
        min_row_width = min(row_widths)
        row_width_ratios = [min_row_width / w for w in row_widths]

        new_img_rows = []
        for row, width_ratio in zip(img_rows, row_width_ratios):
            new_row = []
            for img in row:
                new_width = int(img.width * width_ratio)
                new_height = int(img.height * width_ratio)
                new_row.append(img.resize((new_width, new_height)))

            new_img_rows.append(new_row)

        img_rows = new_img_rows

    row_widths = [sum([img.width for img in row]) for row in img_rows]
    row_heights = [max([img.height for img in row]) for row in img_rows]

    w = ensure_even(min(row_widths))
    h = ensure_even(sum(row_heights))

    result_image = Image.new("RGB", (w, h))
    x_pos, y_pos = 0, 0
    for row in img_rows:
        for img in row:
            result_image.paste(img, (x_pos, y_pos))
            x_pos += img.width
            continue
        y_pos += max([img.height for img in row])
        x_pos = 0
        continue
    return result_image


def make_collage(images: List[bytes], output: str) -> float:
    start = time.time()
    try:
        if len(images) == 1:
            image = Image.open(BytesIO(images[0]))
            save_image(output, image, image.width, image.height)
            return time.time() - start

        pil_images = []
        with concurrent.futures.ThreadPoolExecutor() as executor:
            futures = [
                executor.submit(process_image, image, index)
                for index, image in enumerate(images)
            ]
            for future in concurrent.futures.as_completed(futures):
                pil_images.append(future.result())

        pil_images.sort(key=lambda x: x[1])
        collage = create_collage([img for img, _ in pil_images])

        width = collage.width
        height = collage.height
        if width > width_arg:
            collage = collage.resize(
                (width_arg, int(height / width * width_arg)),
                Image.LANCZOS,
            )
        elif height > height_arg:
            collage = collage.resize(
                (int(width / height * height_arg), height_arg),
                Image.LANCZOS,
            )
        save_image(output, collage, width, height)

        return time.time() - start
    except Exception as e:
        print(f"Error: {e}")
        return -1


if __name__ == "__main__":
    import sys
    import os

    folder_images = sys.argv[1] if len(sys.argv) > 1 else None
    if folder_images:
        files = os.listdir(folder_images)
        files.sort(key=lambda x: int(x.split("_")[1].split(".")[0]))
        images = [Image.open(os.path.join(folder_images, file)) for file in files]
    else:
        images = []
        for i in range(1, 6):
            with open(f"../test_images/{i}.jpg", "rb") as f:
                images.append(f.read())
        print(f"{make_collage(images, 'collage.png')} seconds")
