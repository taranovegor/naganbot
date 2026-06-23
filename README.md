# naganbot
Shoot yourself in Russian roulette

## Set up
### Requirements
- Docker, docker-compose
- PostgreSQL 18 or MySQL 8/MariaDB 11

### Configuration
Copy an instance of the environment file and save it as a file `.env`
```shell
cp .env.dist .env
```
... and configure as you need. Environment variables are described in comments

Build or pull docker images of application
```shell
make container-build
```
```shell
make container-pull
```

## Randomness & Verification

The bot uses [drand](https://drand.love) as its source of randomness — a decentralized, publicly verifiable random beacon operated by Cloudflare, Protocol Labs, and others.

### How it works

1. When a game is played, the bot fetches the latest beacon from drand's **quicknet** network (rounds every 3 seconds)
2. The beacon returns a cryptographically signed `randomness` value
3. The bot computes `SHA-256(randomness + game_id)` and uses the result to determine the outcome:
   - Bytes `[0:4]` → `% 100` → bullet type (< 3 = atomic, otherwise lead)
   - Bytes `[4:12]` → `% N` → victim index (ignored for atomic bullet — everyone dies)
4. A proof link is stored with the game and can be viewed via `/ngjoined`

### Verification

After a game is played, `/ngjoined` includes a link to the drand beacon used. To manually verify:

1. Open the proof link — note the `randomness` hex value and the `game[id]` query parameter
2. Decode `randomness` from hex to bytes
3. Compute `SHA-256(randomness_bytes + game_id_string)`
4. First 4 bytes of hash as big-endian uint32 `% 100` → bullet type
5. Bytes `[4:12]` of hash as big-endian uint64 `% N` (number of players) → victim index

The randomness is generated independently by the drand network before the game requests it — the bot cannot influence the outcome.

### Launch
Run the application using the Make tool
```shell
make start
```
