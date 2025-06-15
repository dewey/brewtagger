class Brewtagger < Formula
  desc "Automatically tag Homebrew Cask applications in your /Applications folder"
  homepage "https://github.com/dewey/brewtagger"
  url "https://github.com/dewey/brewtagger/archive/refs/tags/v0.1.0.tar.gz"
  sha256 "" # We'll need to fill this in after creating the release
  license "MIT"

  depends_on "go" => :build
  depends_on "tag"

  def install
    system "go", "build", "-o", "brewtagger"
    bin.install "brewtagger"
    
    # Install the plist file
    prefix.install "com.github.dewey.brewtagger.plist"
  end

  def post_install
    # Create LaunchAgents directory if it doesn't exist
    system "mkdir", "-p", "#{ENV["HOME"]}/Library/LaunchAgents"
    
    # Copy the plist file
    plist_path = "#{ENV["HOME"]}/Library/LaunchAgents/com.github.dewey.brewtagger.plist"
    system "cp", "#{prefix}/com.github.dewey.brewtagger.plist", plist_path
    
    # Load the launchd service
    system "launchctl", "load", plist_path
  end

  def uninstall
    # Unload the launchd service
    plist_path = "#{ENV["HOME"]}/Library/LaunchAgents/com.github.dewey.brewtagger.plist"
    system "launchctl", "unload", plist_path if File.exist?(plist_path)
    
    # Remove the plist file
    system "rm", "-f", plist_path
  end

  test do
    system "#{bin}/brewtagger", "--version"
  end
end 