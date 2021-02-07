import requests


resp = requests.get(
    'https://server/',
    cert=('./client.pem', './client-key.pem'),
    verify='./ca.pem',
)
if resp.status_code != 200:
    print('expected status_code 200 but got {}'.format(resp.status_code))
    exit(1)
