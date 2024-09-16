# Code Crafters - Kafka Go

[![progress-banner](https://backend.codecrafters.io/progress/kafka/b39a0c02-1cbb-4f05-a2a1-568f670f29f2)](https://app.codecrafters.io/users/codecrafters-bot?r=2qF)

## Learnings

- The `listener.Accept` is **blocking** meaning you most likely want to call it within its own Go routine.

- 8 bits = 1 byte.

- There is _big endian_ and _small endian_ notation.

  - In the _big endian_ notation, **the most significant byte is first**.

  - In the _small endian_ notation, **the most significant byte is last**.
