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
                            OldHeader: 'The-Old-Header'
                            NewHeader: 'The-New-Header'
```
