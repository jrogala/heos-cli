# heos-cli

CLI for controlling Denon HEOS speakers over the local network.

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
| `player info` | Show player details |
| `player play` | Start playback |
| `player pause` | Pause playback |
| `player stop` | Stop playback |
| `player volume` | Get or set volume (0-100) |
| `player mute` | Toggle mute |
| `player now-playing` | Show currently playing track |
| `player queue` | List the play queue |
| `player next` | Skip to next track |
| `player prev` | Skip to previous track |
| `group list` | List speaker groups |
| `group set` | Create or modify a speaker group |
| `group volume` | Get or set group volume |
| `group mute` | Toggle group mute |
| `browse sources` | List available music sources |
| `browse play-url` | Play a URL on a player |
| `browse search` | Search music sources |
| `system heart-beat` | Check connection to HEOS system |
| `system check-account` | Show signed-in HEOS account |
| `system sign-in` | Sign in to HEOS account |

## Examples

```bash
# Discover speakers and save config
heos setup

# Set living room speaker to 40% volume
heos player volume --player "Living Room" 40

# Play a URL on a specific player
heos browse play-url --player "Living Room" "http://stream.example.com/radio"

# Group two speakers together
heos group set --leader "Living Room" --members "Kitchen"

# See what's currently playing
heos player now-playing --player "Living Room"
```

## JSON Output

All commands support `--json` for machine-readable output.
