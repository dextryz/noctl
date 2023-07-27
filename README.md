# Melange

Nostr terminal client the Unix way.

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
melange profile -nsec <nsec...>
melange profile -name "Alice"
```

5. Before you can fetch notes you have to add at least one relay.

```
melange relay -add wss://nostr.ffiat.net
```

6. Now commit your profile to these relays.

```shell
melange profile -commit
```

## Publish Notes

1. Publish your first text note.

```shell
melange event -note "hello friend"
```

2. Request the note your just published using your npub.

```shell
melange req -npub <npub...>
```

## Show Timeline

1. Add users to follow, including yourself

```
melange follow -add <npub>
```

7. Finally, echo your timeline of events pulled from your set of relays and followings.

```shell
melange req -following
```
