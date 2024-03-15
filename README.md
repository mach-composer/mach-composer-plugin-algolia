# Algolia Plugin for Mach Composer 

This repository contains the Algolia plugin for Mach Composer. It requires Mach Composer 3.x

## Usage

```yaml
mach_composer:
  version: 1
  plugins:
    algolia:
      source: mach-composer/algolia
      version: 0.1.0

global:
  # ...

sites:
  - identifier: my-site

    algolia:
      api_key: api-key
      app_id: app-id

```
