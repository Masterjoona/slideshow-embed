# tiktok slideshow collage

Sharing tiktoks can be hard on other platforms. [tiktxk](https://github.com/Britmoji/tiktxk) exists but it has just one problem: it cannot embed slideshows well, it only shows first couple images. So what does this program do? It downloads the images and collages them into a nice single image!

Embedding video tiktoks is also supported.

## How to setup

| Env var               | Default Value | Description                                                                                                       |
| --------------------- | ------------- | ----------------------------------------------------------------------------------------------------------------- |
| `DOMAIN`              | none          | The domain where it will serve                                                                                    |
| `PORT`                | 4232          | The port where it will serve                                                                                      |
| `PUBLIC`              | false         | If the index page is public                                                                                       |
| `LIMIT_PUBLIC_AMOUNT` | -1            | How many links are linked on the index. -1 is unlimited. set to 0 and public to true to only display stats        |
| `FFMPEG`              | false         | collages with sound                                                                                               |
| `FANCY_SLIDESHOW`     | false         | enables video slideshow path, it slides the images with audio. This can and will make requests take a lot longer. |
| `DOWNLOADER`          | false         | enables the downloader path.                                                                                      |
| `TIKTOK_PROVIDER`     | api           | Set the provider for data, either tikwm or ttsave otherwise it will default to tiktoks own api.                   |
| `SUBTITLER`           | false         | burn the auto generated subtitles into the video.                                                                 |
| `TEST_PROVIDERS`      | false         | enables the test providers path.                                                                                  |

### Docker compose

Clone this repo and `cd` into it. edit the compose file for path for the collages. Leave the build arg as tikwm or change it to ttsave. I would have to proxy the videos if i were to use urls in tiktoks html.

```bash
docker compose up --build -d
```

### cli

you probably dont want to do this since you need to run 2 things at once.

<details>
<summary> add fpng_py for fast png encoding</summary>
Add this to the dockerfile

```Dockerfile
RUN git clone --recurse-submodules https://github.com/K0lb3/fpng_py
# for arm64 we disable some build args. what does these flags do? i dont know
RUN sed -i 's/"-msse4.1"/#&/' fpng_py/setup.py
RUN sed -i 's/"-mpclmul"/#&/' fpng_py/setup.py
WORKDIR /app/fpng_py
RUN pip install . --break-system-packages

```

or you can build it yourself

```bash
git clone --recurse-submodules https://github.com/K0lb3/fpng_py
cd fpng_py
pip install .
```

and copy the compiled files to the container

```Dockerfile
COPY ./fpng_py/build/lib.path/fpng_py /app/fpng_py
```

</details>

## What does it look like?

https://github.com/Masterjoona/slideshow-embed/assets/69722179/cb07845d-851d-4cc1-97ed-badf80c37faa

_Yes it is kinda slow but I really cannot affect that_

| url path               | description                                                                                                                                                                     |
| ---------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| /t?v=                  | normal collage or embed a video tiktok                                                                                                                                          |
| /s?v=                  | collage with sound, i guess you can embed a video tiktok as well...                                                                                                             |
| /f?v=                  | slides the images, same for this lol. The requests to this will take a long time to complete, so its better to open this in a browser and when it completes send it to discord. |
| /d?v=                  | show all images and a button download a zip of all images.                                                                                                                      |
| /subs?v=link&lang=lang | burn the subtitles into a video ([language codes](https://www.tiktok.com/player/node/api/translation_language), maybe more accepted as well?)                                   |
| /test?v=               | test the providers, this will show the data that the providers give.                                                                                                            |

## notes

this is a beginner project so there might some insane design choices ![trolley](https://cdn.discordapp.com/emojis/1068825486265942056.webp?size=48&name=trolley&quality=lossless)

One such example is calling [this](https://github.com/twilsonco/PyPhotoCollage) python script to make the collages. Its refactored to be a bit more readable and to have types and removed all the unnecessary stuff.

Massive thanks to [Dmitry Alimov](https://github.com/delimitry) and [Tim Wilson](https://github.com/twilsonco) for inventing the collaging script!

## contributing

be my guest
