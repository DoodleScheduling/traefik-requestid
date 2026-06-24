# Traefik Request ID header

[![release](https://img.shields.io/github/release/DoodleScheduling/traefik-requestid/all.svg)](https://github.com/DoodleScheduling/traefik-requestid/releases)
[![report](https://goreportcard.com/badge/github.com/DoodleScheduling/traefik-requestid)](https://goreportcard.com/report/github.com/DoodleScheduling/traefik-requestid)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/DoodleScheduling/traefik-requestid/badge)](https://api.securityscorecards.dev/projects/github.com/DoodleScheduling/traefik-requestid)
[![Coverage Status](https://coveralls.io/repos/github/DoodleScheduling/traefik-requestid/badge.svg?branch=master)](https://coveralls.io/github/DoodleScheduling/traefik-requestid?branch=master)
[![license](https://img.shields.io/github/license/DoodleScheduling/traefik-requestid.svg)](https://github.com/DoodleScheduling/traefik-requestid/blob/master/LICENSE)

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
