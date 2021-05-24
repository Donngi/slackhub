# Next Step: Using advanced mode
This page describes how to use `Advanced mode`. If you don't know well about slack's lifecycle events, I don't recommend you to use this mode.

To utilize advanced mode, you are able to realize the chain of modal views. One of the example of the tool using advanced mode is SlackHub's official tool `Editor`. If you have any interests, please see the code.

## When you should use advanced mode?
If you want to make these application, you should consider to use this mode.

- App with continuous transition of modals
- App with modal input validation
- App with dynamic modals

## Difference between Normal and Advanced mode
The biggest difference between normal and advanced mode is how to invoke tool's lambda.

| Boot Mode | How to invoke |
| :--- | :--- |
| Normal | Asynchronously |
| Advanced | Synchronously |

You can see the difference graphically in the page below.

[SlackHub sequence](https://github.com/Jimon-s/slackhub/blob/master/documents/slackhub_architecture)

In advanced mode, SlackHub acts as just a proxy. You must implement all request/response sequence by themselves.

## How to implement
To implement your lambda, You should follow these two steps.

1. Change `Boot Mode` to advance.
2. Implement your lambda which supports synchronous invocation.
 
### STEP1 Change Boot Mode
First, you should change Boot Mode by using SlackHub's official tool `Editor`. Of course, you can initially set `Boot Mode = Advanced` by using `Register`.

### STEP2 Implement your lambda
#### Response
Your lambda must return the response to SlackHub like below. This format is `AWS API Gateway proxy response` format. You should use SDK of your language.

```
StatusCode:      200,
IsBase64Encoded: false,
Headers: {
    "Content-Type": "application/json",
},
Body: body
```

In this response, the `body` is slack's `response_action` format. You can see the datail in the page below.

[Slack - Using modals in Slack apps: Modifying Modals](https://api.slack.com/surfaces/modals/using#opening#modifying)

#### External ID
If you send new modal to slack, you must add your tool's `Tool ID` to modal's `external_id` field like below.

Format:
```
YOUR_TOOL_ID:ANYTHING
```

Example of Editor tool:
```
editor:ABCDEFGHIJK
```

SlackHub uses this value to identify the tool tied to the request. 

If you want to see an example of code, see the SlackHub's official tool `Editor`. It's typical advanced mode application written in Go.
