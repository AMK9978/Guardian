# LLM Guardian Central Unit
A plug-and-play end-to-end LLM guardian for seamless integration.
The central unit is in charge of fanning out requests to the connected 
security plugins based on user, and collecting their responses to judge whether a prompt
is malicious or benign. You can also use the central unit as a rate limiter, authentication unit, or quota management.


## Get Started
Docker:
```
git clone git@github.com:LLMGuardian/Guardian.git
docker-compose up --build -d
```

## Architecture
The guardian system uses a micro-kernel architecture designed by the idea of extension in mind:

![Architecture](./docs/arch.jpg)


## Contribution
I am always welcome to your contributions. Open an issue and/or open a PR accordingly.
This system is designed and the central unit is developed by Amir Karimi (@amk9978).

