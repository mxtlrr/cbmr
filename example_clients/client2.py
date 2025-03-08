import http.client

class Client:
  def __init__(self):
    pass
  def sendReq(connect: http.client.HTTPConnection, method: str, string: str) -> str:
    connect.request(method,string)
    return connect.getresponse().read().decode("utf-8")
  

connection = http.client.HTTPConnection("127.0.0.1", "3000", timeout=10)
Client.sendReq(connection, "GET", "/connect?name=testaccount")
while True:
  pass