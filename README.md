# Example AdServer

This project is just a showcase of my coding style. Feel free to do whatever you want with it!

## URLs

### Serve ad:

|name|method|url|example body|result|
|----|------|---|------------|------|
|Create Ad|POST|http://localhost:8000/ads/|`{"image_url":"http://example.org/img.png","clickthrough_url":"http://example.org/landing"}`|*Status*: 201<br>*Location*:Ad server url|
|Update Ad|POST|http://localhost:8000/ads/{id}|`{"image_url":"http://example.org/img.png","clickthrough_url":"http://example.org/landing"}`|*Status*: 204|
|List Active Ads|Get|http://localhost:8000/ads/active|-|*Body*: `[{"ID":1,"ImageURL":"http://example.org/img.png","ClickThroughURL":"http://example.org/landing"}]`|
|Serve Ad|Get|http://localhost:8001/{id}|-|*Body*: `<a href="http://example.org/landing"><img src="http://example.org/img.png"></a>`|



## Testing

You can run the e2e suit with make:

```sh
make test-e2e
```

Unit tests:

```sh
make test-unit
```

## Building and running

You can run the project:

```sh
make run
```

This will build the required images. You can alternatively just build the images:

```sh
make build
```
