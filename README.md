# Traefik redirector plugin

## Configuration

- **clientType** - type of client (`mock | http`)
- **cacheControlMaxAge** - max-age and s-maxage of cache-control header send with redirect response
- **maxAge** - TTL for cache (`seconds`)
- **endpoint** - API with redirect repository 
- **method** - HTTP Method to connect with API
- **data** - Payload for HTTP Request (`POST`)
- **debugParameter** - Name of query parameters use to debug

## Dev

### How to run

- set up the configuration
- build package with `docker-compose build`
- start with `docker-compose build`
- visit `http://whoami.localhost/`
- watch docker console for output

### Debug mode

Endpoint `http://whoami.localhost/?debug=redirects-dump` will show current state of cached redirects
Endpoint `http://whoami.localhost/?debug=redirects-load` will force cache warmup