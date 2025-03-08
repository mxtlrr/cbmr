# Client1 -- A client that opens a connection to the server
# under the name "mingus04", and prompts the player to start
# a match (this is when you'll run client2)
#
# Once the match starts, this client will prompt the user (again)
# if it would like to forefeit
import http.client

class Client:
  def __init__(self):
    pass
  def sendReq(connect: http.client.HTTPConnection, method: str, string: str) -> str:
    connect.request(method,string)
    return connect.getresponse().read().decode("utf-8")
  

connection = http.client.HTTPConnection("127.0.0.1", "3000", timeout=10)
Client.sendReq(connection, "GET", "/connect?name=mingus04")

input("Type anything and press enter to start a match. Client 2 should be running.")
print(Client.sendReq(connection, "GET", "/start_match?player1=mingus04&category=any"))

# Now we should have the match started, ask user to win!
input("Press enter to win the match")

# send post req
headers = { 'Content-Type': 'application/json' }

json_data = {
  "match_result": "win",
  "match_time": "17:26.444",
  "sender": "mingus04"
}
connection.request("POST", "/match_info", __import__("json").dumps(json_data), headers)

# Wait for GET