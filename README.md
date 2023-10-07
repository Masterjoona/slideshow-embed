# tiktok slideshow collage

Sharing tiktoks can be hard on other platforms. [tiktxk](https://github.com/Britmoji/tiktxk) exists but it has just one problem: it cannot embed slideshows well, it only shows first image. So what does this program do? It downloads the images and collages them into a nice single image!

## How to setup

### Dockerfile
Clone this repo and `cd` into it change `YOUR_DOMAIN_HERE` in Dockerfile to some domain. You can also change `false` to `true` if you want to have links on the index page.

```bash
docker build -t <your_image_name> .
docker run -d -p 4232:4232 -v <path>:/app/collages/ <container_name>
```
Where path can be accessed by the domain you specified.

### cli
You'll figure it out. just use docker

## notes
this is a beginner project so there might some insane design choices ![trolley](https://cdn.discordapp.com/emojis/1068825486265942056.webp?size=48&name=trolley&quality=lossless) 
One such example is calling [this](https://github.com/delimitry/collage_maker) python script to make the collages

## contributing

be my guest
