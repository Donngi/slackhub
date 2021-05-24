# STEP1 Create a modal view with slack's official GUI tool
Hi!
In this step, we make a modal view with a GUI tool.

To make a modal view more easily, slack team officially provides an GUI tool :tada:
[Block Kit Builder](https://app.slack.com/block-kit-builder)

It's a powerful and easy to use. Umm, that's why Slack is loved by everyone :+1:

(If you don't need to use a modal view, you can skip this step and see [Step2](https://github.com/Jimon-s/slackhub/blob/master/documents/guide_for_developer/step2_implement_lambda))

## Create basis of a modal view
First, you should visit [Block Kit Builder site](https://api.slack.com/tools/block-kit-builder) and make a basis of a modal view.

![Block Kit Builder](https://github.com/Jimon-s/slackhub/blob/images/guide_developer_1.png)

You can locate blocks by GUI and json code will be automatically created. 

## Insert block_id and action_id
Only using Block Kit Builder, you can get a complete json of a modal view. However, as it is, there are NOT information to identify blocks. 

For example, if your modal view have five text input areas, you can't determine which information is tied to a particular input area.

In order to avoid such a situation, slack provides us functions of **block_id** and **action_id**. We can label each block by using them.

Of course, you don't have to label all blocks. You should insert block_id and action_id **only to blocks which you will use in your tool's implemantation**.

![Which is block](https://github.com/Jimon-s/slackhub/blob/images/guide_developer_2.png)

NOTE: In the modal view, block is a unit of a component. It's correspond to the item which you can select from a left sidebar. You should add block_id to them. In addition, you should also add action_id to the element in the block.

| ID | Target to be added |
| :--- | :--- |
| block_id | each component in the "blocks" array |
| action_id | map "element" in the block |

I definitely recommend you to these approach on [Block Kit Builder](https://api.slack.com/tools/block-kit-builder) because it automatically points out mistakes. It's really useful.

### Example: Inputs - Singleline
Original:
```json
{
    "type": "input",
    "element": {
        "type": "plain_text_input"
    },
    "label": {
        "type": "plain_text",
        "text": "What do you want?",
        "emoji": true
    }
}
```

Edited (Add block_id and action_id)
```json
{
    "block_id": "YOUR_BLOCK_ID_HERE!!!",
    "type": "input",
    "element": {
        "action_id": "YOUR_ACTION_ID_HERE!!!",
        "type": "plain_text_input"
    },
    "label": {
        "type": "plain_text",
        "text": "What do you want?",
        "emoji": true
    }
}
```

### Example: Inputs - Multiline
Original:
```json
{
    "type": "input",
    "label": {
        "type": "plain_text",
        "text": "What do you want?",
        "emoji": true
    },
    "element": {
        "type": "plain_text_input",
        "multiline": true
    },
    "optional": true
}
```

Edited (Add block_id and action_id)
```json
{
    "block_id":"YOUR_BLOCK_ID_HERE!!!",
    "type": "input",
    "label": {
        "type": "plain_text",
        "text": "What do you want?",
        "emoji": true
    },
    "element": {
        "action_id":"YOUR_ACTION_ID_HERE!!!",
        "type": "plain_text_input",
        "multiline": true
    },
    "optional": true
}
```
