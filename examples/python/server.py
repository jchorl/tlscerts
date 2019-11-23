import ssl
from flask import Flask

app = Flask(__name__)

@app.route('/')
def index():
    return 'Hi!'

ctx = ssl.create_default_context(capath='./ca.pem')
ctx.load_cert_chain(
    certfile='./server.pem',
    keyfile='./server-key.pem',
)
app.run(debug=False, ssl_context=ctx)
