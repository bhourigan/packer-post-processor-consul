This packer post processing plugin will insert AMI IDs and other image data into a consul
key for later retrieval by terraform or other software that consumes consul kv data.

It's still alpha, and needs to be polished and support added to nubis-builder. I'm using it
by dropping consul.json in projects/builder in my project.

```
{
"post-processors": [
  {
    "type": "consul",
    "consul_address": "<consul hostname>:8500",
    "consul_token": "<secret>",
    "aws_access_key": "{{user `aws_access_key`}}",
    "aws_secret_key": "{{user `aws_secret_key`}}",
    "project_name": "{{user `project_name`}}",
    "project_version": "{{user `project_version`}}"
  }
]
}
```
