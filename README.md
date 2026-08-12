# Taida Tech Nexus API

Taida Tech Nexus API is a unified AI API gateway and model management platform maintained by Taida Tech for the Thai market. It provides a consistent interface for text, image, video, audio, embedding, and other AI services.

> [!IMPORTANT]
> Taida Tech Nexus API is a branded and modified distribution of [New API](https://github.com/QuantumNous/new-api). It is not the original New API project. This distribution is maintained independently by Taida Tech and is based on New API tag `v1.0.0-rc.21`.

[![License: AGPLv3](https://img.shields.io/badge/License-AGPLv3-blue.svg)](LICENSE)
[![Upstream: New API](https://img.shields.io/badge/Upstream-New%20API-black.svg)](https://github.com/QuantumNous/new-api)

## Overview

Taida Tech Nexus API helps teams connect multiple AI providers through one platform. It includes:

- Unified API access for multiple AI providers and model types
- Text, image, video, audio, and embedding model support
- User, token, quota, billing, and channel management
- Model routing, failover, and provider configuration
- A web console for administrators and users
- Taida-specific provider adapters and model mappings
- English-first and Thai-localized product interfaces

For a detailed record of Taida-specific changes, see [TAIDA_MODIFICATIONS.md](TAIDA_MODIFICATIONS.md).

## Source and Attribution

This repository contains the corresponding source code for the Taida Tech Nexus API distribution.

The upstream project is:

- Project: [New API](https://github.com/QuantumNous/new-api)
- Upstream copyright: QuantumNous and New API contributors
- Upstream license: [GNU Affero General Public License v3.0](LICENSE)
- Taida base version: `v1.0.0-rc.21`

Required upstream attribution:

> Frontend design and development by New API contributors.

Taida Tech has modified the upstream project. The modification history and material differences are documented in [TAIDA_MODIFICATIONS.md](TAIDA_MODIFICATIONS.md). Taida Tech does not claim authorship of the original New API code.

## Commercial Use and License Compliance

Commercial use is permitted under the AGPLv3 only when all applicable license obligations are followed. In particular, distributors and operators of modified network services must:

1. Keep the complete [AGPLv3 license](LICENSE) with the covered work.
2. Preserve the upstream copyright and legal notices in [NOTICE](NOTICE).
3. Preserve the exact attribution notice shown above in an appropriate legal, About, footer, or attribution location.
4. Keep a visible link to the original [New API repository](https://github.com/QuantumNous/new-api).
5. Clearly state that this is a modified distribution and document material changes.
6. Offer users interacting with the software over a network access to the complete corresponding source code as required by AGPLv3 Section 13.
7. Preserve applicable third-party notices and licenses in [THIRD-PARTY-LICENSES.md](THIRD-PARTY-LICENSES.md).

This repository does not grant a separate proprietary or closed-source commercial license. Organizations that require terms outside the AGPLv3 obligations should contact the New API upstream maintainers at `support@quantumnous.com` about a separate commercial license.

## Deployment

Build the Taida branch from source:

```bash
git clone --branch taida https://github.com/leonczh07-byte/new-api.git
cd new-api
docker build -t taida-tech-nexus-api:local .
```

Deployment settings, provider credentials, database credentials, and private infrastructure configuration are environment-specific and are not included in this repository.

For upstream installation and configuration guidance, refer to the [New API documentation](https://docs.newapi.pro/).

## Security

- Never commit API keys, database passwords, email credentials, or other secrets.
- Rotate credentials immediately if they are exposed.
- Restrict administrative endpoints and follow least-privilege access practices.
- Use this software and connected AI services only in compliance with applicable laws, provider policies, and local regulations.

## Trademarks

“Taida Tech Nexus API” and related Taida branding identify Taida Tech's modified distribution. “New API” belongs to its respective owner and is used here only to identify the upstream open-source project.

## Legal Files

- [LICENSE](LICENSE)
- [NOTICE](NOTICE)
- [THIRD-PARTY-LICENSES.md](THIRD-PARTY-LICENSES.md)
- [TAIDA_MODIFICATIONS.md](TAIDA_MODIFICATIONS.md)

If this README conflicts with the license or legal notices, the terms in [LICENSE](LICENSE) and [NOTICE](NOTICE) control.
