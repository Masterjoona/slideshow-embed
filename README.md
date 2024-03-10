# tiktok slideshow collage

Sharing tiktoks can be hard on other platforms. [tiktxk](https://github.com/Britmoji/tiktxk) exists but it has just one problem: it cannot embed slideshows well, it only shows first image. So what does this program do? It downloads the images and collages them into a nice single image!

## How to setup

| Env var | Default Value | Description                                   |
|---------------------|---------------|-----------------------------------------------|
| `DOMAIN`            | YOUR DOMAIN   | The domain where it will serve  |
| `PORT`              | 4232          | What port it is on |
| `PROXITOK_INSTANCE` | none          | If you want to use a custom ProxiTok instance |
| `PUBLIC`            | false         | If the index page is public  |
| `LIMIT_PUBLIC_AMOUNT` | 0           | How many collages are linked on the index. 0 is unlimited |
| `FFMPEG`            | ?         | read important below  |
| `FANCY_SLIDESHOW`   | false         | If the slideshow is fancy. Will make a video slideshow instead of a collage with sound like above. This can and will make requests take a lot longer. |


> [!IMPORTANT] 
> This program will check if `/usr/bin/ffmpeg` or `/usr/local/bin/ffmpeg` exists, if it does, it will enable the sound route. You can override this behaviour by using `FFMPEG` environment var. (You can also set this to true, if ffmpeg is not in those paths)
> Since my vps is an arm machine I use an arm build of ffmpeg, you will have to change the files in `Dockerfile` to something from [here](https://johnvansickle.com/ffmpeg/)

> [!IMPORTANT]
> Due to recent frontend change on TikTok's side, slideshow-embed will fetch a [ProxiTok](https://github.com/pablouser1/ProxiTok) public instance to get the audio and image urls as well details for the slideshow.
> <details>
> <summary>What did TikTok change?</summary>
> Previously I was able to get the image urls directly from the page html, but now they are loaded in a separate request and i cannot figure out how to add the necessary signing query parameters to the request. It might take a bit longer to request a collage now.
> </details>

### Dockerfile
Clone this repo and `cd` into it. 


```bash
docker build -t <container_name> .
docker run -d -e DOMAIN='YOUR_DOMAIN_HERE' -p 4232:4232 -v /path/to/your/collages/:/app/collages/ <container_name>
```
*It will error if you don't add a domain*

Basic reverse proxy config for nginx
```nginx
server {
    server_name tt.example.com;
    location / {
        proxy_pass http://localhost:4232;
    }
}
```
### cli
```bash
go build
DOMAIN='YOUR_DOMAIN_HERE' GIN_MODE=release ./meow
```
*it's meow because :3*

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

or youll figure it out if ou dont use docker
</details>



## What does it look like?


https://github.com/Masterjoona/slideshow-embed/assets/69722179/cb07845d-851d-4cc1-97ed-badf80c37faa

*Yes it is kinda slow but I really cannot affect that*
*old video*


## notes
this is a beginner project so there might some insane design choices ![trolley](https://cdn.discordapp.com/emojis/1068825486265942056.webp?size=48&name=trolley&quality=lossless) 

One such example is calling [this](https://github.com/twilsonco/PyPhotoCollage) python script to make the collages

Massive thanks to [Dmitry Alimov](https://github.com/delimitry) and [Tim Wilson](https://github.com/twilsonco) for making the insane maths for the collaging script!!!

## contributing

be my guest
