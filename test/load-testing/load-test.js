import http from 'k6/http';
import { check } from 'k6';

// Function to generate a random number between 0 and 2^64 - 1
function getRandomNumber() {
  return Math.floor(Math.random() * Math.pow(2, 64));
}

// Define the profile URL as a constant
const BASE_PROFILE_URL = 'http://5.78.90.240/profile';

// Define the latitude and longitude constants
const LATITUDE = 48.8566;
const LONGITUDE = 2.3522;

export let options = {
  scenarios: {
    constant_request_rate: {
      executor: 'constant-arrival-rate',
      rate: 1, // Number of requests per second
      timeUnit: '1s',
      duration: '1m', // Test duration
      preAllocatedVUs: 50, // Initial pool of virtual users
      maxVUs: 100, // Maximum number of virtual users
    },
  },
};

function sendPostRequest() {
  const profileUrlWithRandom = `${BASE_PROFILE_URL}/${getRandomNumber()}`;
  const url = 'http://load-testing-index.murmurations.network/v2/nodes';
  const payload = JSON.stringify({
    profile_url: profileUrlWithRandom,
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  let res = http.post(url, payload, params);

  let success = check(res, {
    'POST request status is 200': (r) => r.status === 200,
  });

  if (!success) {
    console.error(`POST Error: Expected status 200 but got ${res.status} - ${res.body}`);
  }

  return res;
}

function sendGetRequest() {
  const url = `http://load-testing-index.murmurations.network/v2/nodes?lat=${LATITUDE}&lon=${LONGITUDE}&range=10km`;
  let res = http.get(url);

  let success = check(res, {
    'GET request status is 200': (r) => r.status === 200,
  });

  console.log(res.body)

  if (!success) {
    console.error(`GET Error: Expected status 200 but got ${res.status} - ${res.body}`);
  }

  return res;
}

export default function () {
  sendPostRequest();
  sendGetRequest();
}
