import logging
import urllib.request
import json

logger = logging.getLogger(__name__)


def lambda_handler(event, context):

    # Get user input values.
    lunch_input = event["view"]["state"]["values"]["lunch_block"]["lunch_action"]["value"]
    detail_input = event["view"]["state"]["values"]["detail_block"]["detail_action"]["value"]

    # Get private_metadata fields value.
    # NOTE: private_metadata field's type is simply "String". SlackHub stores some convenient values in the form of JSON.
    # In order to use them easily, you should convert private_metadata field value to JSON.
    private_metadata_str = event["view"]["private_metadata"]
    private_metadata = json.loads(private_metadata_str)
    response_url = private_metadata["response_url"]

    # Send a message to Slack.
    msg = ("What is your favorite lunch?\n - {}"
           "\nTell us more!\n - {}").format(lunch_input, detail_input)
    params = {
        "text": msg
    }

    headers = {
        "Content-Type": "application/json"
    }

    req = urllib.request.Request(
        response_url, json.dumps(params).encode(), headers)

    try:
        urllib.request.urlopen(req)
    except urllib.error.HTTPError as err:
        # If status code is 4XX or 5XX
        logger.error(err.code)
    except urllib.error.URLError as err:
        # If HTTP communication somehow failed
        logger.error(err.reason)

    return
