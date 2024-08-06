# How to Perform Load Testing

## Table of Contents
1. [Update the Hosts](#1-update-the-hosts)
2. [Install k6](#2-install-k6)
3. [Write a Script](#3-write-a-script)
4. [Run the Test](#4-run-the-test)
5. [Check the Results](#5-check-the-results)
6. [Note](#6-note)
7. [Debug](#7-debug)

## 1. Update the Hosts

```bash
echo -e "\n142.132.160.156 load-testing-index.murmurations.network\n142.132.160.156 load-testing-library.murmurations.network\n142.132.160.156 load-testing-data-proxy.murmurations.network" | sudo tee -a /etc/hosts
```

## 2. Install k6

```bash
brew install k6
```

## 3. Write a Script

- Use the constant-arrival-rate executor for load testing. Adjust the rate to increase requests per second.
- Virtual users (VUs) are similar to the number of threads executing the function. If the rate exceeds maxVUs, the desired rate won't be achieved.

```javascript
import http from 'k6/http';
import { check } from 'k6';

export let options = {
  scenarios: {
    constant_request_rate: {
      executor: 'constant-arrival-rate',
      rate: 100, // Number of requests per second
      timeUnit: '1s',
      duration: '1m', // Test duration
      preAllocatedVUs: 500, // Initial pool of virtual users
      maxVUs: 5000, // Maximum number of virtual users
    },
  },
};

export default function () {
  let res = http.get('http://load-testing-index.murmurations.network/v2/nodes?lat=51.493518&lon=0.009199&range=10km');
  
  // Check the response status
  let success = check(res, {
    'status is 200': (r) => r.status === 200,
  });

  // If the status is not 200, log an error message
  if (!success) {
    console.error(`Error: Expected status 200 but got ${res.status} - ${res.body}`);
  }
}
```

## 4. Run the Test

```bash
k6 run load-test.js
```

## 5. Check the Results

![image](https://github.com/user-attachments/assets/30cca494-c2f8-486f-b686-544da231b4e3)

### Key Metrics

- **http_req_duration**: Provides a comprehensive view of response times. Aim for < 500 milliseconds.
- **http_req_failed**: Ensures request reliability. Aim for 0%.

## 6. Note

If you hit the rate limit, please update the rate settings.

Update the rate limit in the config file `charts/murmurations/charts/index/templates/index/config.yaml`:

```yaml
GET_RATE_LIMIT_PERIOD: "<some-big-number>-M"
POST_RATE_LIMIT_PERIOD: "<some-big-number>-M"
```

Remember to manually deploy the index server after updating.

## 7. Debug

### Access Kibana

```sh
kubectl port-forward service/index-kibana 5601:5601
```

### Access MongoDB

```sh
kubectl port-forward service/index-mongo 27017:27017
```

### Access NATS

```sh
kubectl port-forward svc/nats 4222:4222 -n murm-queue
nats stream ls
nats stream info nodes
```
