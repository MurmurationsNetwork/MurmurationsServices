<div align="center">
<br/>
<h1>Murmurations Services</h1>
</div>

## Run locally

1. Make sure you have installed [ingress controller](https://kubernetes.github.io/ingress-nginx/deploy/)

2. Start the services

```
make dev
```

## Directory Layout

* **services** contains packages that compile to applications that are long-running processes (such as API servers).
* **common** contains packages that shared across different services.
