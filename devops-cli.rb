class DevopsCli < Formula
  desc "A DevOps automation CLI tool"
  homepage "https://github.com/yourusername/devops-cli"
  version "1.0.0"
  
  if OS.mac?
    url "https://github.com/yourusername/devops-cli/releases/download/v1.0.0/devops-cli-darwin-amd64.tar.gz"
    sha256 "your-sha256-hash"
  elsif OS.linux?
    url "https://github.com/yourusername/devops-cli/releases/download/v1.0.0/devops-cli-linux-amd64.tar.gz"
    sha256 "your-sha256-hash"
  end

  def install
    bin.install "devops-cli"
  end

  test do
    system "#{bin}/devops-cli", "--version"
  end
end 