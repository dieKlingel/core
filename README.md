# dieKlingel Core

[![Get it from the Snap Store](https://snapcraft.io/static/images/badges/en/snap-store-black.svg)](https://snapcraft.io/dieklingel-core)

This repository contains the core application for the dieKlinge project. The core covers features like peer to peer connection through webrtc, storing devices and executing actions to notify you if someone rings your bell. The core is designed to be configurable so that it could be integrated into a smarthome system like fhem. The communication. The startup configuration is set throug a config file. After startup the communcation with core takes place over mqtt or in some cases http. Just install the core on a raspberry pi and run it.

## Getting Started

Install the latest build from snapcraft

```bash
sudo snap install dieklingel-core
```

Configure the core

```bash
sudo set install dieklingel-core core="
actions:
  - trigger: ring
    lane: python3 ./scripts/push-notification.py
  - trigger: unlock2
    lane: |
     touch hallo.txt
     echo unlock2 > hallo.txt 

media:
  video-src: autovideosrc ! video/x-raw, width=1280, height=720, framerate=30/1 ! videoconvert ! x264enc tune=zerolatency bitrate=500 speed-preset=superfast ! appsink name=h264sink
  audio-src: autoaudiosrc ! audioconvert ! opusenc ! appsink name=opussink
  audio-sink: appsrc format=time do-timestamp=true name=opussrc ! application/x-rtp, payload=127, encoding-name=OPUS ! rtpopusdepay ! decodebin ! autoaudiosink

mqtt:
  uri: mqtt://server.dieklingel.com:1883/dieklingel/mayer/kai/
  username: ''
  password: ''

rtc:
  ice-servers:
    - urls: stun:stun1.l.google.com:19302
    - urls: stun:stun1.l.google.com:19302
      username: adsa
      credentials: ada
"
```

Run the core

```bash
sudo snap set dieklingel-core daemon=true
```

## TODO

- add libcamera support
