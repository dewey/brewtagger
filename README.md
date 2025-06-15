# BrewTagger

A simple tool that automatically tags Homebrew Cask applications in your `/Applications` folder. It uses the `tag` command-line tool to add color tags to your applications, making it easy to visually distinguish between different types of apps.

## Prerequisites

- macOS
- [Homebrew](https://brew.sh/) installed
- [tag](https://github.com/jdberry/tag) command-line tool installed

## Installation

1. Install the required dependencies:
```bash
brew install tag
```

2. Build the tool:
```bash
go build -o brewtagger
```

3. Move the binary to a suitable location:
```bash
sudo mv brewtagger /usr/local/bin/
```

## Usage

Run the tool manually:
```bash
brewtagger                    # Tag apps with default color (Yellow)
brewtagger -tag-color=Blue    # Tag apps with Blue color
brewtagger -tag-color=Red     # Tag apps with Red color
```

### Available Colors

The tool supports all colors that the `tag` command-line tool supports. Common colors include:
- Yellow (default)
- Blue
- Green
- Red
- Purple
- Orange
- Gray

## Automatic Tagging with launchd

To automatically tag apps when they're installed or when you log in, you can use launchd:

1. Create the launchd plist file:
```bash
mkdir -p ~/Library/LaunchAgents
cp com.github.dewey.brewtagger.plist ~/Library/LaunchAgents/
```

2. Edit the plist file to point to your binary location:
```bash
sed -i '' "s|/path/to/brewtagger|/usr/local/bin/brewtagger|g" ~/Library/LaunchAgents/com.github.dewey.brewtagger.plist
```

3. Load the launchd job:
```bash
launchctl load ~/Library/LaunchAgents/com.github.dewey.brewtagger.plist
```

The tool will now run:
- When you log in
- Every 6 hours
- When you install new Homebrew Cask applications

## Uninstallation

1. Unload the launchd job:
```bash
launchctl unload ~/Library/LaunchAgents/com.github.dewey.brewtagger.plist
```

2. Remove the files:
```bash
rm ~/Library/LaunchAgents/com.github.dewey.brewtagger.plist
sudo rm /usr/local/bin/brewtagger
```

## License

MIT License 