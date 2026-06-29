# Changelog

## [4.9.0](https://github.com/rudderlabs/terraform-provider-rudderstack/compare/v4.8.0...v4.9.0) (2026-06-29)


### Features

* require non-empty consents in consent management schema [SDK-4965] ([#276](https://github.com/rudderlabs/terraform-provider-rudderstack/issues/276)) ([3b37744](https://github.com/rudderlabs/terraform-provider-rudderstack/commit/3b37744009f759d23f7c9984152018f464509d50))

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
