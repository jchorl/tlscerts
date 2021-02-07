import ssl
from flask import Flask

app = Flask(__name__)

port = 443

@app.route('/')
def index():
    return 'Hi!'

ctx = ssl.create_default_context(purpose=ssl.Purpose.CLIENT_AUTH)
ctx.load_verify_locations(cafile='./ca.pem')
ctx.load_cert_chain(
    certfile='./server.pem',
    keyfile='./server-key.pem',
)
ctx.verify_mode = ssl.CERT_REQUIRED

app.run(debug=True, ssl_context=ctx, host="0.0.0.0", port=port)
