FROM klakegg/hugo:0.89.4 AS builder

WORKDIR /src
COPY . .

RUN ls -al /src

RUN hugo

FROM nginx:alpine
COPY --from=builder /src/public /usr/share/nginx/html