# peer-finder

kubernetes peer-finder


# usage



Dockerfile
```
...

# script
ADD bootstrap.sh /bootstrap.sh
ADD on-change.sh /on-change.sh
RUN chmod +x /bootstrap.sh
RUN chmod +x /on-change.sh

ENTRYPOINT ["/bootstrap.sh"]
```

bootstrap.sh
```
/peer-finder -on-change=/on-change.sh -on-start=/on-change.sh -service=sfs -ns=default -dns-suffix=svc.cluster.local

while true; do sleep 1000; done
```

on-change.sh
```
#! /bin/bash

> /etc/peer_config
while read -ra LINE; do
    IP=${LINE#*,}
    DNS=${LINE%%,*}
    HOST=${LINE%%.*}

    PEERS=("${PEERS[@]}" ${DNS})

    echo "${DNS}" >> /etc/peer_config
done

echo ${PEERS}
```

`kubectl logs sfs-1`

```
2017/04/17 14:03:14 Peer list updated
iam sfs-1.sfs.default.svc.cluster.local,172.1.56.2
was []
now [sfs-0.sfs.default.svc.cluster.local,172.1.65.7 sfs-1.sfs.default.svc.cluster.local,172.1.56.2]
2017/04/17 14:03:14 execing: /on-change.sh with stdin: sfs-0.sfs.default.svc.cluster.local,172.1.65.7
sfs-1.sfs.default.svc.cluster.local,172.1.56.2
2017/04/17 14:03:14 sfs-0.sfs.default.svc.cluster.local
[dev_dean@VM_61_2_centos statefulset]$ kubectl logs sfs-0
2017/04/17 14:03:08 lookup sfs on 182.1.0.100:53: server misbehaving
2017/04/17 14:03:09 Peer list updated
iam sfs-0.sfs.default.svc.cluster.local,172.1.65.7
was []
now [sfs-0.sfs.default.svc.cluster.local,172.1.65.7]
2017/04/17 14:03:09 execing: /on-change.sh with stdin: sfs-0.sfs.default.svc.cluster.local,172.1.65.7
2017/04/17 14:03:09 sfs-0.sfs.default.svc.cluster.local
2017/04/17 14:03:15 Peer list updated
iam sfs-0.sfs.default.svc.cluster.local,172.1.65.7
was [sfs-0.sfs.default.svc.cluster.local,172.1.65.7]
now [sfs-0.sfs.default.svc.cluster.local,172.1.65.7 sfs-1.sfs.default.svc.cluster.local,172.1.56.2]
2017/04/17 14:03:15 execing: /on-change.sh with stdin: sfs-0.sfs.default.svc.cluster.local,172.1.65.7
sfs-1.sfs.default.svc.cluster.local,172.1.56.2
2017/04/17 14:03:15 sfs-0.sfs.default.svc.cluster.local
```

`kubectl exec sfs-1 -- cat /etc/peer_config`

```
sfs-0.sfs.default.svc.cluster.local
sfs-1.sfs.default.svc.cluster.local
```
