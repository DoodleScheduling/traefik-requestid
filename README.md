# Traefik Request ID header

[![release](https://img.shields.io/github/release/DoodleScheduling/traefik-envoy-ratelimit/all.svg)](https://github.com/DoodleScheduling/traefik-envoy-ratelimit/releases)
[![report](https://goreportcard.com/badge/github.com/DoodleScheduling/traefik-envoy-ratelimit)](https://goreportcard.com/report/github.com/DoodleScheduling/traefik-envoy-ratelimit)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/DoodleScheduling/traefik-envoy-ratelimit/badge)](https://api.securityscorecards.dev/projects/github.com/DoodleScheduling/traefik-envoy-ratelimit)
[![Coverage Status](https://coveralls.io/repos/github/DoodleScheduling/traefik-envoy-ratelimit/badge.svg?branch=master)](https://coveralls.io/github/DoodleScheduling/traefik-envoy-ratelimit?branch=master)
[![license](https://img.shields.io/github/license/DoodleScheduling/traefik-envoy-ratelimit.svg)](https://github.com/DoodleScheduling/traefik-envoy-ratelimit/blob/master/LICENSE)

Create a unique identifier and attach a header both up and downstream.
By default the header is called `X-Request-Id`.

## Configuration

```yaml
apiVersion: traefik.io/v1alpha1
kind: Middleware
metadata:
  name: requestid
spec:
  plugin:
    requestid: {}
```
