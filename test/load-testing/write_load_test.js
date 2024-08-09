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

// Define the profile URLs as constants
const POST_URL = 'http://load-testing-index.murmurations.network';
const SMALL_PROFILE_URL = 'http://5.78.90.240/profile/small';
const MEDIUM_PROFILE_URL = 'http://5.78.90.240/profile/medium';
const LARGE_PROFILE_URL = 'http://5.78.90.240/profile/large';

// Define the options for the test
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

function selectProfileUrl() {
  const rand = Math.random();

  if (rand < 0.4) {
    return `${SMALL_PROFILE_URL}/${generateUUID()}`;
  } else if (rand < 0.7) {
    return `${MEDIUM_PROFILE_URL}/${generateUUID()}`;
  } else {
    return `${LARGE_PROFILE_URL}/${generateUUID()}`;
  }
}

function sendPostRequest() {
  const profileUrlWithUUID = selectProfileUrl();
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

export default function () {
  sendPostRequest();
}
