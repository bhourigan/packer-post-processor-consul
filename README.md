This packer post processing plugin will insert AMI IDs and other image data into a consul
key for later retrieval by terraform or other software that consumes consul.

It's still alpha (pre alpha?), and needs some polishing. The plan is to integrate support
for this directly into nubis-builder so no end user configuration other than the address
and token will be required.

Here's how I am using it during development.

$ go get github.com/bhourigan/post-processor-consul
$ cd ~/gocode/src/github.com/bhourigan/post-processor-consul && make

Copy the resulting 'post-processor-consul' binary to your packer binary directory (such as
~/gocode/bin if you're building packer from source) or into ~/.packer.d/plugins as
'packer-post-processor-consul'

If you're using nubis-builder all you need to do is add consul_address into your secrets,
and it will automatically load packer/post-processors/consul.json during run time (take a
peek for the config parameters that get passed)

If you're using it outside of nubis-builder, you'll need to add this to your packer json
file:

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
