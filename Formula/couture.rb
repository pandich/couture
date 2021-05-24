# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Couture < Formula
  desc "Allows for tailing multiple event sources."
  homepage ""
  version "0.0.4"
  bottle :unneeded

  if OS.mac? && Hardware::CPU.intel?
    url "https://github.com/gaggle-net/couture/releases/download/v0.0.4/couture_0.0.4_Darwin_x86_64.tar.gz", :using => GitHubPrivateRepositoryReleaseDownloadStrategy
    sha256 "857ab40460fdf4b6fcb24f21fdff88988690243d68da3a93ef94b781170c69f2"
  end
  if OS.mac? && Hardware::CPU.arm?
    url "https://github.com/gaggle-net/couture/releases/download/v0.0.4/couture_0.0.4_Darwin_arm64.tar.gz", :using => GitHubPrivateRepositoryReleaseDownloadStrategy
    sha256 "a7a8d0599b9a98ba9abb6cecc9a56d39213fa97d57d25aa531eb4b340c96c747"
  end
  if OS.linux? && Hardware::CPU.intel?
    url "https://github.com/gaggle-net/couture/releases/download/v0.0.4/couture_0.0.4_Linux_x86_64.tar.gz", :using => GitHubPrivateRepositoryReleaseDownloadStrategy
    sha256 "cff8b909f897a0f181ce2e9f2213126d5e699c9f9102caad91fdfc9507132d3c"
  end
  if OS.linux? && Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
    url "https://github.com/gaggle-net/couture/releases/download/v0.0.4/couture_0.0.4_Linux_arm64.tar.gz", :using => GitHubPrivateRepositoryReleaseDownloadStrategy
    sha256 "6807f5619ff87078b58a31e9b946d2181325168f6e6f1ea7312e625c69cc6c3e"
  end

  def install
    bin.install "couture"
  end
end
