# How to Create a Telegram Bot

## 1. Create a Telegram Bot

1. Search for the `@BotFather` username in your Telegram application.
2. Click "Start" to begin a conversation with `@BotFather`.
3. Send the command `/newbot` to `@BotFather`.
4. Provide a name and a unique username for your bot.
5. BotFather will then give you a `botToken`. Save this token as it will be used to configure Alertmanager.

## 2. Create a Group for Receiving Alerts and Add the Bot

1. Create a new group in your Telegram application. This group will be used for receiving alerts.
2. Add your bot to the group you just created.
3. Send a message in the group, such as `/hello @your_bot_username`.
4. Open a web browser and go to `https://api.telegram.org/bot<your-bot-token>/getUpdates`, replacing `<your-bot-token>` with the token you received from @BotFather.
5. After sending a message to the group, refresh the browser page. Look for the chat JSON object in the response to find the group's `chat_id`. This ID will be a negative number and is required for Alertmanager configuration.

## 3. Configure Alertmanager

Now that you have the `botToken` and the group `chat_id`, you can use them to configure Alertmanager by following this doc [How to Receive Alerts](how-to-receive-alerts.md).
