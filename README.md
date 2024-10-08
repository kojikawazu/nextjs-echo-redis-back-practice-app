# Next.js(TypeScript) + Echo(Go) + Redisを使ったキャッシュサーバー学習用(バックエンド用)

バックエンドAPIにキャッシュサーバーを導入したい。
キャッシュサーバーに慣れるために簡単なバックエンドAPIを構築する。

## Tech

- Frontend
  - Next.js
  - TypeScript
  - TailWindCSS
- BackEnd
  - Echo
  - Go
- Cloud
  - AWS
    - VPC
    - App Runner
    - ElastiCache
    - ECR
- Other
  - Redis
  - Supabase

## Architecture

![アーキテクチャ図](./drawio/elasticache.drawio.png)

## 注意事項

詳細は [こちらのCautionファイル](./manuals/caution.md) を参照してください。

## Todo

- AWS ElasticCacheを最終的に使う。
  - このキャッシュサーバーを使えるようにする。
- AWS App Runnerに構築する。
  - ALB + ECSが理想。しかし、まずはキャッシュサーバーに慣れるためにApp Runnerを使用する。
- バックエンド中心
  - フロントエンドは使わないかもしれない。(キャッシュサーバーに慣れるため)