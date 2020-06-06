# STEP2 Register your own SlackHub App in slack console
In this step, you'll resister SlackHub to your slack workspace.

Unfortunately, there is NOT a way to do this process automatically. You must manually set up in slack console.

Oh! Don't drop your shoulders! Look on the bright side.
If you don't use SlackHub, all of your team members must do this process every time to make apps. Your help will save them a lot of work. You are like a HERO!

## Create an app
First, open [https://api.slack.com/apps](https://api.slack.com/apps) and click the "Create New App" button.

![Create New App](https://github.com/nicoJN/slackhub/blob/images/guide_admin_1.png)

Then enter the App name (SlackHub) and select Workspace.

Next, you should set your app an icon. Scroll down and click to upload.

![Set an icon](https://github.com/nicoJN/slackhub/blob/images/guide_admin_2.png)

I provide you an icon of SlackHub. See project root directory and get `logo_512.png`.

## Get your secter and register it to the SSM Parameter Store
Now that, it's a time to get your secret. Select `Settings - Basic Information` tab in the left sidebar and search `Signing Secret`.

NOTE: In this page, you can see many token and secrets. Please don't make a mistake!

Next, you should run the following command by using aws cli (replace placeholder to your secret). This command encrypt and register your signing secret to the AWS SSM Parameter Store.

```
$ aws ssm put-parameter --name "slackhub-signing-secret" --type "SecureString" --value "YOUR_SINGNING_SECRET_HERE!!!"
```

## Enable event subscriptions
SlackHub uses 2 functions of Slack.

- Event subscriptions
- Interactivity

First, you should enable event subscription. Select `Features - Event Subscriptions` tab in the left sidebar.

![Event Subscriptions](https://github.com/nicoJN/slackhub/blob/images/guide_admin_3.png)

Please set these 3 setting.

1. Turn the toggle to ON
2. Enter the Request URL field. Set `SlackHubStack.SlackHubInitialEndpointXXXXXXXXXX` you got in the STEP1.
3. Add `app_mention` event in the `Subscribe to bot events` section.

## Enable interactivity
Now, it's turn of interactivity. Select `Features - Interactivity & Shortcuts` tab in the left sidebar.

![Interactivity & Shortcuts](https://github.com/nicoJN/slackhub/blob/images/guide_admin_4.png)

Please set these 2 setting.

1. Turn the toggle to ON
2. Enter the Request URL field. Set `SlackHubStack.SlackHubInteractiveEndpointXXXXXXXXXX` you got in the STEP1.

## Get bot user access token and register it to the SSM Parameter Store
There are just a little more settings! Select `Features - OAuth & Permissions` tab in the left sidebar.

![OAuth & Permissions](https://github.com/nicoJN/slackhub/blob/images/guide_admin_5.png)

Please click `Install App to Workspace` button. 

After that, `Bot User OAuth Access Token` will be appeared, you should copy it and please excute a command below by using aws cli.

```
$ aws ssm put-parameter --name "slackhub-bot-user-auth-token" --type "SecureString" --value "YOUR_BOT_USER_AUTH_TOKEN_HERE!!!" 
```

## Add permission scopes to SlackHub
Now, it's last process! Keep to stay `Features - OAuth & Permissions` tab and scroll down.

![OAuth & Permissions - Scopes](https://github.com/nicoJN/slackhub/blob/images/guide_admin_6.png)

Please add these 3 permissin scopes in Scopes section.

1. `app_mentions:read`: This scope is auto matically added.
2. `chat:write`
3. `users:read`

That's all of setup process! You can call SlackHub by sending @mention to SlackHub. 
And also, your members will be able to resister their convinient tools to SlackHub.

Let's try to mention to the SlackHub and select the tool `SlackHub - Catalog`. You can experience a fundamental modal view application :+1:

Have a good day!
