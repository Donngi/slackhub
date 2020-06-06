# SlackHub Architecture

In this section, you can see SlackHub's architecture.

## Architecture overview
![Architecture overview](https://github.com/nicoJN/slackhub/blob/images/architecture.png)

SlackHub has two main streams in its architecture.
One is `Initial` stream which sends a tool list menu to slack. Another is `Interactive` stream which sends a modal view to slack and invokes tool's lambda.

## Sequence (Normal Mode)
If tool developers select `Boot Mode = Normal` (Default), the lifecycle of the tool is below.

![Sequence Normal Mode](https://github.com/nicoJN/slackhub/blob/images/sequence_normal_mode.png)

In Normal mode, SlackHub handles almost all events instead of you. First, when users send @mention, SlackHub receives the message and do

- Authorize the message
- Identify slack's event type
- Create a menu of tool
- Send the menu to slack

Then, users receive the menu and select the tool they want to use. Once users choice the tool, SlackHub receives the message and do

- Authorize the message
- Authorize the user whether he/she is permitted to use the seleceted tool
- Identify slack's event type
- Create a modal view tied to the tool (load metadata from DynamoDB)
- Send the modal view to slack

After that, users can see the pop up modal view and input some fields. Then SlackHub receives the message and do

- Authorize the message
- Identify slack's event type
- Invoke tool's lambda **asynchronously** with the JSON payload which is sent from slack
- Send a invocation message to slack

The reason why SlackHub sends a invocation message is to avoid abuse of tools. I think that all users in workspace should know when and who use the tool because the tools may contains some sensitive tools like deploy or data recovery tool.

After tool's lambda receive the JSON from SlackHub, tools are eble to do anything as tool's developers like. I strictly recommend developers to send a compression message in order to make more user friendly application.

## Sequence (Advanced Mode)
If tool developers select `Boot Mode = Advanced`, the lifecycle of the tool is below. The difference between normal mode and advance mode is that developers are able to controll all lifecycle events.

![Sequence Normal Mode](https://github.com/nicoJN/slackhub/blob/images/sequence_advanced_mode.png)

It is same as normal mode until users submit a modal view inputs. After receiving a submission event, SlackHub do

- Authorize the message
- Identify slack's event type
- Invoke tool's lambda **synchronously** with tha JSON payload which is sent from slack

In advanced mode, SlackHub acts as just a proxy. Developers must implement all request/response sequence by themselves.

To utilize advanced mode, developers are able to realize chain of modal views. One of example of the tool using advanced mode is SlackHub's official tool `Editor`. If you have any interests, please see the code.
