# numero

numero is a math library for Go. 

This README is largely a work in progress document. I am writing down my thoughts as I go along.

### motivation

I am largely just looking to do some recreational programming. I know that mature math libraries exist, but I want to build one myself and maybe, just maybe, I'll arrive at something useful for everybody else.

### packages to make

1. nlog - logger (TODO)
2. nrepl - repl (TODO)
3. nparser - equation parser (TODO)
4. nplot - graphing (TODO)
5. nmath - math functions (TODO)
6. nserver - server (TODO)
7. nutil - utils (TODO)
6. nweb - web interface (TODO)

### areas to be covered by nmath

1. function differentiation
2. function integration
3. matrices
4. vectors
5. statistics
6. probability

### interface

All of these packages listed above can act as useful components to be combined with natural language.

To begin with, there can be a web based tool wherein you could input something like: "Solve " and then press "@" which pops up a little menu, you pick "equation" from that menu, write your equation. And the equation parsing does not happen via an AI. It happens via the tools we build. So there is determinism. This could be extended to a usecase where I would want to say "draw y = x". And yes, an AI could do this. But all of it takes tokens, API keys, guardrails and what not. Who wants that? 
