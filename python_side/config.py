from PIL.Image import Image as ImageType
from PIL import Image
from typing import List, Tuple
import io
import os.path


def save_with_fpng(
    filepath: str, image: ImageType, width: int, height: int, resize: bool
) -> None:
    fpng_py.fpng_encode_image_to_file(filepath, image.tobytes(), width, height)  # type: ignore


def save_with_pillow(
    filepath: str, image: ImageType, width: int, height: int, resize: bool
) -> None:
    if resize:
        image = image.resize((width, height), Image.LANCZOS)
    image.save(filepath)


def find_largest_image_dimensions_pillow(buffers: List[bytes]) -> Tuple[int, int]:
    max_width = 0
    max_height = 0

    for buffer in buffers:
        image = Image.open(io.BytesIO(buffer))
        width, height = image.size

        if width <= 0 or height <= 0:
            raise ValueError("Invalid image dimensions")

        max_width = max(max_width, width)
        max_height = max(max_height, height)

    return max_width, max_height


def find_largest_image_dimensions_fpng_py(buffers: List[bytes]) -> Tuple[int, int]:
    max_width = 0
    max_height = 0

    for buffer in buffers:
        width, height, _ = fpng_py.fpng_get_info(buffer)  # type: ignore
        if width <= 0 or height <= 0:
            raise ValueError("Invalid image dimensions")

        max_width = max(max_width, width)
        max_height = max(max_height, height)

    return max_width, max_height


has_fpng_module = False
save_image = save_with_pillow
find_largest_image_dimensions = find_largest_image_dimensions_pillow
try:
    import fpng_py

    save_image = save_with_fpng
    # find_largest_image_dimensions = find_largest_image_dimensions_fpng_py
    # mfw it exits but he doesnt export the function https://github.com/K0lb3/fpng_py/blob/a826a55f682cead912d05eb0bad52f857b4bc726/fpng_py/fpng_py.cpp#L200
except ImportError:
    pass

width_arg = 5000
height_arg = 5000
init_height = 500
is_docker = os.path.isfile("/.dockerenv")
collages_path = "/app/collages" if is_docker else "../collages"
temporary_dir = "/tmp/collages" if is_docker else "../tmp/collages"
