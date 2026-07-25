# Taida Tech Modifications / 泰达科技修改说明

This repository is a fork of [QuantumNous/new-api](https://github.com/QuantumNous/new-api),
based on tag **v1.0.0-rc.21** (commit `bde9b2f4`). It is published in compliance with the
**GNU Affero General Public License v3 (AGPLv3)**, Section 13: the modified source below is the
complete corresponding source of the gateway that powers the production service operated by
Taida Tech (xw taida tech / taida tech API relay platform).

本仓库是 [QuantumNous/new-api](https://github.com/QuantumNous/new-api) 的 fork，
基于标签 **v1.0.0-rc.21**（commit `bde9b2f4`）。依据 **AGPLv3 第 13 条**发布：
以下修改即为我们生产环境（taida tech / xw taida tech API 中转平台）所运行网关的完整对应源码。

The fork remains under **AGPLv3**. See [LICENSE](LICENSE). Original copyright and notices are
preserved. 本 fork 继续采用 AGPLv3 协议，原始版权与声明全部保留。

## Modifications / 修改清单

All changes are confined to Aliyun Bailian (DashScope) model adaptation, request routing,
and accounting. No changes to core billing logic, pricing, or authentication semantics.

全部修改仅限于阿里百炼（DashScope）模型适配、请求路由与用量记录，
不涉及核心计费逻辑、定价或认证语义的改动。

### New files / 新增文件

| File | Purpose |
|------|---------|
| `relay/channel/ali/embedding_vision.go` | Aliyun multimodal embedding adaptor（`tongyi-embedding-vision-flash/plus`）|
| `relay/channel/ali/tts.go` | Aliyun TTS adaptor（`qwen3-tts-flash` 语音合成）|
| `service/auto_route.go` | `auto` virtual model: task-type based routing across upstream models（auto 虚拟模型智能路由）|

### Modified files / 修改文件

| File | Change summary |
|------|----------------|
| `relay/channel/ali/adaptor.go` | Register embedding-vision & TTS capabilities; request/response adaptation for new Bailian models（注册新模型能力）|
| `relay/channel/task/ali/adaptor.go` | Video/task async API adaptation: happyhorse & wan series model endpoints, task polling and result parsing（视频异步任务适配）|
| `relay/channel/task/ali/constants.go` | New task model constants（新增任务模型常量）|
| `relay/channel/task/ali/adaptor_test.go` | Tests for the above（对应测试）|
| `controller/relay.go` | Hook for `auto` virtual model routing before channel selection（auto 路由钩子）|
| `service/billing_session.go` | Usage accounting fields for orchestrated multi-model sessions（编排会话用量记录字段）|
| `model/token.go` | Token model-allowlist handling for the `auto` virtual model（auto 模型的令牌模型名单处理）|
| `model/user.go` | Minor user-query adjustment used by the console（配合管理台的查询微调）|
| `model/errors.go` | Additional error mapping for new upstream error codes（新增上游错误码映射）|
| `model/user_update_test.go` | Tests for the above（对应测试）|

## Build / 构建

Built with the stock `Dockerfile` in this repository (multi-stage Go + Node build),
tagged internally as `meimo/new-api:v1.0.0-rc.21-nmd-models.*`.
使用仓库自带的 `Dockerfile` 构建，未修改构建脚本。

## Contact / 联系

Operated by Taida Tech. For compliance questions regarding this source release,
please open an issue in this repository.
如有关于本次源码发布的合规问题，请在本仓库提交 issue。
