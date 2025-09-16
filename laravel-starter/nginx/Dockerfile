FROM public.ecr.aws/nginx/nginx:1.24.0

COPY ./nginx.conf /etc/nginx/
COPY ./default.conf /etc/nginx/conf.d/

WORKDIR /usr/share/nginx/html

CMD ["nginx", "-g", "daemon off;"]