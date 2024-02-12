# DiscordAI

A Discord bot powered by artificial intelligence, rebuilt in Go from [DiscordGPT](https://github.com/nekogravitycat/DiscordGPT). Designed for modularity and extensibility.

## Environment Variables

This program requires the following environment variables to be set in order to function properly:

- **DISCORDBOT_TOKEN**: Discord bot token for authentication.
- **OPENAI_TOKEN**: OpenAI API token for accessing AI capabilities.

You can set these environment variables by either creating a `.env` file at the project root and specifying them there, or by setting them directly using system environment variables.

## Commands

### Regular User Commands:
- **/activate-gpt**: Start ChatGPT on the current channel.
- **/deactivate-gpt**: Stop ChatGPT on the current channel.
- **/credits**: Check your remaining credits.
- **/add-gpt-image**: Add an image URL for GPT Vision to analyze.
  - `image-url`: The URL of the image (required).
- **/gpt-sys-prompt**: View the current GPT system prompt for this channel.
- **/set-gpt-sys-prompt**: Set the GPT system prompt for this channel.
  - `sys-prompt`: The system prompt of GPT (required).
- **/reset-gpt-sys-prompt**: Reset the GPT system prompt for this channel to default.
- **/set-gpt-model**: Set the GPT model for your usage.
  - `model`: The GPT model to use (required). Available options: `gpt-3.5-turbo`, `gpt-4-turbo-preview`, `gpt-4-vision-preview`.
- **/clear-gpt-history**: Clear the GPT chat history for this channel.

### Admin Commands:
- **/add-credits**: Add credits for a user.
  - `user-id`: User ID (required).
  - `amount`: Amount to add (in USD) (required).
- **/set-user-privilege**: Set privilege level for a user.
  - `user-id`: User ID (required).
  - `privilege-level`: Privilege level (required).

Here's the addition to the README regarding deploying the app with Docker using the provided docker-compose file:

## Docker Deployment

To deploy DiscordAI using Docker, you can utilize the provided docker-compose file. Ensure you have Docker installed on your system before proceeding.

1. Clone the GitHub repository for the docker-compose file:
    ```bash
    git clone https://github.com/nekogravitycat/DiscordAI-Server
    ```

2. Navigate to the directory containing the docker-compose file:
    ```bash
    cd DiscordAI-Server
    ```

3. Create and modify the `.env` file within the project directory to include your Discord bot token and OpenAI API token.

4. Run the following command to start DiscordAI:
    ```bash
    docker-compose up -d
    ```

Please note that the docker-compose file mounts the `configs` and `data` directories to persist configuration and data between container restarts. 