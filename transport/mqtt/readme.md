### LoRA transport

This Transport pipes messages from a MQTT data source to a GCloud pubsub

If required the Transport also listens for Downlink messages and republishes downlink ACK, CONFIRMED etc

Config files have the form:
```
{
  "projectName": "**google project name**",
  "mqtt": {
    "appID": "**mqtt appID**",
    "username": "**mqtt appID**",
    "address": "**mqtt address e.g: ssl://eu1.cloud.thethings.network:8883**"
  },
  "topics": {
    "joins": "**gcloud joins topic**",
    "uplinks": "**gcloud uplinks topic**",
    "downlinks": "**gcloud downlinks topics (used to create sub)**",
    "downlinkReceipts": "**gcloud downlink receipts topic**"
  },
  "subscriptions": {
    "downlinks": "**gcloud downlinks subscription**"
  }
}
```

We use secrets to secure access to API keys etc

The secret is set up independently an accessed in code
```
  "secret": {
    "name": "ttn-milesight",
    "version": 1
  },
```