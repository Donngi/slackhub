# STEP3 Register your tool with SlackHub's Register tool
Hi, it's a final step!
Let's register your tool!

## Register
First, you should call SlackHub in your slack channel and select the tool `SlackHub - Resister`.

![SlackHub - Resister](https://github.com/nicoJN/slackhub/blob/images/guide_developer_3.png)

You can register your tool by only using this tool. Please enter the each field.

| Field | Detail |
| :--- | :--- |
| Tool ID | ID of the tool. No duplication allowed. |
| Display Name | The name displayed for users. |
| Description | The description of your tool. It will be appear when users use Catalog tool. |
| Arn of callee lambda | Your lambda's arn like `arn:aws:lambda:YOUR_REGION:YOUR_ACCOUNT:function:XXXXXXXXXXXX` |
| Administrators | If you select administrators, no one else will be able to edit the tool. Default is none. |
| Authorized Users | If you select authorized users, no one else will be able to use the tool. Default is none. |
| Modal JSON | Your modal JSON created in STEP1. If you don't register any modal, SlackHub won't send modal to Slack and will simply pass the request to your tool's lambda. | 
| Boot Mode | This is **Advanced Setting**. Default value is "Normal" and I strictly recommend you to keep "Normal" if you are NOT familiar with Slack's lifecycle sequence. When you set "Advanced", SlackHub will call your lambda synchronously, and you must manage all request/response by yourself. |

After entering all fields, push `Submit` button. 

That's all of registering. You can use your tool through SlackHub!

## After registering ...
If you want to edit, delete and list your tools, you can use SlackHub's official tools below.

| Tool Name | Function |
| :--- | :--- |
| Editor | Edit your existing tools |
| Eraser | Delete your existing tools |
| Catalog | List your tools |
