# tiktok slideshow collage

Sharing tiktoks can be hard on other platforms. [tiktxk](https://github.com/Britmoji/tiktxk) exists but it has just one problem: it cannot embed slideshows well, it only shows first image. So what does this program do? It downloads the images and collages them into a nice single image!

Embedding normal tiktoks is also supported.


## How to setup

| Env var | Default Value | Description                                   |
|---------------------|---------------|-----------------------------------------------|
| `DOMAIN`            | ""   | The domain where it will serve  |
| `PORT`              | 4232          | What port it is on |
| `PUBLIC`            | false         | If the index page is public  |
| `LIMIT_PUBLIC_AMOUNT` | 0           | How many collages are linked on the index. 0 is unlimited |
| `FFMPEG`            | false         | collages with sound  |
| `FANCY_SLIDESHOW`   | false         | will make a video slideshow instead of a collage with sound like above. This can and will make requests take a lot longer. |
| `INSTALL_ID`      | "" | if you have acquired a tiktok install id, you can put it here. |


### Docker compose
Clone this repo and `cd` into it. edit the compose file for path for the collages. You can also the uncommenn the build tag if you want to scrape ttsave.app instead of relying on the api endpoint.

```bash
docker compose build && docker compose up -d
```

### cli
you probably dont want to do this

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

*Yes it is kinda slow but I really cannot affect that*


| url path | description                            |
|----------|----------------------------------------|
| /t?v=    | normal collage or embed a video tiktok | 
| /s?v=    | collage with sound, i guess you can embed a video tiktok as well...                     |
| /f?v=    | slides the images, same for this lol                      |

## notes
this is a beginner project so there might some insane design choices ![trolley](https://cdn.discordapp.com/emojis/1068825486265942056.webp?size=48&name=trolley&quality=lossless) 


One such example is calling [this](https://github.com/twilsonco/PyPhotoCollage) python script to make the collages. actually it doenst call it anymore, there are two containers, one for http server and a python one for collaging and resizing.

Massive thanks to [Dmitry Alimov](https://github.com/delimitry) and [Tim Wilson](https://github.com/twilsonco) for inventing the collaging script!

## contributing

be my guest
