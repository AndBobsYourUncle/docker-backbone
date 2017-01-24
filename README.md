# Docker Backbone
- Facilitates the deployment of Docker images straight from repos

### Usage
```
git clone https://github.com/AndBobsYourUncle/docker-backbone
cd docker-backbone
```

You will then need to edit the docker-compose.yml and replace "deployer.domainname.com" with the subdomain and domain you want the deployer to listen on.


```
docker network create traffic_front
docker network create traffic_back
docker-compose up -d
docker exec -it docker-backbone bash
cat config/default.yml
```

Write down that token in a safe place. It is what will be used to remotely deploy projects.

### Deploying

To deploy a project to your server, here is a demo script:

```
#!/bin/bash

DEPLOY_FILE=`base64 docker-compose.yml`

cat <<EOF > deploy.json
{
  "project": "$APP_NAME",
  "compose_file": "$DEPLOY_FILE",
  "registry": {
    "url": "$REGISTRY_URL",
    "login": "$REGISTRY_USER",
    "password": "$REGISTRY_PASS"
  },
  "extra": {
    "TAG": "$CIRCLE_TAG"
  }
}
EOF

curl -H "Auth-Token: $DEPLOY_TOKEN" -X POST -d @deploy.json $DEPLOYER_URL
```