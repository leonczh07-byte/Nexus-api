# Taida Tech Modifications

This repository is a fork of [QuantumNous/new-api](https://github.com/QuantumNous/new-api),
based on tag **v1.0.0-rc.21** (commit `bde9b2f4`). It is published in compliance with the
**GNU Affero General Public License v3 (AGPLv3)**, Section 13: the modified source below is the
complete corresponding source of the gateway that powers the production service operated by
Taida Tech (xw taida tech / taida tech API relay platform).

The fork remains under **AGPLv3**. See [LICENSE](LICENSE). Original copyright and notices are
preserved.

## Modifications

All changes are confined to Aliyun Bailian (DashScope) model adaptation, request routing,
and accounting. No changes to core billing logic, pricing, or authentication semantics.

### New files

| File | Purpose |
|------|---------|
| `relay/channel/ali/embedding_vision.go` | Aliyun multimodal embedding adaptor for `tongyi-embedding-vision-flash/plus` |
| `relay/channel/ali/tts.go` | Aliyun TTS adaptor for `qwen3-tts-flash` speech synthesis |
| `service/auto_route.go` | Task-type routing across upstream models for the `auto` virtual model |

### Modified files

| File | Change summary |
|------|----------------|
| `relay/channel/ali/adaptor.go` | Register embedding-vision and TTS capabilities; adapt requests and responses for new Bailian models |
| `relay/channel/task/ali/adaptor.go` | Adapt asynchronous video tasks for HappyHorse and Wan model endpoints, polling, and result parsing |
| `relay/channel/task/ali/constants.go` | Add constants for the new task models |
| `relay/channel/task/ali/adaptor_test.go` | Add tests for the task-model adaptations |
| `controller/relay.go` | Add the `auto` virtual-model routing hook before channel selection |
| `service/billing_session.go` | Add usage-accounting fields for orchestrated multi-model sessions |
| `model/token.go` | Handle token model allowlists for the `auto` virtual model |
| `model/user.go` | Adjust the user query used by the management console |
| `model/errors.go` | Map additional upstream error codes |
| `model/user_update_test.go` | Add tests for the user-query adjustment |

## Build

Built with the stock `Dockerfile` in this repository (multi-stage Go + Node build),
tagged internally as `meimo/new-api:v1.0.0-rc.21-nmd-models.*`. The build scripts are unchanged.

## Contact

Operated by Taida Tech. For compliance questions regarding this source release,
please open an issue in this repository.
