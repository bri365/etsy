# etsy - and etcd playground

## Dependencies

### Environment variables 
The following terraform variables are expected to be populated
```bash
TF_VAR_guest_password
TF_VAR_vsphere_password
```

### Binaries
The terraform scripts rely on a local binary

```bash
export RELEASE="3.4.13"
wget https://github.com/etcd-io/etcd/releases/download/v${RELEASE}/etcd-v${RELEASE}-linux-amd64.tar.gz
```

