# tiktok slideshow collage

Sharing tiktoks can be hard on other platforms. [tiktxk](https://github.com/Britmoji/tiktxk) exists but it has just one problem: it cannot embed slideshows well, it only shows first image. So what does this program do? It downloads the images and collages them into a nice single image!

## How to setup
> [!NOTE]
> Domain, public and port are all configurable via environment variables. Only domain is required. Public will show link on the index page, by default is false. Port is 4232 by default.
### Dockerfile
Clone this repo and `cd` into it. 
Changing the false to true will install ffmpeg and allow you to use the video collages. (will make your container bigger)

```bash
docker build --build-arg INSTALL_FFMPEG=false -t <container_name> .
docker run -d -e DOMAIN='YOUR_DOMAIN_HERE' -p 4232:4232 -v /path/to/your/collages/:/app/collages/ <container_name>
```

Basic reverse proxy config for nginx
```nginx
server {
    server_name tt.example.com;

    location / {
        proxy_pass http://localhost:4232;
        root /path/to/your/collages;
    }
}
```
### cli
```bash
DOMAIN='YOUR_DOMAIN_HERE' go run main.go
```

## notes
this is a beginner project so there might some insane design choices ![trolley](https://cdn.discordapp.com/emojis/1068825486265942056.webp?size=48&name=trolley&quality=lossless) 
One such example is calling [this](https://github.com/delimitry/collage_maker) python script to make the collages

## contributing

be my guest
