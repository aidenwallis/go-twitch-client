# go-twitch-client

[![Go Reference](https://pkg.go.dev/badge/github.com/aidenwallis/go-twitch-client.svg)](https://pkg.go.dev/github.com/aidenwallis/go-twitch-client) [![codecov](https://codecov.io/gh/aidenwallis/go-twitch-client/branch/main/graph/badge.svg?token=s6fH5g5GG0)](https://codecov.io/gh/aidenwallis/go-twitch-client)

A simple, low level wrapper for the [Twitch API](https://dev.twitch.tv). The package deliberately does not make considerations around rate-limiting, token management, or other abstracted behaviour. It aims to be a very simple, thin layer between you and the Twitch API.

For that reason, you should use this package behind any rate-limiting logic, or authorization logic you build into your apps.

Each API is split into it's own package, documentation relevant to Helix lives in the [helix](helix/README.md) directory.

This package is built using Generics, and thus requires Go 1.18 or later.
