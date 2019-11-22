require 'faraday'
require 'openssl'

cert = OpenSSL::X509::Certificate.new File.read './client.pem'
pkey = OpenSSL::PKey::EC.new File.read './client-key.pem'
connection = Faraday::Connection.new 'https://server/', :ssl => {
    :client_cert => cert,
    :client_key => pkey,
    :ca_file => './ca.pem'
}
resp = connection.get '/'
if resp.status != 200
    exit 1
end
exit 0
