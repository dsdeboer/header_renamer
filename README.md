Treafik Plugin
--------------

Working sample:

```yaml
# Template for configuration

http:
    middlewares:
        traefik_header_rename:
            plugin:
                traefik_header_rename:
                    Rules:
                        -   Rule:
                            Name: 'Header transformation'
                            Header: 'The-Old-Header'
                            Value: 'The-New-Header'
```
