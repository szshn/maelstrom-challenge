# Distributed Systems Challenges

This repository contains my Go implementations for the **Gossip Glomers** distributed systems challenge series by [Fly.io](https://fly.io/dist-sys/), in collaboration with Kyle Kingsbury, creator of [Jepsen](https://github.com/jepsen-io/jepsen).  

These challenges are built on top of the **Maelstrom** testing framework, which simulates network behavior, injects failures, and verifies system properties based on specific consistency guarantees.

## 💡 Challenges Covered
1. [Echo](/maelstrom-echo)
1. [Unique ID Generation](/maelstrom-unique-ids)
1. [Broadcast](/maelstrom-broadcast/)
1. Grow-Only Counter
1. Kafka-Style Log
1. Totally-Available Transactions

Each challenge explores key concepts in distributed systems such as fault tolerance, consistency, and replication.

## ⚙️ Tech Stack
- **Language:** Go
- **Framework:** Maelstrom (Jepsen workbench)

## 📚 Resources
- [Fly.io Challenge Docs](https://fly.io/dist-sys/)
- [Maelstrom GitHub Repository](https://github.com/jepsen-io/maelstrom)
- [Jepsen Testing Library](https://github.com/jepsen-io/jepsen)

## 🗂️ Repo Structure
Each challenge is implemented in its own directory with documentation and test scripts for Maelstrom verification.

## 📄 License & Attribution
This project is licensed under the MIT License. Challenge descriptions and Maelstrom framework are provided by Fly.io and Jepsen respectively.
