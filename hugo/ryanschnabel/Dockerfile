FROM klakegg/hugo:0.111.3-ext-ubuntu-onbuild AS builder

WORKDIR /src
COPY . .

RUN hugo --config hugo.toml

FROM arm64v8/nginx:1.24-perl
COPY --from=builder /src/public /usr/share/nginx/html