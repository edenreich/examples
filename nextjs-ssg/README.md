
## NextJS Static Site Generator(SSG)

Short example on how to develop a NextJS application with docker.
This example focuses on Static Site Generator(SSG), and its purpose is to give an overview of the progress for the following [blog post](https://www.eden-reich.com/engineering-blog/build-a-blazing-fast-blog-using-nextjs-ssg/).

To work on development:

```sh
docker build -t blog --target development .
docker run --rm -it --name blog -p 3000:3000 -v ${PWD}:/app -w /app blog
yarn
yarn dev
```

To package it for production:

```sh
docker build -t blog --target production .
```
