# DiscordAI

A Discord bot powered by artificial intelligence, rebuilt in Go from [DiscordGPT](https://github.com/nekogravitycat/DiscordGPT). Designed for modularity and extensibility.

## Commands

### Regular Commands

- `/gpt activate`: Start ChatGPT on this channel.
- `/gpt deactivate`: Stop ChatGPT on this channel.
- `/gpt set-model`: Set the GPT model for the user.
- `/gpt clear-history`: Clear GPT chat history for this channel.
- `/gpt sys-prompt show`: Show GPT system prompt for this channel.
- `/gpt sys-prompt set`: Set GPT system prompt for this channel.
- `/gpt sys-prompt reset`: Reset GPT system prompt for this channel to default.
- `/credits`: Check user credits.
- `/dall-e-2-generate`: Generate an image using DALL·E 2.
- `/dall-e-3-generate`: Generate an image using DALL·E 3.

## Admin Commands

- `/add-credits`: Add credits for a user.
- `/set-user-privilege`: Set privilege level for the user.


## Environment Variables

This program requires the following environment variables to be set in order to function properly:

- **DISCORDBOT_TOKEN**: Discord bot token for authentication.
- **OPENAI_TOKEN**: OpenAI API token for accessing AI capabilities.

You can set these environment variables by either creating a `.env` file at the project root and specifying them there, or by setting them directly using system environment variables.

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
    docker compose up -d
    ```

Please note that the docker-compose file mounts the `configs` and `data` directories to persist configuration and data between container restarts.
