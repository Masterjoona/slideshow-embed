import aiohttp
import asyncio
import sys
import json

aweme_id = sys.argv[1] if len(sys.argv) > 1 else None
output = sys.argv[2] if len(sys.argv) > 2 else None
place_holder_url = "https://www.tiktok.com/@placeholder/video/"


async def download_image(image_dict) -> bytes:
    url = image_dict["imageURL"]["urlList"][0]
    async with aiohttp.ClientSession() as session:
        async with session.get(url) as response:
            return await response.read()


async def fetch_api(aweme_id: str):
    async with aiohttp.ClientSession() as session:
        async with session.options(
            f"https://api22-normal-c-alisg.tiktokv.com/aweme/v1/feed/?aweme_id={aweme_id}"
        ) as response:
            if response.status != 200:
                print("Failed to fetch")
                return
            return await response.json()


async def main():
    if not aweme_id:
        print("Please provide an video id")
        return
    resp = await fetch_api(aweme_id)
    with open("resp.json", "w") as f:
        json.dump(resp["aweme_list"][0], f)


asyncio.run(main())
