import os, json, requests
from datetime import datetime

USERNAME = "<your username>"
PASSWORD = "<your password>"
PREFIX = "dieklingel/mayer/kai"
SIGN = os.getenv("SIGN")

result = requests.request(
  "GET",
  "http://localhost:8081/" + PREFIX + "/devices",
  headers = {
    "Username": USERNAME,
    "Password": PASSWORD
  }
)

if result.status_code != 200:
  print("early exit, erro while fetching devices:", result.status_code, result.text)
  exit(1)


print("received ring for ", SIGN)

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
  'id': '100',
  'title': 'Jemand steht vor deiner Tür!',
  'body': now.strftime("%d.%m.%Y %H:%M:%S")
}

requests.post("https://fcm-worker.dieklingel.workers.dev/fcm/send", json = payload)


result = requests.request(
  "GET",
  "http://localhost:8081/" + PREFIX + "/camera/snapshot",
  headers = {
    "Username": USERNAME,
    "Password": PASSWORD
  }
)

if result.status_code != 200:
  print("early exit, erro while fetching snapshot:", result.status_code, result.text)
  exit(0)


payload = {
  'tokens': tokens,
  'id': '100',
  'title': 'Jemand steht vor deiner Tür!',
  'body': now.strftime("%d.%m.%Y %H:%M:%S"),
  'image': result.text
}

requests.post("https://fcm-worker.dieklingel.workers.dev/fcm/send", json = payload)
