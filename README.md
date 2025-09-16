# word-hero

```shell
make dockertag
docker rm -f word-hero
docker run --name word-hero -p 8080:8080 --restart always -d word-hero:latest
```