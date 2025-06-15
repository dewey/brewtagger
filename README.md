# AppTagger

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
go build -o apptagger
```

3. Move the binary to a suitable location:
```bash
sudo mv apptagger /usr/local/bin/
```

## Usage

Run the tool manually:
```bash
apptagger                    # Tag apps with default color (Yellow)
apptagger -tag-color=Blue    # Tag apps with Blue color
apptagger -tag-color=Red     # Tag apps with Red color
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
sudo mkdir -p /Library/LaunchDaemons
sudo cp com.github.dewey.apptagger.plist /Library/LaunchDaemons/
```

2. Edit the plist file to point to your binary location:
```bash
sudo sed -i '' "s|/path/to/apptagger|/usr/local/bin/apptagger|g" /Library/LaunchDaemons/com.github.dewey.apptagger.plist
```

3. Load the launchd job:
```bash
sudo launchctl load /Library/LaunchDaemons/com.github.dewey.apptagger.plist
```

The tool will now run:
- When you log in
- Every 6 hours
- When you install new Homebrew Cask applications

## Uninstallation

1. Unload the launchd job:
```bash
sudo launchctl unload /Library/LaunchDaemons/com.github.dewey.apptagger.plist
```

2. Remove the files:
```bash
sudo rm /Library/LaunchDaemons/com.github.dewey.apptagger.plist
sudo rm /usr/local/bin/apptagger
```

## License

MIT License 