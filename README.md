This packer post processing plugin will insert AMI IDs and other image data into a consul
key for later retrieval by terraform or other software that consumes consul.

It's still alpha (pre alpha?), and needs some polishing. The plan is to integrate support
for this directly into nubis-builder so no end user configuration other than the address
and token will be required.

Here's how I am using it during development.

Build this module and drop the binary into ~/.packer.d/plugins, add the following post
processor into $project_path/nubis/builder

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
