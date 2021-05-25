# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Couture < Formula
  desc "Allows for tailing multiple event sources."
  homepage ""
  version "0.0.10"
  bottle :unneeded

  if OS.mac? && Hardware::CPU.intel?
    url "https://github.com/gaggle-net/couture/releases/download/v0.0.10/couture_0.0.10_Darwin_x86_64.tar.gz", :using => GitHubPrivateRepositoryReleaseDownloadStrategy
    sha256 "9c976e7f11337781903f426ea252581585dcfcd43a00956898cccd70bddeb4dd"
  end
  if OS.mac? && Hardware::CPU.arm?
    url "https://github.com/gaggle-net/couture/releases/download/v0.0.10/couture_0.0.10_Darwin_arm64.tar.gz", :using => GitHubPrivateRepositoryReleaseDownloadStrategy
    sha256 "fd0e192e37d92d378fdecf9ca51788ec307717897d50be2695d16707e01a18cb"
  end
  if OS.linux? && Hardware::CPU.intel?
    url "https://github.com/gaggle-net/couture/releases/download/v0.0.10/couture_0.0.10_Linux_x86_64.tar.gz", :using => GitHubPrivateRepositoryReleaseDownloadStrategy
    sha256 "ded139591fef79478651131e1c4ce009027da1572f7705067edbca21b131b286"
  end
  if OS.linux? && Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
    url "https://github.com/gaggle-net/couture/releases/download/v0.0.10/couture_0.0.10_Linux_arm64.tar.gz", :using => GitHubPrivateRepositoryReleaseDownloadStrategy
    sha256 "4c875e0b7d9d1e339817d3561f584d75bb2b6b5e0a983a611315a15a32dc11eb"
  end

  def install
    bin.install "couture"
  end
end