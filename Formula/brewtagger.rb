# typed: true
# frozen_string_literal: true

# This file was generated by Homebrew Releaser. DO NOT EDIT.
class Brewtagger < Formula
  desc "Tag casks with a color in /applications so we know they are homebrew managed"
  homepage "https://github.com/dewey/brewtagger"
  url "https://github.com/dewey/brewtagger/archive/refs/tags/v0.3.0.tar.gz"
  sha256 "e8dd9e460fbe2d8ad11cc99fef5594d27eb0bc5b4f5b6b3d68ca3ea42f62694b"

  depends_on "go" => :build
  depends_on "tag"

  def install
    system "go", "build", "-o", "brewtagger"
    bin.install "brewtagger"
    prefix.install "com.github.dewey.brewtagger.plist"
  end

  test do
    assert_match("brewtagger version", shell_output("brewtagger --version"))
  end
end
