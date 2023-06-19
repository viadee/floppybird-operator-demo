# fast and small python base image alpine
FROM python:3.9.17-alpine

RUN adduser -D floppybird-server
USER floppybird-server

WORKDIR /opt/
COPY . .

# define env var FLOPPY_DISTRO which can be set before entrypoint is called
ENV FLOPPY_DISTRO="floppybird"
# other FLOPPY_DISTRO: "flappy-docker", "flappy-math-saga", "flappydragon", 
CMD echo "Available distros:" && ls && echo "Choose distro: ${FLOPPY_DISTRO}" && python3 -m http.server 8000 --directory ./${FLOPPY_DISTRO}

EXPOSE 8000
