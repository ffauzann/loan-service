# policy-read-only.hcl
path "my-service/*" {
    capabilities = ["read", "list"]
}