const fs = require('fs');
const https = require('https');

https.get(
  {
    hostname: 'server',
    port: 443,
    path: '/',
    method: 'GET',
    cert: fs.readFileSync('./client.pem'),
    key: fs.readFileSync('./client-key.pem'),
    ca: fs.readFileSync('./ca.pem')
  },
  res => {
    if (res.statusCode != 200) {
      console.error(`expected status 200 but found ${res.statusCode}`);
      return process.exit(1);
    }
  }
).on('error', error => {
  console.error(error);
  process.exit(1);
});
