# Twintail

Tailscale Serve Status Dashboard - Tailscaleでホストされているサービスのステータスを監視するWebアプリケーション

## 技術スタック

- Go 1.25
- Echo v5 (Webフレームワーク)
- Tailwind CSS
- DaisyUI (darkテーマ)

## ビルド手順

### 初回セットアップ

Node.jsの依存関係とair（ホットリロードツール）をインストールします：

```bash
make install-deps
```

### 開発サーバー

Go + Tailwind CSSのホットリロード開発サーバーを起動します：

```bash
make dev
```

- Goファイル・HTMLファイルの変更時に自動で再ビルド＆再起動
- Tailwind CSSの変更時に自動で再ビルド
- ブラウザで `http://localhost:8077` にアクセス

### 本番ビルド

CSSビルドとGoバイナリのビルドを一括で実行します：

```bash
make build
```

本番ビルドでは静的ファイルがバイナリに埋め込まれます。

## 実行方法

```bash
./twintail
```

ブラウザで `http://localhost:8077` にアクセスしてください。

## プロジェクト構造

```
twintail/
├── assets/
│   └── css/
│       └── input.css       # Tailwind CSSのエントリーポイント
├── static/
│   └── css/
│       └── output.css      # ビルドされたCSS（バイナリに埋め込み）
├── views/
│   └── index.html          # メインページテンプレート
├── controllers/            # コントローラ層
├── services/               # サービス層
├── server.go               # サーバーエントリーポイント
├── tailwind.config.js      # Tailwind CSS設定
├── package.json            # Node.js依存関係
├── Makefile                # ビルドスクリプト
└── go.mod/go.sum           # Goモジュール
```

## 開発メモ

- 本番ビルド時のみ `//go:embed static` により静的ファイルがバイナリに埋め込まれます
- 開発時は `os.DirFS` を使用し、ファイル変更が即座に反映されます
- ビルドされたバイナリは単体で実行可能です（staticディレクトリ不要）
- DaisyUIのdarkテーマが適用されています
