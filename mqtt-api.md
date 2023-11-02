# MQTT API

## Operations

### __PUBLISH__ to: `$CORE_DEVICE_ID/connections/offer`

Description: start a new call with the core

```javascript
{
    headers: {
        senderDeviceId: String,
        senderSessionId: String,
    },
    body: {
        iceCandidate: {
            candidate: String,
            sdpMid: Number,
            sdpMLineIndex: String,
        }
    }
}
```

---

### __PUBLISH__ to: `$CORE_DEVICE_ID/connections/answer`

Description: Answer a call, wich was started by the core

```javascript
{
    headers: {
        senderDeviceId: String,
        senderSessionId: String,
        sessionId: String
    },
    body: {
        type: Number,
        sdp: String
    }
}
```

---

### __PUBLISH__ to: `$CORE_DEVICE_ID/connections/candidate`

Description: Add a candidate to an ongoing call

```javascript
{
    headers: {
        senderDeviceId: String,
        senderSessionId: String,
        sessionId: String
    },
    body: {
        type: Number,
        sdp: String
    }
}
```

---

### __PUBLISH__ to: `$CORE_DEVICE_ID/connections/close`

Description: Close an ongoing call, with the core

```javascript
{
    headers: {
        senderDeviceId: String,
        senderSessionId: String,
        sessionId: String
    }
}
```

---

### __SUBSCRIBE__ to: `$SENDER_DEVICE_ID/connections/offer`

Description: start a new call with the client

```javascript
{
    headers: {
        senderDeviceId: String,
        senderSessionId: String,
    },
    body: {
        iceCandidate: {
            candidate: String,
            sdpMid: Number,
            sdpMLineIndex: String,
        }
    }
}
```

---

### __SUBSCRIBE__ to: `$SENDER_DEVICE_ID/connections/close`

Description: Close an ongoing call, with the client

```javascript
{
    headers: {
        senderDeviceId: String,
        senderSessionId: String,
        sessionId: String
    }
}
```

---

### __PUBLISH__ to: `$CORE_DEVICE_ID/apps/add`

Description: sets a device to active

```javascript
{
    headers: {
        senderDeviceId: String
    },
    body: {
        token: String
        timeToLive: Number
    }
}
```

---

### __PUBLISH__ to: `$CORE_DEVICE_ID/apps/remove`

Description: sets a device to inactive

```javascript
{
    headers: {
        senderDeviceId: String
    },
    body: {
        token: String,
    }
}
```

---

### __PUBLISH__ to: `$CORE_DEVICE_ID/actions/trigger`

Description: trigger an action on the system.

```javascript
{
    headers: {},
    body: {
        pattern: String,
        environment: {
            [key: String]: String
        }
    }
}
```
