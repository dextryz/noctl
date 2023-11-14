# Ixian

[Ixian](https://dune.fandom.com/wiki/Ix) inspired by the Dune universe is a [Nostr](https://nostr.com) commandline tool done the [Unix way](https://en.wikipedia.org/wiki/Unix_philosophy).

## TODO

- [X] Request relay information via NIP-11

## Create Profile

1. Create a new profile configuration.

```shell
mkdir -p $HOME/.config/nostr
export CONFIG_NOSTR=$HOME/.config/nostr/alice.json
touch $CONFIG_NOSTR
```

2. Generate a private-public key pair with the [cipher](https://github.com/ffiat/cipher) commandline tool.

```shell
lemon -keygen
```

3. Add this generated secret (`nsec`) to your local profile.

```shell
lemon profile -nsec <nsec...>
lemon profile -name "Alice"
```

5. Before you can fetch notes you have to add at least one relay.

```
lemon relay -add wss://nostr.ffiat.net
```

6. Now commit your profile to these relays.

```shell
lemon profile -commit
```

## Publish Notes

1. Publish your first text note.

```shell
lemon event -note "hello friend"
```

2. Request the note your just published using your npub.

```shell
lemon req -npub <npub...>
```

3. Request a note via a specific event ID.

```shell
lemon req -id <eventId>
```

## Show Timeline

1. Add users to follow, including yourself

```
lemon follow -add <npub>
```

7. Finally, echo your timeline of events pulled from your set of relays and followings.

```shell
lemon req -following
```

## Key Management

```shell
> lemon encode -note 9ccec662f0a0bb3e00231134af8e7222249073bd30896a62fc1fcd5de513f8ef
note1nn8vvchs5zanuqprzy62lrnjygjfquaaxzyk5churlx4megnlrhsf44pp7

> lemon encode -npub a07b13f189d309c36a399d208208d09863b231c3745bd9d63d4de4339ab540a5
npub15pa38uvf6vyux63en5sgyzxsnp3myvwrw3dan43afhjr8x44gzjse75plm
```

```shell
> lemon key -new
nsec: nsec1qfczae3envq7rldt8rpvtf99xu0rq7kq8ua26tw89r0c57722kuqhd4lsc
npub: npub1yejfe5ujj03elm62m9pek9wau8tp3qzd3cjvmn8gqlnx3846qz4s453w6w
```

```shell
> lemon key -decode npub15pa38uvf6vyux63en5sgyzxsnp3myvwrw3dan43afhjr8x44gzjse75plm
a07b13f189d309c36a399d208208d09863b231c3745bd9d63d4de4339ab540a5
```
