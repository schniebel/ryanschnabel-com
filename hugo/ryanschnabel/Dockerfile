FROM klakegg/hugo:0.89.4 AS builder

WORKDIR /src
COPY . .

RUN hugo --config hugo.toml

FROM nginx:alpine
COPY --from=builder /src/public /usr/share/nginx/html