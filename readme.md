# Terraform httpdownload provider
This provider can be used to download files in a controlled manner on the host that's executing Terraform. Terraform will perform checksum validation and replace the file if it doesn't match.

## Why?
During a recent project I had the use case to execute a local binary during a Terraform run. Since I wanted to have the entire deploy as easy as possible and idempotent, I wanted to manage the binary using Terraform.

## Build from source
```
go build -o terraform-provider-httpdownload
cp terraform-provider-httpdownload ${your_terraform_code_dir}/terraform.d/plugins/darwin_amd64/ # this is for osx.
```

## Example
```
resource "httpdownload" "openshift-client-mac" {
  remote_url    = "https://mirror.openshift.com/pub/openshift-v4/clients/ocp/4.1.22/openshift-client-mac-4.1.22.tar.gz"
  filename      = "/tmp/openshift-client-mac-4.1.22.tar.gz"
  checksum      = "455234ae43a0a1a361ad474d785f0e1fadaac53d120f5444b610255fbe4f7a02"
  checksum_type = "sha256"
}
```

## Disclaimer
This is the first time I wrote a Terraform provider or used Go. Use at your own risk! Improvements are welcome via pull-requests
