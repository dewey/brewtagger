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

2. Install brewtagger:
```bash
brew install dewey/brewtagger/brewtagger
```

This will:
- Install the binary
- Set up the launchd service to run at login
- Configure it to watch for new applications

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

## Automatic Tagging

The tool will automatically run:
- When you log in
- Every 6 hours
- When you install new Homebrew Cask applications

## Uninstallation

```bash
brew uninstall brewtagger
```

This will:
- Unload the launchd service
- Remove the plist file
- Remove the binary

## License

MIT License 