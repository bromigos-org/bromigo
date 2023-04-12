require('dotenv').config()
import { Client, GatewayIntentBits, Message } from 'discord.js';
import { Configuration, OpenAIApi } from "openai";

const configuration = new Configuration({
  apiKey: process.env.OPENAI_API_KEY,
});
const openai = new OpenAIApi(configuration);

const client = <any>new Client({ intents: [GatewayIntentBits.Guilds, GatewayIntentBits.GuildMessages, GatewayIntentBits.MessageContent] });
const prefix = process.env.BOTID!;

client.on('ready', () => {
  console.log(
    `Logged in as ${client.user.tag} (${client.user.id})`
  );
  client.user.setActivity(
    `I am a Bot using ChatGPT`
  );
});


client.on('messageCreate', async (msg: Message) => {
  if (!msg.content.startsWith(prefix) || msg.author.bot) return;

  const args = <any>msg.content.slice(prefix.length).split(/ +/);
  const command = args[1].toLowerCase();

  const requestText = args.join(" ");

  try {
    const response = async (): Promise<{ choices: { message?: { content?: string } }[] }> => {
      const response = await openai.createChatCompletion({
        model: "gpt-3.5-turbo",
        messages: [
          { "role": "system", "content": "You are a discord bot that talks and acts like PC principal from South Park. You are an all knowing discord bot" },
          { "role": "user", "content": "What is something you would get mad at?" },
          { "role": "assistant", "content": "As PC Principal, I get pretty fired up about any kind of hate speech, discrimination, or prejudice. I have zero tolerance for any form of bigotry or intolerance, and will do everything in my power to make sure that this server is a safe and welcoming space for everyone, regardless of their race, gender, sexual orientation, religion, or any other identity. I also get pretty angry when people use offensive language or slurs, as those have no place in civil discourse. Additionally, any kind of bullying or harassment will not be tolerated under my watch. So let's all work together to keep this server a positive and respectful community!" },
          { "role": "user", "content": requestText }
        ],
        max_tokens: 1900,
      });

      return response.data;

    }

    msg.reply((await response()).choices[0]?.message?.content!);
  } catch (error) {
    console.error(error);
  }
});
client.login(process.env.DISCORD_API_TOKEN);
