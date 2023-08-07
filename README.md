# dieKlingel Core

## Install

1. install go

   <https://go.dev/doc/install>

2. dependencies:

   ```bash
   sudo apt-get install libgtk-3-dev
   ```

   ```bash
   sudo apt-get install libgstreamer1.0-dev libgstreamer-plugins-base1.0-dev libgstreamer-plugins-bad1.0-dev gstreamer1.0-plugins-base gstreamer1.0-plugins-good gstreamer1.0-plugins-bad gstreamer1.0-plugins-ugly gstreamer1.0-libav gstreamer1.0-tools gstreamer1.0-x gstreamer1.0-alsa gstreamer1.0-gl gstreamer1.0-gtk3 gstreamer1.0-qt5 gstreamer1.0-pulseaudio
   ```

## Config

### Media

Sample media config:

```yaml
media:
  video-src: autovideosrc ! x264enc tune=zerolatency bitrate=500 speed-preset=superfast
  audio-src: autoaudiosrc ! audioconvert ! opusenc
```

## TODO

- add libcamera support
- gstreamer microphone
- gstreamer audio output
