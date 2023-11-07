FROM debian:buster as STAGEONE

# install hugo
ENV HUGO_VERSION=0.120.2
ADD https://github.com/gohugoio/hugo/releases/download/v${HUGO_VERSION}/hugo_extended_${HUGO_VERSION}_linux-arm64.tar.gz /tmp
RUN tar -xf /tmp/hugo_extended_${HUGO_VERSION}_linux-arm64.tar.gz -C /usr/local/bin/

# install syntax highlighting
RUN apt-get update
RUN apt-get install -y python3-pygments

# build site
WORKDIR /src
COPY . .
RUN hugo --config hugo.toml --logLevel info

FROM arm64v8/nginx:1.24-perl
COPY --from=STAGEONE /src/public /usr/share/nginx/html