# Load Testing Guide

## Table of Contents

1. [Update Hosts](#1-update-hosts)
2. [Install k6](#2-install-k6)
3. [Modify Settings](#3-modify-settings)
4. [Run the Test](#4-run-the-test)
5. [Check the Results](#5-check-the-results)
6. [Adjust Rate Limits](#6-adjust-rate-limits)
7. [Debugging](#7-debugging)
8. [Summary](#8-summary)

## 1. Update Hosts

Update the `/etc/hosts` file with the following:

```bash
echo -e "\n142.132.160.156 load-testing-index.murmurations.network" | sudo tee -a /etc/hosts
echo -e "142.132.160.156 load-testing-library.murmurations.network" | sudo tee -a /etc/hosts
echo -e "142.132.160.156 load-testing-data-proxy.murmurations.network" | sudo tee -a /etc/hosts
```

## 2. Install k6

Install k6 using Homebrew:

```bash
brew install k6
```

## 3. Modify Settings

### For Servers

You have three options to modify:

1. **Replicas**: Adjust the number of replicas.
2. **CPU Limits**: Set the maximum CPU usage for the server.
3. **Memory Limits**: Set the maximum memory usage for the server.

To make these changes, update the following files:

- `charts/murmurations/charts/index/templates/index/dpl.yaml`
- `charts/murmurations/charts/validation/templates/validation/dpl.yaml`

### For Scripts

You can modify **the number of requests per second**.

To do this, edit these files:

- `test/load-testing/read_load_test.js`
- `test/load-testing/write_load_test.js`

## 4. Run the Test

Run the load test scripts:

```bash
k6 run test/load-testing/read_load_test.js
k6 run test/load-testing/write_load_test.js
```

## 5. Check the Results

Review the test results. Key metrics to focus on:

- **http_req_duration**: Measures response times. Target: < 500 milliseconds.
- **http_req_failed**: Measures request reliability. Target: 0% failures.
- **Process Rate for Adding New Nodes**: Check `http://load-testing-index.murmurations.network/v2/nodes` to ensure `number_of_results` matches the total requests. If not, go back to step 3 to adjust settings.

## 6. Adjust Rate Limits

If you hit the rate limit, update the rate settings in `charts/murmurations/charts/index/templates/index/config.yaml`:

```yaml
GET_RATE_LIMIT_PERIOD: "<some-big-number>-M"
POST_RATE_LIMIT_PERIOD: "<some-big-number>-M"
```

**Note:** Remember to manually deploy the index server after making these changes.

## 7. Debugging

### Access Kibana

```bash
kubectl port-forward service/index-kibana 5601:5601
```

Access Kibana at `http://localhost:5601`

### Access MongoDB

```bash
kubectl port-forward service/index-mongo 27017:27017
```

### Access Validation Redis

```bash
kubectl port-forward service/validation-redis 6379:6379
```

### Access NATS

```bash
kubectl port-forward svc/nats 4222:4222 -n murm-queue
nats stream ls # List all streams and see the number of messages remaining in each
```

## 8. Summary

1. **Validation Server Capacity:** A validation server with 64Mi memory and 128m CPU can process up to 4 requests per second. Increasing the server's CPU or memory, or adding more index servers, won’t significantly improve this rate.

2. **Balancing Index and Validation Servers:** The number of index servers should match the number of validation services to prevent performance issues. If the index servers are fewer, it may lead to a drop in the write queue.

3. **Index Server Performance:** An index server with 256Mi memory can handle requests at the following rates, depending on its CPU allocation:
   - 10 requests per second with 256m CPU.
   - 20 requests per second with 512m CPU.
   - 30 requests per second with 1024m CPU.

   The read performance is mainly limited by the total CPU capacity of the server.

### Example Configurations

- **20 Write Requests Per Second:** Deploy 5 index servers and 5 validation services.
- **60 Read Requests Per Second:** You have two configuration options:
  1. Use 6 index servers, each with 256Mi memory and 256m CPU.
  2. Use 2 index servers, each with 1024Mi memory and 1024m CPU.
