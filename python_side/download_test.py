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


async def fetch_html(url:str):
    print(f"Fetching {url}")
    async with aiohttp.ClientSession() as session:
        async with session.get(url) as response:
            text = await response.text()
            text = text.split('<script id="__UNIVERSAL_DATA_FOR_REHYDRATION__" type="application/json">')[1].split('</script>')[0]
            return json.loads(text)


async def main():
    if not aweme_id:
        print("Please provide an video id")
        return
    script_data = await fetch_html(place_holder_url + aweme_id)
    images = script_data["__DEFAULT_SCOPE__"]["webapp.video-detail"]["itemInfo"]["itemStruct"]["imagePost"]["images"]
    tasks = [download_image(image) for image in images]
    buffers = await asyncio.gather(*tasks)
    for i, buffer in enumerate(buffers):
        with open(f"{output}/{i+1}.jpg", "wb") as f:
            f.write(buffer)


asyncio.run(main())
