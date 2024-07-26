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
