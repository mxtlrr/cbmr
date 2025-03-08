
from flask import Flask, request, jsonify
import json, urllib.parse
import requests

app = Flask(__name__)

@app.route('/match_info', methods=['GET'])
def receive_get():
	data_str = request.args.get("data")  # Extract the "data" parameter
	if not data_str:
		return jsonify({"error": "Missing or empty data parameter"}), 400

	data_str = urllib.parse.unquote(data_str)

	if data_str.startswith('"') and data_str.endswith('"'):
		data_str = data_str[1:-1]

	data = json.loads(data_str)
	print(f"Parsed JSON: {data}")

	json_data = {
		"winner": "mingus04",
		"time": "17:26.444",
		"type": "loss",
		"accepted": True
	}
	headers = {"Content-Type": "application/json"}
	requests.post(f"http://127.0.0.1:3000/end_match",
									data=json.dumps(json_data),
									headers=headers)
	return ""

if __name__ == "__main__":
    app.run(host='0.0.0.0', port=8080)