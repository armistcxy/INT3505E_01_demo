## Use Docker to quickly creating Swagger UI container and render Swagger UI

Run this command
```bash
docker run -p 8081:8080 -e SWAGGER_JSON=/docs/openapi.yaml -v $(pwd):/docs swaggerapi/swagger-ui
```

## Exposing to Internet

Use ngrok:

```bash
ngrok http 8081
```