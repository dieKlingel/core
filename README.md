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
  - trigger: "unlock"
    environment: bash
    script: |
      echo "Hallo Welt!"
  - trigger: "ring"
    environment: python
    script: |
      print('Hallo Welt!')

media:
  video:
    src: autovideosrc ! video/x-raw, width=1280, height=720, framerate=30/1 ! videoconvert ! x264enc tune=zerolatency bitrate=500 speed-preset=superfast ! appsink name=h264sink
  audio:
    src: autoaudiosrc ! audioconvert ! opusenc ! appsink name=opussink
    sink: appsrc format=time do-timestamp=true name=opussrc ! application/x-rtp, payload=127, encoding-name=OPUS ! rtpopusdepay ! decodebin ! autoaudiosink

mqtt:
  uri: mqtt://server.dieklingel.com:1883
  username: ''
  password: ''

rtc:
  ice-servers:
    - urls: stun:stun1.l.google.com:19302
    - urls: stun:stun1.l.google.com:19302
      username: ''
      credentials: ''
"
```

Run the core

```bash
sudo snap set dieklingel-core daemon=true
```

## Roadmap

### streamer boy -- 0.3.1 (release: 2023-09-15)

- [x] mjpeg api for video

### new born baby -- 0.3.2 (release: ?)

- [x] establish a call over mqtt
- [x] execute actions over mqtt
- [x] use bash or python as action environment
- [x] add microphone support

### app combat -- 0.4.0 (release: ?)

- [ ] add db support for devices/apps
- [ ] allow to push mqtt topics to inactive devices/apps
- [ ] store devices by last-will-topic (inactive message)

### environment explorer -- 0.5.0 (release: ?)

- [ ] emit events on well-known actions
- [ ] add gpio (Raspberry Pi) support

### camera combat -- 0.6.0 (release: ?)

- [ ] add rtsp support
- [ ] save camera capture on movement

### call a friend -- 0.7.0 (release: ?)

- [ ] add sip support

### Future Versions

- [ ] add libcamera support
