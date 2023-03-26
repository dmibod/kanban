docker run --name proxy -d -p 80:80 -v "./proxy/conf:/etc/nginx/conf.d" -v "./proxy/www:/usr/share/nginx/html:ro" nginx
