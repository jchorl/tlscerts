require 'openssl'
require 'webrick'
require 'webrick/https'

cert = OpenSSL::X509::Certificate.new File.read './server.pem'
pkey = OpenSSL::PKey::EC.new File.read './server-key.pem'

server = WEBrick::HTTPServer.new(:Port => 443,
                                 :SSLEnable => true,
                                 :SSLCertificate => cert,
                                 :SSLPrivateKey => pkey,
                                 :SSLCACertificateFile => './ca.pem')
server.mount_proc '/' do |req, res|
  res.body = 'Hi!'
end
server.start
