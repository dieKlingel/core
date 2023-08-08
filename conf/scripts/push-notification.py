import os, json, requests
from datetime import datetime

SIGN = os.getenv("SIGN")

result = requests.request("GET", "http://localhost:8081/proxy/dieklingel/mayer/kai/devices")
if result.status_code != 200:
	exit()

devices = result.json()
devices = filter(lambda device: SIGN in device['signs'], devices)

tokens = list(map(lambda device: device['token'], devices))
print("send push notification to:", tokens)

if not tokens:
    print("no tokens: early exit!")
    exit()

now = datetime.now()

payload = {
  'tokens': tokens,
  'id': 'random123',
  'title': 'Jemand steht vor deiner Tür!',
  'body': now.strftime("%d.%m.%Y %H:%M:%S")
}

requests.post("https://fcm-worker.dieklingel.workers.dev/fcm/send", json = payload)
