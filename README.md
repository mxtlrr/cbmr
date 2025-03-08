# cbmr
A server for SCP: Containment Breach ranked matches. Handles match data, including
- whether or not a player has forfeited / made a draw
- generating the seed for the specific category
- sending data back to the ranked client via HTTP requests
- storing player data in a database

## Architecture
![](./img/architecture.png)


# Usage of `cbmr`
TODO.

<!-- TODO: modernc.org/sqlite for sqlite -->

# Upcoming Features
- [ ] Add skill-based matchmaking (based off ELOs).