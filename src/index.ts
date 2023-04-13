require('dotenv').config()
import { Client, GatewayIntentBits, Message } from 'discord.js';
import { createChatCompletion } from './chatgpt';


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


    msg.reply(await createChatCompletion(requestText));
  } catch (error) {
    console.error(error);
  }
});
client.login(process.env.DISCORD_API_TOKEN);
