# STEP2 Implement tool's lambda in your favorite language
In this step, you'll implement your tool's lambda.

SlackHub is designed to work regardless of the language of the latter program. **You can use any language you love!**

## JSON sent from SlackHub to your lambda
![Role of SlackHub](https://github.com/Jimon-s/slackhub/blob/images/role.png)

Once users submit a modal view, SlackHub authorize the message and pass the input params to your lambda.

You can see a sample of the JSON in [Slack's developer website](https://api.slack.com/reference/interaction-payloads/views#view_submission) (See view.submission event).

Remember that user's input data is in
```
view.state
```
field. 

### Convenient parameter
SlackHub adds some convenient params in the JSON.

In `view.private_metadata` field, you can see below.

```json
{
    "response_url": "https://xxxxxxxxxxx",
    "channel_id"  : "xxxxxxxxxxxxxxxxxxx"
}
```

`response_url`

This field includes the url tied to the slack channel. You can send a message to slack only sending POST request to this url.

```
POST https://xxxxxxxxxxx
Content-type: application/json
{
    "text": "YOUR_MESSAGE_HERE!"
}
```

You can send messages **within 30 minutes up to 5 times**. If you send a message by response_url, it replaces the slackhub's wait message. (@USER invokes XXXX. Please wait a moment.)

`channel_id`

This field includes the channel ID of the slack channel. If you learn Slack API and want to send a message by API, you can use this ID.

## Best Practice
To create more user friendly tools, I recommend you to **send a compression message**.

When you send, please use the `response_url` descrebed above.

## Implement your tool
That ’s all for the explanation. Let's implement your tool's lambda and deploy to your AWS account!

We provide you some example codes in this repository. Please see [the sample directly](https://github.com/Jimon-s/slackhub/blob/master/examples).
