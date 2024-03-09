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
import re
import argparse
import os
import random
import math
from PIL import ImageOps
from PIL import Image
from typing import List
import io
import time

fpng_py_imported = False
try:
    import fpng_py

    fpng_py_imported = True
except ImportError:
    pass

# got idea from https://medium.com/@jtreitz/the-algorithm-for-a-perfectly-balanced-photo-gallery-914c94a5d8af


# start partition problem algorithm from https://stackoverflow.com/a/7942946
# modified to act on list of images rather than the weights themselves
# more info on the partition problem http://www8.cs.umu.se/kurser/TDBAfl/VT06/algorithms/BOOK/BOOK2/NODE45.HTM


def linear_partition(seq, k, dataList=None):
    if k <= 0:
        return []
    n = len(seq) - 1
    if k > n:
        return map(lambda x: [x], seq)
    _, solution = linear_partition_table(seq, k)
    k, ans = k - 2, []
    if dataList == None or len(dataList) != len(seq):
        while k >= 0:
            ans = [[seq[i] for i in range(solution[n - 1][k] + 1, n + 1)]] + ans
            n, k = solution[n - 1][k], k - 1
        ans = [[seq[i] for i in range(0, n + 1)]] + ans
    else:
        while k >= 0:
            ans = [[dataList[i] for i in range(solution[n - 1][k] + 1, n + 1)]] + ans
            n, k = solution[n - 1][k], k - 1
        ans = [[dataList[i] for i in range(0, n + 1)]] + ans
    return ans


def linear_partition_table(seq, k):
    n = len(seq)
    table = [[0] * k for x in range(n)]
    solution = [[0] * (k - 1) for x in range(n - 1)]
    for i in range(n):
        table[i][0] = seq[i] + (table[i - 1][0] if i else 0)
    for j in range(k):
        table[0][j] = seq[0]
    for i in range(1, n):
        for j in range(1, k):
            table[i][j], solution[i - 1][j - 1] = min(
                (
                    (max(table[x][j - 1], table[i][0] - table[x][0]), x)
                    for x in range(i)
                ),
                key=lambda t: t[0],
            )
    return (table, solution)


# end partition problem algorithm


def clamp(v, l, h):
    return l if v < l else h if v > h else v


# takes list of PIL image objects and returns the collage as a PIL image object
def makeCollage(
    imgList, spacing=0, antialias=False, background=(0, 0, 0), aspectratiofactor=1.0
):
    # first downscale all images according to the minimum height of any image
    #     minHeight = min([img.height for img in imgList])
    #     if antialias:
    #         imgList = [img.resize((int(img.width / img.height * minHeight),minHeight), Image.LANCZOS) if img.height > minHeight else img for img in imgList]
    #     else:
    #         imgList = [img.resize((int(img.width / img.height * minHeight),minHeight)) if img.height > minHeight else img for img in imgList]

    # first upscale all images according to the maximum height of any image (downscaling would result in a terrible quality image if a very short image was included in the batch)
    maxHeight = max([img.height for img in imgList])
    if antialias:
        imgList = [
            (
                img.resize(
                    (int(img.width / img.height * maxHeight), maxHeight), Image.LANCZOS
                )
                if img.height < maxHeight
                else img
            )
            for img in imgList
        ]
    else:
        imgList = [
            (
                img.resize((int(img.width / img.height * maxHeight), maxHeight))
                if img.height < maxHeight
                else img
            )
            for img in imgList
        ]
    # generate the input for the partition problem algorithm
    # need list of aspect ratios and number of rows (partitions)
    # imgHeights = [img.height for img in imgList]
    totalWidth = sum([img.width for img in imgList])
    avgWidth = totalWidth / len(imgList)
    targetWidth = avgWidth * math.sqrt(len(imgList) * aspectratiofactor)

    numRows = clamp(int(round(totalWidth / targetWidth)), 1, len(imgList))
    if numRows == 1:
        imgRows = [imgList]
    elif numRows == len(imgList):
        imgRows = [[img] for img in imgList]
    else:
        aspectRatios = [int(img.width / img.height * 100) for img in imgList]

        # get nested list of images (each sublist is a row in the collage)
        imgRows = linear_partition(aspectRatios, numRows, imgList)

        # scale down larger rows to match the minimum row width
        rowWidths = [
            sum([img.width + spacing for img in row]) - spacing for row in imgRows
        ]
        minRowWidth = min(rowWidths)
        rowWidthRatios = [minRowWidth / w for w in rowWidths]
        if antialias:
            imgRows = [
                [
                    img.resize(
                        (int(img.width * widthRatio), int(img.height * widthRatio)),
                        Image.LANCZOS,
                    )
                    for img in row
                ]
                for row, widthRatio in zip(imgRows, rowWidthRatios)
            ]
        else:
            imgRows = [
                [
                    img.resize(
                        (int(img.width * widthRatio), int(img.height * widthRatio))
                    )
                    for img in row
                ]
                for row, widthRatio in zip(imgRows, rowWidthRatios)
            ]

    # pupulate new image
    rowWidths = [sum([img.width + spacing for img in row]) - spacing for row in imgRows]
    rowHeights = [max([img.height for img in row]) for row in imgRows]
    minRowWidth = min(rowWidths)
    w, h = (minRowWidth, sum(rowHeights) + spacing * (numRows - 1))

    if w % 2 != 0:
        w += 1
    if h % 2 != 0:
        h += 1

    if background == (0, 0, 0):
        background += tuple([0])
    else:
        background += tuple([255])

    outImg = Image.new("RGBA", (w, h), background)
    xPos, yPos = (0, 0)

    for row in imgRows:
        for img in row:
            outImg.paste(img, (xPos, yPos))
            xPos += img.width + spacing
            continue
        yPos += max([img.height for img in row]) + spacing
        xPos = 0
        continue
    return outImg


# modified (significantly) from https://github.com/delimitry/collage_maker
# this main function is for the CLI implementation
def main():
    def rgb(s):
        try:
            rgb = (
                0 if v < 0 else 255 if v > 255 else v for v in map(int, s.split(","))
            )
            return rgb
        except:
            raise argparse.ArgumentTypeError(
                'Background must be (r,g,b) --> "(0,0,0)" to "(255,255,255)"'
            )

    parse = argparse.ArgumentParser(description="Photo collage maker")
    parse.add_argument(
        "-f",
        "--folder",
        dest="folder",
        help="folder with images (*.jpg, *.jpeg, *.png)",
        default=False,
    )
    parse.add_argument(
        "-F",
        "--file",
        dest="file",
        help="file with newline separated list of files",
        default=False,
    )
    parse.add_argument(
        "-o",
        "--output",
        dest="output",
        help="output collage image filename",
        default="collage.png",
    )
    parse.add_argument(
        "-W",
        "--width",
        dest="width",
        type=int,
        help="resulting collage image height (mutually exclusive with --height)",
        default=5000,
    )
    parse.add_argument(
        "-H",
        "--height",
        dest="height",
        type=int,
        help="resulting collage image height (mutually exclusive with --width)",
        default=5000,
    )
    parse.add_argument(
        "-i",
        "--initheight",
        dest="initheight",
        type=int,
        help="resize images on input to set height",
        default=500,
    )
    parse.add_argument(
        "-s",
        "--shuffle",
        action="store_true",
        dest="shuffle",
        help="enable images shuffle",
    )
    parse.add_argument(
        "-g",
        "--gap-between-images",
        dest="imagegap",
        type=int,
        help="number of pixels of transparent space (if saving as png file; otherwise black or specified background color) to add between neighboring images",
        default=3,
    )
    parse.add_argument(
        "-b",
        "--background-color",
        dest="background",
        type=rgb,
        help="color (r,g,b) to use for background if spacing is added between images",
        default=(0, 0, 0),
    )
    parse.add_argument(
        "-c",
        "--count",
        dest="count",
        type=int,
        help="count of images to use",
        default=0,
    )
    parse.add_argument(
        "-r",
        "--scale-aspect-ratio",
        dest="aspectratiofactor",
        type=float,
        help="aspect ratio scaling factor, multiplied by the average aspect ratio of the input images to determine the output aspect ratio",
        default=1.0,
    )
    parse.add_argument(
        "-a",
        "--no-antialias-when-resizing",
        dest="noantialias",
        action="store_false",
        help="disable antialiasing on intermediate resizing of images (runs faster but output image looks worse; final resize is always antialiased)",
    )
    parse.add_argument(
        "-resize",
        "--resize-images",
        dest="resize",
        action="store_true",
        help="resize images to the same size",
    )
    parse.add_argument("files", nargs="*")

    args = parse.parse_args()

    if args.resize:
        resize_images(args.folder)
        return

    if not args.file and not args.folder and not args.files:
        parse.print_help()
        exit(1)

    # get images
    images = args.files
    if args.folder:
        files = [os.path.join(args.folder, fn) for fn in os.listdir(args.folder)]
        images = [
            fn
            for fn in files
            if os.path.splitext(fn)[1].lower() in (".jpg", ".jpeg", ".png")
        ]

    if args.file:
        images = []
        with open(args.file, "r") as f:
            for line in f:
                images.append(line.strip())
    elif args.folder:
        images = []
        for root, _, files in os.walk(args.folder):
            for name in files:
                if re.findall("jpg|png|jpeg", name.split(".")[-1]):
                    fname = os.path.join(root, name)
                    images.append(fname)

    """
    Doesn't seem to matter?
    if len(images) < 3:
        print("Need to use 3 or more images. Try again")
        return
    """

    if len(images) == 1:
        img = Image.open(images[0])
        img.save(args.output)
        print(f"Collage is ready at {args.output}!")
        return

    # shuffle images if needed
    if args.shuffle:
        random.shuffle(images)
    else:
        try:
            images = sorted(images, key=lambda x: int(x.split("/")[1].split(".jpg")[0]))
        except:
            pass
    if args.count > 2:
        images = images[: args.count]

    # get PIL image objects for all the photos
    print("Loading photos...")
    pilImages = []
    sum_time = 0
    for f in images:

        img = Image.open(f)

        exif = img.getexif()
        # Remove unwanted EXIF tags
        for k in exif.keys():
            if k != 0x0112:
                exif[k] = (
                    None  # If I don't set it to None first (or print it) the del fails for some reason.
                )
            del exif[k]

        img.info["exif"] = exif.tobytes()
        # Rotate the image based on EXIF orientation data
        ImageOps.exif_transpose(img)
        # Calculate the initial height only once
        initial_height = args.initheight

        # Resize the image if needed
        if initial_height > 2 and img.height > initial_height:
            # Choose the appropriate resampling method based on 'noantialias' option
            resampling_method = Image.NEAREST if args.noantialias else Image.LANCZOS
            new_width = int(img.width / img.height * initial_height)
            pilImages.append(img.resize((new_width, initial_height), resampling_method))
        else:
            pilImages.append(img)

    print("Making collage...")
    st = time.time()
    collage = makeCollage(
        pilImages,
        args.imagegap,
        not args.noantialias,
        args.background,
        args.aspectratiofactor,
    )
    print(f"Collage took {time.time() - st} seconds")
    if args.width > 0 and collage.width > args.width:
        collage = collage.resize(
            (args.width, int(collage.height / collage.width * args.width)),
            Image.LANCZOS,
        )
        pass
    elif args.height > 0 and collage.height > args.height:
        collage = collage.resize(
            (int(collage.width / collage.height * args.height), args.height),
            Image.LANCZOS,
        )
        pass
    output = args.output

    if fpng_py_imported:
        fpng_time = time.time()
        fpng_py.fpng_encode_image_to_file(
            output, collage.tobytes(), collage.width, collage.height
        )
        print(f"fpng took {time.time() - fpng_time} seconds")
        # PIL for 22 images 1.4 seconds, fpng 0.068 seconds
    else:
        collage.save(output)

    print(f"Collage is ready at {output}!")


class SlideshowImageDimensions:
    def __init__(self):
        self.maxWidth = 0
        self.maxHeight = 0


def getImageDimensions(buffers: List[bytes]) -> SlideshowImageDimensions:
    dimensions = SlideshowImageDimensions()

    for buffer in buffers:
        image = Image.open(io.BytesIO(buffer))
        width, height = image.size

        if not width or not height:
            raise Exception("Could not get image dimensions")

        if width > dimensions.maxWidth:
            dimensions.maxWidth = width

        if height > dimensions.maxHeight:
            dimensions.maxHeight = height
    return dimensions


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


def resize_contain(image, size, resample=Image.LANCZOS, bg_color=(255, 255, 255, 0)):
    """
    Resize image according to size.
    image:      a Pillow image instance
    size:       a list of two integers [width, height]
    """
    img_format = image.format
    img = image.copy()
    img.thumbnail((size[0], size[1]), resample)
    background = Image.new("RGBA", (size[0], size[1]), bg_color)
    img_position = (
        int(math.ceil((size[0] - img.size[0]) / 2)),
        int(math.ceil((size[1] - img.size[1]) / 2)),
    )
    background.paste(img, img_position)
    background.format = img_format
    return background.convert("RGBA")


def resize_images(path):
    # get amount of images
    images = []
    for root, _, files in os.walk(path):
        for name in files:
            if re.findall("jpg|png|jpeg", name.split(".")[-1]):
                fname = os.path.join(root, name)
                images.append(fname)
    width, heigth = 0, 0
    images_buffer = []
    for i in range(1, len(images) + 1):
        with open(f"{path}/{i}.jpg", "rb") as f:
            images_buffer.append(f.read())
    dimensions = getImageDimensions(images_buffer)
    width = dimensions.maxWidth
    heigth = dimensions.maxHeight

    if width % 2 != 0:
        width += 1
    if heigth % 2 != 0:
        heigth += 1

    for i in range(1, len(images) + 1):
        with open(f"{path}/{i}.jpg", "r+b") as f:
            with Image.open(f) as image:
                cover = resize_contain(
                    image,
                    [width, heigth],
                    resample=Image.LANCZOS,
                    bg_color=(0, 0, 0),
                )
                cover.convert("RGB").save(f"{path}/{i}.jpg", "JPEG")


if __name__ == "__main__":
    main()
