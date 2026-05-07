# Algolia Plugin for Mach Composer 

This repository contains the Algolia plugin for Mach Composer. It requires Mach Composer 3.x

## Usage

### Single application

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

### Multiple applications

Use `applications` when a single Mach Composer site needs to configure multiple
Algolia applications. The application `name` is linked to the component with the
same name and must be a valid Terraform provider alias.

> [!WARNING]
> Do not switch an existing site from single application configuration to
> multiple applications without a migration. The Terraform resource keys differ
> between the two modes, so switching directly would make Terraform delete the
> existing configured Algolia data and recreate it under the new keys.

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
      applications:
        - name: lab_digital
          api_key: lab-digital-api-key
          app_id: lab-digital-app-id
        - name: mach_composer
          api_key: mach-composer-api-key
          app_id: mach-composer-app-id
```

Components named `lab_digital` and `mach_composer` will receive the matching
aliased Algolia provider:

```hcl
providers = {
  algolia = algolia.lab_digital
}
```
