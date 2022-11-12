# pull the latest official nginx image
FROM nginx:stable
# run docker service on HTTP
EXPOSE 80
# copy static maintanence
COPY /html/welcome.html /usr/share/nginx/html/index.html
# copy favicon
COPY /assets/favicon.ico /usr/share/nginx/html/favicon.ico
# copy images
COPY /assets/prof-pic.png /usr/share/nginx/html/prof-pic.png
COPY /assets/github_logo.png /usr/share/nginx/html/github_logo.png
COPY /assets/linkedin_logo.png /usr/share/nginx/html/linkedin_logo.png 
# STOPSIGNAL SIGQUIT
CMD ["nginx", "-g", "daemon off;"]
