FROM busybox  
ADD index.html /www/index.html 
EXPOSE 9999 
CMD httpd -p 9999 -h /www ; tail -f /dev/null
