resource "http-download" "readme_from_website" {
  remote_url    = "https://mirror.openshift.com/pub/openshift-v4/clients/ocp/4.1.22/openshift-client-mac-4.1.22.tar.gz"
  filename      = "/tmp/openshift-client-mac-4.1.22.tar.gz"
  checksum      = "455234ae43a0a1a361ad474d785f0e1fadaac53d120f5444b610255fbe4f7a02"
  checksum_type = "sha256"
}
