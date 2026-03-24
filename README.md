# heos-cli

CLI for controlling Denon HEOS speakers over the local network.

## Install

Download a binary from the [latest release](https://github.com/jrogala/heos-cli/releases/latest), or install with Go:

```bash
go install github.com/jrogala/heos-cli@latest
```

## Setup

Set `HEOS_HOST` env var or configure via `~/.config/heos-cli/config.yaml`:

```yaml
host: 192.168.1.x
```

Run `heos setup` to auto-discover speakers on the network.

## Commands

| Command | Description |
|---|---|
| `setup` | Auto-discover HEOS speakers and save config |
| `player list` | List all available players |
| `player info --pid PID` | Show player details |
| `player play --pid PID` | Start playback |
| `player pause --pid PID` | Pause playback |
| `player stop --pid PID` | Stop playback |
| `player get-volume --pid PID` | Get player volume |
| `player set-volume --pid PID --level N` | Set volume (0-100) |
| `player mute --pid PID` | Mute player |
| `player unmute --pid PID` | Unmute player |
| `player now-playing --pid PID` | Show currently playing track |
| `player queue --pid PID` | List the play queue |
| `player next --pid PID` | Skip to next track |
| `player previous --pid PID` | Skip to previous track |
| `group list` | List speaker groups |
| `group set-volume --gid GID --level N` | Set group volume |
| `browse sources` | List available music sources |
| `browse play-url --pid PID --url URL` | Play a URL on a player |
| `system heart-beat` | Check connection to HEOS system |
| `system check-account` | Show signed-in HEOS account |

## Examples

```bash
$ heos player list
PID  NAME         MODEL           NETWORK  SERIAL
1    Living Room  Denon Home 150  wifi     AABBCC
2    Kitchen      Denon Home 150  wifi     DDEEFF

$ heos player get-volume --pid 1
40

$ heos player now-playing --pid 1
Type:    song
Song:    Bohemian Rhapsody
Artist:  Queen
Album:   A Night at the Opera

$ heos player set-volume --pid 1 --level 30
Volume set to 30
```

## JSON Output

All commands support `--json` for machine-readable output.
