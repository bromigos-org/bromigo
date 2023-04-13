import { Configuration, OpenAIApi } from "openai";

const configuration = new Configuration({
    apiKey: process.env.OPENAI_API_KEY,
});
const openai = new OpenAIApi(configuration);

const createChatCompletion = async (requestText: string): Promise<string> => {
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

    return response.data.choices[0]?.message?.content!;

}

export { createChatCompletion }