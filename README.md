# Twintail

Tailscale Services Dashboard - [Tailscale Services](https://tailscale.com/kb/1552/tailscale-services)を管理するWebアプリケーション

> [!WARNING]
> このプロジェクトはVibe Codingにより生成されたものです。テストコードを含むすべてのソースコードはAIによって生成されており、十分なレビューが行われていない可能性があります。本番環境での使用には注意してください。

## 技術スタック

- Go 1.25
- Echo v5 (Webフレームワーク)
- Vite
- Tailwind CSS v4
- DaisyUI v5

## ビルド手順

### 初回セットアップ

Node.jsの依存関係とair（ホットリロードツール）をインストールします：

```bash
make install-deps
```

### 開発サーバー

Go + Viteのホットリロード開発サーバーを起動します：

```bash
make dev
```

- Goファイル・HTMLファイルの変更時に自動で再ビルド＆再起動
- CSS/TSの変更時に自動で再ビルド
- ブラウザで `http://localhost:8077` にアクセス

### 本番ビルド

フロントエンドビルドとGoバイナリのビルドを一括で実行します：

```bash
make build
```

本番ビルドでは静的ファイルがバイナリに埋め込まれます。

## 実行方法

```bash
./twintail-linux-arm64  # ARM64の場合
./twintail-linux-amd64  # AMD64の場合
```

ブラウザで `http://localhost:8077` にアクセスしてください。

## プロジェクト構造

```
twintail/
├── assets/
│   └── css/
│       └── input.css         # Tailwind CSSのエントリーポイント
├── static/
│   └── dist/                 # Viteビルド出力（バイナリに埋め込み）
├── views/
│   ├── layouts/
│   │   └── base.html         # ベースレイアウト
│   ├── index.html            # サービス一覧
│   ├── show_service.html     # サービス詳細
│   ├── new_service.html      # サービス作成
│   ├── confirm_delete.html   # サービス削除確認
│   ├── new_endpoint.html     # エンドポイント作成
│   └── confirm_delete_endpoint.html  # エンドポイント削除確認
├── controllers/              # コントローラ層
├── services/                 # サービス層（Tailscale CLI連携）
├── server.go                 # サーバーエントリーポイント
├── vite.config.ts            # Vite設定
├── package.json              # Node.js依存関係
├── Makefile                  # ビルドスクリプト
└── go.mod/go.sum             # Goモジュール
```

## 開発メモ

- 本番ビルド時のみ `//go:embed static` により静的ファイルがバイナリに埋め込まれます
- 開発時は `os.DirFS` を使用し、ファイル変更が即座に反映されます
- ビルドされたバイナリは単体で実行可能です（staticディレクトリ不要）
