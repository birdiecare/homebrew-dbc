# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Dbc < Formula
  desc "Database Connect"
  homepage "https://github.com/birdicare/homebrew-dbc"
  version "0.3.12"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/birdiecare/homebrew-dbc/releases/download/v0.3.12/birdiecare_dbc_0.3.12_darwin_arm64.tar.gz"
      sha256 "862bae0b4f3af57407dc28ca64e33f3bb2e813db78b6ca98127d2d43a97bec2c"

      def install
        bin.install "dbc"
      end
    end
  end

  on_linux do
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/birdiecare/homebrew-dbc/releases/download/v0.3.12/birdiecare_dbc_0.3.12_x86_64_arm64.tar.gz"
      sha256 "89f4920e9fd7dfba3c565e42d281be0431839bcc4c876c79d93374f1fc871262"

      def install
        bin.install "dbc"
      end
    end
  end
end
