![SlackHub](https://github.com/nicoJN/slackhub/blob/images/banner.png)

ChatOps acceleration tool for the team! Easy to create slack modal applications.
===

SlackHub enables your team to create **slack modal application** like below.

![Demo](https://github.com/nicoJN/slackhub/blob/images/demo.gif)

## Features
**Make it easy to operate your systems**. It's truly what you want to do. 

However, it takes a long time to be able to create ChatOps tools itself. SlackHub supports people who have such concerns.

By using SlackHub, you can easily make an app with a modal view like below.

![SlackHub application](https://github.com/nicoJN/slackhub/blob/images/flow.png)

Or you can also make an app without a modal view.

![SlackHub application](https://github.com/nicoJN/slackhub/blob/images/flow_without_modal.png)

Almost all you have to do in this flow is **ONLY implementing tool's lambda which conducts what you really want to do!**. SlackHub manages almost all processes between slack and your application server. 

### Important points

- `Easy to implement`: You can make tools by simply implementing your lambda and making a modal view appearance with GUI tool.
- `Low cost`: It's fully composed of AWS serverless components, so running costs is extreamly low!
- `Authentication`: SlackHub provides an authentication function. You can easily limit the users who are able to use each tool.
- `Official Tools`: We initially provide you some tools like new tool register, editor, catalog and eraser. You can manage your tools by only ChatOps.

### Architecture Overview
![Architecture overview](https://github.com/nicoJN/slackhub/blob/images/architecture.png)

## Instration
1. [SlackHub initial setup](https://github.com/nicoJN/slackhub/blob/master/documents/guide_for_admin) (Only once! By your team's representative.)
2. [Create your tools!](https://github.com/nicoJN/slackhub/blob/master/documents/guide_for_developer)

We provide you some guides with screenshots :+1:

## Contribution
You are more than welcome to contribute to this project. Fork and make a Pull Request, or create an Issue if you see any problem. 

Please run the following command before you make pull request.

```
make pr-prep
```

I'm a biginner of Go and not good at writing English. Refactoring of code and documents are much appreciated!

## License
MIT (NOTE: You must not use SlackHub's logo and some images to other things which are not related to SlackHub)

This project uses some repositories. See [dependencies](https://github.com/nicoJN/slackhub/blob/master/dependency) directory.
