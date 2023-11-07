# install hugo
ENV HUGO_VERSION=0.41
ADD https://github.com/gohugoio/hugo/releases/download/v${HUGO_VERSION}/hugo_${HUGO_VERSION}_Linux-64bit.tar.gz /tmp/
RUN tar -xf /tmp/hugo_${HUGO_VERSION}_Linux-64bit.tar.gz -C /usr/local/bin/

# install syntax highlighting
RUN apt-get update
RUN apt-get install -y python3-pygments

# build site
COPY source /source
RUN hugo --source=/source/ --destination=/public/

# nginx
FROM arm64v8/nginx:1.24-perl
RUN apk --update add curl bash
RUN rm /etc/nginx/conf.d/default.conf
COPY modules/nginx.luiscachog.io.conf /etc/nginx/conf.d/luiscachog.io.conf
COPY --from=STAGEONE /public/ /usr/share/nginx/html/
EXPOSE 80 443