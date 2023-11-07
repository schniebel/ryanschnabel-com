FROM klakegg/hugo:0.89.4 AS builder

WORKDIR /src
COPY . .

RUN hugo --config hugo.toml

RUN ls -la /src/public

FROM arm64v8/nginx:1.24-perl
COPY --from=builder /src/public /usr/share/nginx/html
