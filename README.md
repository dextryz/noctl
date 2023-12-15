# NIX

Nostr Information eXtraction (NIX) via the Unix command line.

## Create Profile

1. Create a new profile configuration.

```shell
mkdir -p $HOME/.config/nostr
export CONFIG_NOSTR=$HOME/.config/nostr/alice.json
touch $CONFIG_NOSTR
```

2. Generate a private-public key pair with the [cipher](https://github.com/ffiat/cipher) commandline tool.

```shell
nix -keygen
```

3. Add this generated secret (`nsec`) to your local profile.

```shell
nix profile -nsec <nsec...>
nix profile -name "Alice"
```

5. Before you can fetch notes you have to add at least one relay.

```
nix relay -add wss://nostr.ffiat.net
```

6. Now commit your profile to these relays.

```shell
nix profile -commit
```

## Publish Notes

1. Publish your first text note.

```shell
nix event -note "hello friend"
```

2. Request the note your just published using your npub.

```shell
nix req -npub <npub...>
```

3. Request a note via a specific event ID.

```shell
nix req -id <eventId>
```

## Show Timeline

1. Add users to follow, including yourself

```
nix follow -add <npub>
```

7. Finally, echo your timeline of events pulled from your set of relays and followings.

```shell
nix req -following
```

## Key Management

```shell
> nix encode -note 9ccec662f0a0bb3e00231134af8e7222249073bd30896a62fc1fcd5de513f8ef
note1nn8vvchs5zanuqprzy62lrnjygjfquaaxzyk5churlx4megnlrhsf44pp7

> nix encode -npub a07b13f189d309c36a399d208208d09863b231c3745bd9d63d4de4339ab540a5
npub15pa38uvf6vyux63en5sgyzxsnp3myvwrw3dan43afhjr8x44gzjse75plm
```

```shell
> nix key -new
nsec: nsec1qfczae3envq7rldt8rpvtf99xu0rq7kq8ua26tw89r0c57722kuqhd4lsc
npub: npub1yejfe5ujj03elm62m9pek9wau8tp3qzd3cjvmn8gqlnx3846qz4s453w6w
```

```shell
> nix key -decode npub15pa38uvf6vyux63en5sgyzxsnp3myvwrw3dan43afhjr8x44gzjse75plm
a07b13f189d309c36a399d208208d09863b231c3745bd9d63d4de4339ab540a5
```
