#!/usr/bin/env bash

# from https://medium.com/nirman-tech-blog/setting-up-etcd-cluster-with-tls-authentication-enabled-49c44e4151bb

# initialize certificate authority
cfssl gencert -initca ca-csr.json | cfssljson -bare ca -

# generate cluster certificate for client access
cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=peer cluster-csr.json | cfssljson -bare cluster

# generate client certificate
cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=client client-csr.json | cfssljson -bare client

# zip for easier transfer to servers
tar czf etcd.tar.gz ca.pem cluster*pem client*pem
