FROM debian:buster as STAGEONE

# install hugo
ENV HUGO_VERSION=0.120.3
ADD https://github.com/gohugoio/hugo/releases/download/v${HUGO_VERSION}/hugo_${HUGO_VERSION}_linux-arm64.tar.gz /tmp
RUN tar -xf /tmp/hugo_${HUGO_VERSION}_linux-arm64.tar.gz -C /usr/local/bin/

# install syntax highlighting
RUN apt-get update
RUN apt-get install -y python3-pygments

# build site
RUN hugo --config hugo.toml

#COPY source /source
#RUN hugo --source=/source/ --destination=/public/

#RUN ls -la /source/public

# nginx
# FROM arm64v8/nginx:1.24-perl
# RUN apk --update add curl bash
# RUN rm /etc/nginx/conf.d/default.conf
# COPY modules/nginx.luiscachog.io.conf /etc/nginx/conf.d/luiscachog.io.conf
# COPY --from=STAGEONE /public/ /usr/share/nginx/html/
# EXPOSE 80 443