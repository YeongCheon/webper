## Deploy

``` bash
BUCKET_NAME=YOUR_OWN_BUCKET_NAME
DOCKER_TAG=YOUR_OWN_DOCKER_TAG

go build -o webper cmd/main.go
docker build -t $DOCKER_TAG .
docker push $DOCKER_TAG
gcloud run deploy webper --image $DOCKER_TAG --region asia-northeast3 --update-env-vars BUCKET_NAME=$BUCKET_NAME
```
