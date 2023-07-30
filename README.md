A plugin to kuboard cicd.

See [kuboard cicd](https://kuboard.cn/guide/cicd/) and [kuboard swagger](https://demo.kuboard.cn/swagger/index.html).

# Usage

The following settings changes this plugin's behavior.

* cluster (required) -- cluster name.
* kind (required) -- workload kind.
* name (required) -- workload name.
* namespace (required) -- workload namespace.
* image (required) -- the image to be changed.
* tag (required) -- to set the version of the image.
* kuboard_uri (required) -- kuboard uri.
* kuboard_username (required) -- kuboard username.
* kuboard_key (required) -- kuboard access key.

Below is an example `.drone.yml` that uses this plugin.

```yaml
kind: pipeline
name: default

steps:
- name: run kuboard plugin
  image: suyar/drone-kuboard
  pull: if-not-exists
  settings:
    cluster: dev
    kind: deployments
    name: api
    namespace: api-dev
    image: suyar/php
    tag: 8.2
    kuboard_uri: https://demo.kuboard.cn/
    kuboard_username: admin
    kuboard_key: exxrkyj8wkbkczb2n6fbcehxenkefza2
```
