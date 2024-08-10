import http from 'k6/http';
import { check } from 'k6';

// Define the profile URL as a constant
const BASE_URL = 'http://load-testing-index.murmurations.network/v2/nodes';

// Define the options for the test
export let options = {
  scenarios: {
    constant_request_rate: {
      executor: 'constant-arrival-rate',
      rate: 130, // Number of requests per second
      timeUnit: '1s',
      duration: '1m', // Test duration
      preAllocatedVUs: 1000, // Initial pool of virtual users
      maxVUs: 5000, // Maximum number of virtual users
    },
  },
};

function sendGetRequest() {
  // URLs to be requested with different probabilities
  const url50 = `${BASE_URL}`;
  const url25 = `${BASE_URL}?page_size=100`;
  const url15 = `${BASE_URL}?page_size=500`;
  const url10 = `${BASE_URL}?page_size=1000`;

  // Generate a random number between 0 and 100
  let random = Math.random() * 100;
  let url;

  // Select the URL based on the random number to match the specified percentages
  if (random < 50) {
    url = url50; // 50% of the time
  } else if (random < 75) {
    url = url25; // 25% of the time
  } else if (random < 90) {
    url = url15; // 15% of the time
  } else {
    url = url10; // 10% of the time
  }

  let res = http.get(url);

  let success = check(res, {
    'GET request status is 200': (r) => r.status === 200,
  });

  if (!success) {
    console.error(`GET Error: Expected status 200 but got ${res.status} - ${res.body}`);
  }

  return res;
}

export default function () {
  sendGetRequest();
}
