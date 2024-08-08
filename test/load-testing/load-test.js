import http from 'k6/http';
import { check } from 'k6';

// Function to generate a UUID
function generateUUID() {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
    var r = (Math.random() * 16) | 0,
      v = c == 'x' ? r : (r & 0x3) | 0x8;
    return v.toString(16);
  });
}

// Define the profile URL as a constant
const POST_URL = 'http://load-testing-index.murmurations.network';
const BASE_PROFILE_URL = 'http://5.78.90.240/profile';

// Define the latitude and longitude constants
const LATITUDE = 48.8566;
const LONGITUDE = 2.3522;

export let options = {
  scenarios: {
    constant_request_rate: {
      executor: 'constant-arrival-rate',
      rate: 20, // Number of requests per second
      timeUnit: '1s',
      duration: '1m', // Test duration
      preAllocatedVUs: 1000, // Initial pool of virtual users
      maxVUs: 5000, // Maximum number of virtual users
    },
  },
};

function sendPostRequest() {
  const profileUrlWithUUID = `${BASE_PROFILE_URL}/${generateUUID()}`;
  const url = `${POST_URL}/v2/nodes`;
  const payload = JSON.stringify({
    profile_url: profileUrlWithUUID,
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

  console.log(res.body);

  if (!success) {
    console.error(`GET Error: Expected status 200 but got ${res.status} - ${res.body}`);
  }

  return res;
}

export default function () {
  sendPostRequest();
  sendGetRequest();
}
