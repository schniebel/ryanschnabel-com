# pull the latest official nginx image
FROM nginx:stable
# run docker service on HTTP
EXPOSE 80
# copy static maintanence
COPY html/welcome.html /usr/share/nginx/html/index.html
# copy favicon
COPY assets/favicon.ico /usr/share/nginx/html/favicon.ico
# STOPSIGNAL SIGQUIT
CMD ["nginx", "-g", "daemon off;"]
