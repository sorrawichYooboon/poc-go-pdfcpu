Go build

``` bash
env GOOS=linux GOARCH=amd64 go build -o app .
```

Docker build
``` bash
docker build -t app
```

Docker run
``` bash
docker run --cpus=2 -p 8080:8080 app
```

Example Request
``` bash
curl -X POST http://localhost:8080/generate-pdf --output result.pdf
```

