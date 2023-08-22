# Traefik redirector plugin

## Configuration

- **clientType** = type of client (`mock | http`)
- **maxAge** = TTL for cache (`seconds`)
- **endpoint** = API with redirect repository 
- **method** = HTTP Method to connect with API
- **data** = Payload for HTTP Request (`POST`)

## Dev

### How to run

- set up the configuration
- build package with `docker-compose build`
- start with `docker-compose build`
- visit `http://whoami.localhost/`
- watch docker console for output

### Debug mode

Endpoint `http://whoami.localhost/?unicorn=redirects` will show current state of cached redirects