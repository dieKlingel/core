import os, json, requests, sys

url = sys.argv[1]

response = requests.request("GET", url)
token = response.headers["X-FHEM-csrfToken"]
requests.request("GET", url + "&XHR=1&fwcsrf=" + token)
