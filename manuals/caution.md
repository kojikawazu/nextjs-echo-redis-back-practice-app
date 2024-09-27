# 注意点

ここに実装、構築中に発生した内容を書いていく。

## App RunnerからElastiCacheへのアクセス

- AWS App RunnerからElastiCacheへアクセスするには、VPC Connectorが必要。
- ただし、VPC Connectorを使用すると、App Runnerがパブリックインターネットに直接アクセスできなくなり、実質的にプライベートサブネットとなる。

## Supabaseのアクセス

**現象**

- AWS App RunnerにバックエンドAPIをデプロイした後、Supabaseに接続できなくなった。

**原因**  

- App RunnerからElastiCacheに接続するためにVPC Connectorを追加したが、この変更によってApp Runnerがパブリックインターネット経由でSupabaseにアクセスできなくなった。
  - VPC Connectorを追加すると、App Runnerのトラフィックは指定したVPCを経由するため、パブリックなインターネット接続が無効化されてしまう。

**解消方法**

- NAT GatewayをVPCに追加して、VPC内のリソース（App Runner）がインターネットにアクセスできるようにした。
  - NAT Gateway経由でSupabaseのパブリックエンドポイントにアクセス可能となり、問題が解消された。
  - NAT Gatewayは、プライベートサブネットにあるリソースがインターネットアクセスできるようにするためのサービスであり、パブリックなSupabaseへの接続が可能となる。

**今後の注意点**

- VPC Connectorを使用する場合、外部パブリックサービス（例: Supabase）へのアクセスを確保するために、NAT Gatewayや他の通信経路を適切に構築する必要がある。