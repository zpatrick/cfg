# Example Setup for a Go Application


## Setup:
cmd/
    |- <app>/
        |- main.go
internal/
    |- config/
        |- config.go
    |- db/
    |- svr/

## Rules

### Rule 1: Only the main package imports the config package.
- Packages should use configuration inversion: your package should be telling others exactly what it needs: how to be configured (configuration inversion), and which dependenc(ies) it needs to be executed (dependency inversion). 
- It can be easy to create a circular dependency when you try and pass config package objects around (e.g. a config.ServerConfig). That's why in this pattern we only import/create a config.Config instance in main and use that to construct the rest of the application.
- This decouples your packages from the config package. This makes it easy to swap in different means-of-config for running your application in scenarios where we don't necessarily need/want all of the baggage that comes with the config package, e.g. in unit tests.

### Rule 2: Ensure each configuration value is canonically defined.
- Load once on startup, fail fast if configuration is invalid. 
- Single 'setting' in application, 1-many 'providers' for a given setting.
- Order of providers is consistent.
- Keep dev/prod parity for main-path executions. 