# Floppybird image

The code base is copied on 19.06.2023 from the following Repositories:

- floppybird: <https://github.com/nebez/floppybird>
- flappy-docker: <https://github.com/mikesir87/flappy-dock>
- flappy-math-saga: <https://github.com/tikwid/flappy-math-saga>
- flappydragon: <https://github.com/iarunava/flappydragon>

Many thanks to the contributor

To disable caching the following was added to index.html files:

```html
      <!-- no caching for demo-->>
      <meta http-equiv="Cache-control" content="no-cache, no-store, must-revalidate">
      <meta http-equiv="Pragma" content="no-cache">
      <meta http-equiv="expires" content="0">
```

## start

```shell
docker run -p 8000:8000 --rm --name floppybird --env FLOPPY_DISTRO=floppybird -it crowdsalat/floppybird-demo:latest
docker run -p 8000:8000 --rm --name floppybird --env FLOPPY_DISTRO=flappy-docker -it crowdsalat/floppybird-demo:latest
docker run -p 8000:8000 --rm --name floppybird --env FLOPPY_DISTRO=flappydragon -it crowdsalat/floppybird-demo:latest
docker run -p 8000:8000 --rm --name floppybird --env FLOPPY_DISTRO=flappy-math-saga -it crowdsalat/floppybird-demo:latest
```

## localbuild and push

```shell
docker login 
docker-buildx build --platform linux/amd64,linux/arm64 --push -t crowdsalat/floppybird-demo:0.0.4 -t crowdsalat/floppybird-demo:latest .
```
