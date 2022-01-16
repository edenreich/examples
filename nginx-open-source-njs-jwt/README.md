## NGINX Open-Source with JWT Authentication

An example for NGINX Open-Source with JWT Authentication using built-in NJS crypto module.

Run the example with:

```sh
docker run --rm -it \
  -v ${PWD}/html:/var/www/html \
  -v ${PWD}/conf.d:/etc/nginx/conf.d \
  -v ${PWD}/njs:/etc/nginx/njs \
  -v ${PWD}/nginx.conf:/etc/nginx/nginx.conf \
  -v ${PWD}/data:/etc/nginx/data \
  -w /etc/nginx \
  -p 80:80 \
  --name nginx \
  nginx:alpine
```

Create a user with the util:

```sh
./utils/add_user.js
```

This example is not considering external user's data store, for simplicity purposes.
