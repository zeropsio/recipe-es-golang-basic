# ZEROPS RECIPES

The concept of pre-prepared skeletons demonstrates the way how to set up and use technologies Zerops is supporting.

## Recipe: Basic integration of Elasticsearch & Golang

Pre-prepared skeleton demonstrating the way how to set up and use Golang and [Elasticsearch](https://www.elastic.co/elasticsearch) in Zerops. The main focus is on the client [Go SDK](https://www.elastic.co/guide/en/elasticsearch/client/go-api/current/index.html) (here, **v7.15** is used) configuration to connect the Elasticsearch service and show a simple example of how to insert a new document.

## Zerops import syntax

```yaml
services:
# Service will be accessible through zcli VPN under: http://recipees:9200
- hostname: recipees
  # Type and version of a used service.
  type: elasticsearch@7
  # Whether the service will be run on one or multiple containers.
  # Since this is a utility service, using only one container is fine.
  mode: NON_HA
# Service will be accessible through zcli VPN under: http://repiceesnodejs:3000
- hostname: recipeesgolang
  # Type and version of a used service.
  type: go@1
  # Whether the service will be run on one or multiple containers.
  # Since this is a utility service, using only one container is fine.
  mode: NON_HA
  # Service port definition (in Golang it's a user configurable parameter).
  ports:
  - port: 8080
    # If enabled, it means that a web server runs on the port (HTTP application protocol is supported).
    # It also means that you can enable a Zerops subdomain and access the service from the Internet.
    # You can even map public Internet domains with the option of automatic support for SSL certificates.
    httpSupport: true
  # A command that should start your service.
  startCommand: ./bin/server
  # Repository that contains Golang code with build and deploy instructions.
  buildFromGit: https://github.com/zeropsio/recipe-es-golang-basic@main
```

Copy & paste the import snippet above into the dialog of **Import service** functionality.

![Import](./images/Zerops-Import-Services-Dialog.png "Import Service Dialog")

See the [Zerops documentation](https://docs.zerops.io/documentation/export-import/project-service-export-import.html) to understand how to use it.
