# simple-images-api

[![Build Status](https://travis-ci.com/norbjd/simple-images-api.svg?branch=master)](https://travis-ci.com/norbjd/simple-images-api)

A simple API to store and retrieve images, with Minio as a backend.

# Run locally

You can start API, Minio backend and an OpenAPI UI by running :

```
make run
```

It requires `docker-compose`.

By default, API run on `localhost:8080`, Minio on `localhost:9000`.

# Test locally

Once the stack started :

```
make openapi-ui-open
```

starts a browser pointing to the OpenAPI UI. It requires `chromium`. From there, you can test the API.

Or you can do `curl` requests to `localhost:8080`.