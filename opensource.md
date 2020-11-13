---
layout: page
title: Open Source
permalink: /opensource/
---

### Go
* [crypto/x509: return errors instead of panicking](https://github.com/golang/go/commit/27280d8c14331c1c46cd90206be9f3c924f6b4c4)  
* [runtime: improve error messages after allocating a stack that is too big](https://github.com/golang/go/commit/331614c4daa5504ddfe35a96371cc34783d14cf1)  
* [src/go.mod, net/http: update bundled and latest golang.org/x/net](https://github.com/golang/go/commit/62fe10bf4e62c97af3bb8eb2ef72d9224a8752ba)  
* [net/http: reject HTTP/1.1 Content-Length with sign in response](https://github.com/golang/go/commit/8da78625b1fe2a6141d331f54248913936dc49c7)  
* [cmd/go: clean -cache -n should not delete cache](https://github.com/golang/go/commit/c0e8e405c02388fb8e7d3bea092f5aa8b19b2ad9)  
* [time: fix time.Before to reuse t.sec(), u.sec()](https://github.com/golang/go/commit/ff1eb428654b7815a8fc825f1cc29d6cf72cc2f7)  

### hashicorp
* [Expose Envoy's /stats for statsd agents (#7173)](https://github.com/hashicorp/consul/commit/a335aa57c54ffd19283815db23581765f93d588e)  
* [Add packer.ExpandUser() function to support tilde in usage of config.ValidationKeyPath (#8657)](https://github.com/hashicorp/packer/commit/beca6de71ba2c87c981cf12a997eb6008984c801)  

### beatlabs/patron
* [Sync goroutine assertion in component/http with caller test (#269)](https://github.com/beatlabs/patron/commit/ef6b531c6d33f4c4e53f86914c50af3986c58794)  
* [Retrieve error causing consumer group to close (#262)](https://github.com/beatlabs/patron/commit/6bdcd50ed2fd7289433a9a9e8625cfc4945840ee)  
* [Fix logger initialization in patron-cli (#258)](https://github.com/beatlabs/patron/commit/1bbe88cf64f86ebe16e02773daf103b916022649)  
* [Improve URL parameter handling for RawRoutes (#225)](https://github.com/beatlabs/patron/commit/d51c80b45417c5aa83439f5bc8a32f0e6f36ef28)  
* [Update tracing and metrics dependencies (#190)](https://github.com/beatlabs/patron/commit/776b5c6906bdb4660585b3d8d2c733993f547cf5)  
* [Add service builder (#159)](https://github.com/beatlabs/patron/commit/0b82771ef8d2e9ac3b32e0f195bc936cc7729902)  
* [Kafka operational metrics (#154)](https://github.com/beatlabs/patron/commit/6e1b30d4da762ce947b832def01769d266bea78d)  
* [Clean Up trace pkg (#147)](https://github.com/beatlabs/patron/commit/8c5dfab792d6cacc793b5da409528d121eb12746)  
* [Add generic cache package (#138)](https://github.com/beatlabs/patron/commit/aef37dadfb7dbd08a9071a15165c1d7dc9dd370d)  
* [Kafka producer builder (#131)](https://github.com/beatlabs/patron/commit/145aad510bbed6e69d78372f0aa8d7c37453eb9d)  
* [Inject encoder for kafka producer (#126)](https://github.com/beatlabs/patron/commit/9a1ce7e88753167a39042cc9200ec6d53ddd5fb0)  
* [HTTP Component Builder (#115)](https://github.com/beatlabs/patron/commit/26338ecddfc81993253c79f5ee2554aa1421c2eb)  
* [Parse DSN to populate Connection Information (#116)](https://github.com/beatlabs/patron/commit/39e1de9f59b69236a8a3c2597fd7bc6a9e8c1dd3)  

### mattermost
* [Mm 19027 (#12498)](https://github.com/mattermost/mattermost-server/commit/946a5c1417e659516363a866856a1ec10a624ff5)  
* [Migrate tests from model/compliance_test.go to use testify (#12497)](https://github.com/mattermost/mattermost-server/commit/23d495becaaeaf2677959986b7902754ac679c10)  
