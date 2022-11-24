## Container

```toml
kind = "container"
name = "busybox"
image = "docker.io/library/busybox:1.35.0"
command = "sleep 1m"
```

----

## Pod

```toml
kind = "pod"
name = "busybox"

[[containers]]
name = "busybox"
image = "docker.io/library/busybox:1.35.0"
command = "sleep 20"

[[containers]]
name = "busybox-2"
image = "docker.io/library/busybox:1.35.0"
command = "sleep 10"
```
