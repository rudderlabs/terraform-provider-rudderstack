# Changelog

## [4.13.0](https://github.com/rudderlabs/terraform-provider-rudderstack/compare/v4.12.0...v4.13.0) (2026-08-28)


### Features

* **destinations:** add offline conversions destinations ([#316](https://github.com/rudderlabs/terraform-provider-rudderstack/issues/316)) ([d520068](https://github.com/rudderlabs/terraform-provider-rudderstack/commit/d5200684d1e8011bbb2a8a3b515a2f963b002ca0))
* expose customerio connection mode config ([#317](https://github.com/rudderlabs/terraform-provider-rudderstack/issues/317)) ([655b6b9](https://github.com/rudderlabs/terraform-provider-rudderstack/commit/655b6b9d90468c13db477a099025d53a26551b20))

## [4.12.0](https://github.com/rudderlabs/terraform-provider-rudderstack/compare/v4.11.0...v4.12.0) (2026-08-24)


### Features

* **retl:** support Customer.io event object syncs ([#302](https://github.com/rudderlabs/terraform-provider-rudderstack/issues/302)) ([875841a](https://github.com/rudderlabs/terraform-provider-rudderstack/commit/875841add11efef6102d11f85ed3d370304b0faf))


### Bug Fixes

* handle backend secret redaction (destination read + e2e verify) ([#301](https://github.com/rudderlabs/terraform-provider-rudderstack/issues/301)) ([b9cd742](https://github.com/rudderlabs/terraform-provider-rudderstack/commit/b9cd7425e03378b6d9fdcb2baaa958f4cf7ceb54))
* remove HubSpot legacy API-key auth option ([#312](https://github.com/rudderlabs/terraform-provider-rudderstack/issues/312)) ([c893f62](https://github.com/rudderlabs/terraform-provider-rudderstack/commit/c893f62ac154b7ae2ce1f265ed90fd164dbebf99))


### Miscellaneous

* **deps:** bump google.golang.org/grpc ([35b97ea](https://github.com/rudderlabs/terraform-provider-rudderstack/commit/35b97ea28a10998aec8fc1c769efb4d7cb8a78d6))
* **deps:** bump google.golang.org/grpc from 1.80.0 to 1.82.1 in the go_modules group across 1 directory ([#224](https://github.com/rudderlabs/terraform-provider-rudderstack/issues/224)) ([35b97ea](https://github.com/rudderlabs/terraform-provider-rudderstack/commit/35b97ea28a10998aec8fc1c769efb4d7cb8a78d6))
* **deps:** bump the actions group across 1 directory with 5 updates ([#280](https://github.com/rudderlabs/terraform-provider-rudderstack/issues/280)) ([592b7a1](https://github.com/rudderlabs/terraform-provider-rudderstack/commit/592b7a171a302aa2d5816bcf7bf21ad520cb9393))
* upgrade rudder-iac to v0.22.0 ([#299](https://github.com/rudderlabs/terraform-provider-rudderstack/issues/299)) ([672fa72](https://github.com/rudderlabs/terraform-provider-rudderstack/commit/672fa72b9a092d5591eee6c107d282f8dbc164c1))

## [4.11.0](https://github.com/rudderlabs/terraform-provider-rudderstack/compare/v4.10.0...v4.11.0) (2026-07-30)


### Features

* add Amplitude Browser SDK v2 auto_capture fields ([6859a48](https://github.com/rudderlabs/terraform-provider-rudderstack/commit/6859a4871eb16d9df047d0032388f8248a5a4608))
* **amplitude:** add Browser SDK v2 auto_capture fields ([#298](https://github.com/rudderlabs/terraform-provider-rudderstack/issues/298)) ([6859a48](https://github.com/rudderlabs/terraform-provider-rudderstack/commit/6859a4871eb16d9df047d0032388f8248a5a4608))
* **amplitude:** default sdk_version to browser sdk v2 ([#294](https://github.com/rudderlabs/terraform-provider-rudderstack/issues/294)) ([5024fa0](https://github.com/rudderlabs/terraform-provider-rudderstack/commit/5024fa0438952838a0c2dc931b2c8320ff3ab9bb))
* destination version scaffolding for ConfigMeta and generatetf (INT-6494) ([#289](https://github.com/rudderlabs/terraform-provider-rudderstack/issues/289)) ([6f0a77f](https://github.com/rudderlabs/terraform-provider-rudderstack/commit/6f0a77fdc232f5ccc7f5c267c1f6c72b77c2e16b))

## [4.10.0](https://github.com/rudderlabs/terraform-provider-rudderstack/compare/v4.9.0...v4.10.0) (2026-07-10)


### Features

* add missing source type support in destinations ([#292](https://github.com/rudderlabs/terraform-provider-rudderstack/issues/292)) ([9fd66af](https://github.com/rudderlabs/terraform-provider-rudderstack/commit/9fd66af65ec790d2305319ae993ee2907ed71626))

## [4.9.0](https://github.com/rudderlabs/terraform-provider-rudderstack/compare/v4.8.0...v4.9.0) (2026-07-07)


### Features

* add Confluent Cloud destination support ([#268](https://github.com/rudderlabs/terraform-provider-rudderstack/issues/268)) ([ff157c8](https://github.com/rudderlabs/terraform-provider-rudderstack/commit/ff157c853fe4d7ef3bf0049d8ddd0971dd00df6f))
* require non-empty consents in consent management schema [SDK-4965] ([#276](https://github.com/rudderlabs/terraform-provider-rudderstack/issues/276)) ([3b37744](https://github.com/rudderlabs/terraform-provider-rudderstack/commit/3b37744009f759d23f7c9984152018f464509d50))


### Miscellaneous

* **e2e:** BigQuery → Customer.io Audience smoke + per-scenario harness refactor ([#288](https://github.com/rudderlabs/terraform-provider-rudderstack/issues/288)) ([fd048e3](https://github.com/rudderlabs/terraform-provider-rudderstack/commit/fd048e36e0b7b42cc0b387ee778fdcd2cba84c40))

## [4.8.0](https://github.com/rudderlabs/terraform-provider-rudderstack/compare/v4.7.0...v4.8.0) (2026-06-25)


### Features

* **retl:** add rudderstack_retl_connection_customerio resource ([#275](https://github.com/rudderlabs/terraform-provider-rudderstack/issues/275)) ([a779a18](https://github.com/rudderlabs/terraform-provider-rudderstack/commit/a779a18ae74bb745869ace65bf30b68f9244304f))


### Bug Fixes

* **accounts:** BigQuery account uses options.project, not projectId ([c4834ca](https://github.com/rudderlabs/terraform-provider-rudderstack/commit/c4834ca6932613b3c4abed6db16ce49e7dd927cc))


### Miscellaneous

* retl-bq-account support ([#277](https://github.com/rudderlabs/terraform-provider-rudderstack/issues/277)) ([c4834ca](https://github.com/rudderlabs/terraform-provider-rudderstack/commit/c4834ca6932613b3c4abed6db16ce49e7dd927cc))

## [4.7.0](https://github.com/rudderlabs/terraform-provider-rudderstack/compare/v4.6.0...v4.7.0) (2026-06-17)


### Features

* align snowflake destination schema with config parity ([#267](https://github.com/rudderlabs/terraform-provider-rudderstack/issues/267)) ([1060519](https://github.com/rudderlabs/terraform-provider-rudderstack/commit/1060519cc5352767d6a6a339f3c829cd705f818d))

## [4.6.0](https://github.com/rudderlabs/terraform-provider-rudderstack/compare/v4.5.1...v4.6.0) (2026-06-12)


### Features

* add Amplitude sdk_version selector ([#257](https://github.com/rudderlabs/terraform-provider-rudderstack/issues/257)) ([558dec6](https://github.com/rudderlabs/terraform-provider-rudderstack/commit/558dec6ea701f9226ed82b83b1fce812229da922))
* validate duplicate consent providers in plan diff ([#254](https://github.com/rudderlabs/terraform-provider-rudderstack/issues/254)) ([ed5229a](https://github.com/rudderlabs/terraform-provider-rudderstack/commit/ed5229aec6bcd98630a2010b0e0f53acca70794a))


### Miscellaneous

* upgrade vulnerabilities ([#252](https://github.com/rudderlabs/terraform-provider-rudderstack/issues/252)) ([2253c38](https://github.com/rudderlabs/terraform-provider-rudderstack/commit/2253c38361204eb590c18a2d556a5816b1831d93))


### Documentation

* **e2e:** add operational guide for E2E test suite ([#250](https://github.com/rudderlabs/terraform-provider-rudderstack/issues/250)) ([f72d00d](https://github.com/rudderlabs/terraform-provider-rudderstack/commit/f72d00d70b6b7dbed73e88f1825b202f72b70392))
* refresh knowledge docs for repo-local skill ([#255](https://github.com/rudderlabs/terraform-provider-rudderstack/issues/255)) ([3c4a52e](https://github.com/rudderlabs/terraform-provider-rudderstack/commit/3c4a52ec8619ce799f8a5ef19e77c1d2f133abd1))

## [4.5.1](https://github.com/rudderlabs/terraform-provider-rudderstack/compare/v4.5.0...v4.5.1) (2026-05-19)


### Miscellaneous

* **release:** adopt release-please for automated releases ([#247](https://github.com/rudderlabs/terraform-provider-rudderstack/issues/247)) ([0f8812d](https://github.com/rudderlabs/terraform-provider-rudderstack/commit/0f8812d638bf074d96d5c1d78c8d1e83ff4ac525))
