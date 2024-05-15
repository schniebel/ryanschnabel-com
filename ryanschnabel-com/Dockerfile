# Using a multi-stage build

# Stage 1: Building the site with Hugo
FROM debian:buster as STAGEONE

# Install Hugo for amd64
ENV HUGO_VERSION=0.120.2
ADD https://github.com/gohugoio/hugo/releases/download/v${HUGO_VERSION}/hugo_extended_${HUGO_VERSION}_Linux-64bit.tar.gz /tmp
RUN tar -xf /tmp/hugo_extended_${HUGO_VERSION}_Linux-64bit.tar.gz -C /usr/local/bin/

# Install syntax highlighting
RUN apt-get update && \
    apt-get install -y python3-pygments

# Set working directory and copy the source code
WORKDIR /src
COPY . .

# Build the site
RUN hugo --config hugo.toml --logLevel info

# Stage 2: Serving the site with Nginx
FROM nginx:1.24-perl
COPY --from=STAGEONE /src/public /usr/share/nginx/html