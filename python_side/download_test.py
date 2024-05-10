import aiohttp
import asyncio
import sys
import os

install_id = sys.argv[1] if len(sys.argv) > 1 else None
tk_id = sys.argv[2] if len(sys.argv) > 2 else None
output = sys.argv[3] if len(sys.argv) > 3 else None

if tk_id is None or install_id is None or output is None:
    print("Usage: download_test.py <install_id> <aweme_id> <output>")
    sys.exit(1)

tk_api_url = f"https://api22-normal-c-alisg.tiktokv.com/aweme/v1/feed/?iid={install_id}&device_id=7351044760062330401&channel=googleplay&app_name=musical_ly&version_code=300904&device_platform=android&device_type=SM-ASUS_Z01QD&os_version=9&aweme_id={tk_id}"


async def download_image(image_dict) -> bytes:
    url = image_dict["display_image"]["url_list"][1]
    async with aiohttp.ClientSession() as session:
        async with session.get(url) as response:
            return await response.read()


async def fetch_api(url: str) -> dict:
    print(f"Fetching {url}")
    async with aiohttp.ClientSession() as session:
        async with session.get(url) as response:
            return await response.json()


async def main():
    api_data = await fetch_api(tk_api_url)
    images = api_data["aweme_list"][0]["image_post_info"]["images"]
    tasks = [download_image(image) for image in images]
    buffers = await asyncio.gather(*tasks)
    if not os.path.exists(output):
        os.makedirs(output)
    for i, buffer in enumerate(buffers):
        with open(f"{output}/{i+1}.jpg", "wb") as f:
            f.write(buffer)


asyncio.run(main())
