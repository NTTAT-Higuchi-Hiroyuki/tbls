FROM mcr.microsoft.com/devcontainers/go:1-1.24-bookworm

# 必要なツールをインストール
RUN apt-get update && apt-get install -y \
    postgresql-client \
    default-mysql-client \
    sqlite3 \
    curl \
    && rm -rf /var/lib/apt/lists/*

# Node.jsをインストール（NodeSource経由）
RUN curl -fsSL https://deb.nodesource.com/setup_lts.x | bash - \
    && apt-get install -y nodejs \
    && rm -rf /var/lib/apt/lists/*

# GitHub CLIをインストール
RUN curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg \
    && chmod go+r /usr/share/keyrings/githubcli-archive-keyring.gpg \
    && echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | tee /etc/apt/sources.list.d/github-cli.list > /dev/null \
    && apt-get update \
    && apt-get install -y gh \
    && rm -rf /var/lib/apt/lists/*

# Go開発ツールをインストール
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Claude Codeをインストール
RUN npm install -g @anthropic-ai/claude-code && npm install -g ccmanager

# 作業ディレクトリを設定
WORKDIR /workspace
