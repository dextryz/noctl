# Ixian

[Ixian](https://dune.fandom.com/wiki/Ix) inspired by the Dune universe is a [Nostr](https://nostr.com) commandline tool done the [Unix way](https://en.wikipedia.org/wiki/Unix_philosophy).

## Create Profile

1. Create a new profile configuration.

```shell
mkdir -p $HOME/.config/nostr
export CONFIG_NOSTR=$HOME/.config/nostr/alice.json
touch $CONFIG_NOSTR
```

2. Generate a private-public key pair with the [cipher](https://github.com/ffiat/cipher) commandline tool.

```shell
cipher -keygen
```

3. Add this generated secret (`nsec`) to your local profile.

```shell
ix profile -nsec <nsec...>
ix profile -name "Alice"
```

5. Before you can fetch notes you have to add at least one relay.

```
ix relay -add wss://nostr.ffiat.net
```

6. Now commit your profile to these relays.

```shell
ix profile -commit
```

## Publish Notes

1. Publish your first text note.

```shell
ix event -note "hello friend"
```

2. Request the note your just published using your npub.

```shell
ix req -npub <npub...>
```

3. Request a note via a specific event ID.

```shell
ix req -id <eventId>
```

## Show Timeline

1. Add users to follow, including yourself

```
ix follow -add <npub>
```

7. Finally, echo your timeline of events pulled from your set of relays and followings.

```shell
ix req -following
```
