require('dotenv').config()
import { Client, GatewayIntentBits, Message } from 'discord.js';
import axios from 'axios';

const client = <any>new Client({ intents: [GatewayIntentBits.Guilds, GatewayIntentBits.GuildMessages, GatewayIntentBits.MessageContent] });
const prefix = '<@1069873816907558922>'

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
    const apiKey = process.env.OPENAI_API_KEY;
    
    try {
        const response = await axios.post('https://api.openai.com/v1/completions', {
            model: "text-davinci-003",
            prompt: requestText,
            max_tokens: 4000,
        }, {
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${apiKey}`,
            },
        });

        msg.reply(response.data.choices[0].text);
    } catch (error) {
        console.error(error);
    }
});
client.login(process.env.DISCORD_API_TOKEN);
