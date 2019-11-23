const fs = require('fs');
const https = require('https');

const port = 443;

const server = https.createServer(
  {
    cert: fs.readFileSync('./server.pem'),
    key: fs.readFileSync('./server-key.pem'),
    ca: fs.readFileSync('./ca.pem'),
  },
  (req, res) => {
    res.statusCode = 200;
    res.setHeader('Content-Type', 'text/plain');
    res.end('Hi!\n');
  }
);

server.listen(port);
