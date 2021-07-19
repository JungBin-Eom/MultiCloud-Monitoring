import urllib
import hmac
import hashlib
import base64

def make_request(requests, secretKey):
  request = zip(requests.keys(), requests.values())
  request.sort(key=lambda x: str.lower(x[0]))

  requestUrl = "&".join(["=".join([r[0], urllib.quote_plus(str(r[1]))]) for r in request])
  hashStr = "&".join(["=".join([str.lower(r[0]), str.lower(urllib.quote_plus(str(r[1]))).replace("+","%20")]) for r in request])
  sig = urllib.quote_plus(base64.encodestring(hmac.new(secretKey, hashStr, hashlib.sha1).digest()).strip())
  print("Signature:", sig)
  requestUrl += "&signature="
  requestUrl += sig
  print(requestUrl)

if __name__ == '__main__':
  requests = {
    "apiKey": "0OcHRmqlLKxseRjIRoqW2sBtpIbaDDvnUElpbZVedZIVoZ1F11fcKi1n1MDGNuDWDXxBnG6Ba-jMFqSpAi5Tfg",
    "response" : "json",
    "command" : "listHostsMetrics"
  }
  secretKey = "xtbZVaUeYuds-ke_lCyRh0pZSdKdzUNHufwJeSvynO6847jJpWEb_aODEvsuHZ10os--xVFRAl3jepBiA33BAA"
  make_request(requests, secretKey)