FROM klakegg/hugo:0.111.3-ubuntu AS builder

WORKDIR /src
COPY . .

RUN ls -al /src

RUN hugo --config hugo.toml

FROM nginx:alpine
COPY --from=builder /src/public /usr/share/nginx/html