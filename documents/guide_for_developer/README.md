# Guide for developers
Hi! Thank you for your interest!

By using SlackHub, you can easily make an app with a modal view like below.

![SlackHub application](https://github.com/nicoJN/slackhub/blob/images/flow.png)

Or you can also make an app without a modal view.

![SlackHub application](https://github.com/nicoJN/slackhub/blob/images/flow_without_modal.png)

Almost all you have to do in this flow is **ONLY implementing tool's lambda!**. Let's take a look at how to make your app.

## Overview
![Role of SlackHub](https://github.com/nicoJN/slackhub/blob/images/role.png)

The overview of mechanism is that SlackHub works like a proxy and manage almost all slack apps lifecycle event instead. (message authorization, send a selectable menu, send a modal view and so on ...)

If you are interested in SlackHub's architecture and want to know more, see the page of [slackhub's architeture](https://github.com/nicoJN/slackhub/blob/master/documents/slackhub_architecture). 

---
You can easily make a tool with **3 steps**.

1. [Create a modal view with slack's official GUI tool](https://github.com/nicoJN/slackhub/blob/master/documents/guide_for_developer/step1_create_modal_view)
2. [Implemete tool's lambda in your favorite language!](https://github.com/nicoJN/slackhub/blob/master/documents/guide_for_developer/step2_implement_lambda)
3. [Register your tool with SlackHub's Register tool](https://github.com/nicoJN/slackhub/blob/master/documents/guide_for_developer/step3_register_your_tool)

Please see each page!
