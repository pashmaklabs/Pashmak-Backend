# Changelog

All notable changes to this project will be documented in this file. The format is based on [Keep a Changelog](https://keepachangelog.com/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.2.0] - 2025-08-01

### Added

- Integrated system-wide telemetry monitoring utilizing Prometheus and Grafana dashboards (#39).
- Added third-party Google Reviews injection endpoints into the internal Comment API (#96).
- Implemented image URL extraction workers targeting raw `greviews` information tables.
- Integrated `pgvector` engine expansion packages and setup localized vector database structures (#38).
- Implemented multidimensional vector similarity searches supporting semantic spatial discovery (#96).
- Implemented algorithmic query text category classification and query filtering.
- Created standalone feature routes for adding custom tracking coordinates and mapping locations (#36).

### Fixed

- Fixed cascading referential integrity delete failures across saved locations on user label deletions (#97).
- Fixed product optimization faults within the pashmak analysis package module (#97).
- Resolved validation binding issues targeting text embedding processing pipelines (#96).
- Fixed localized search point returns when fetching data records via `getPlaceByID` (#96).

---

## [1.1.0] - 2025-07-21

### Added

- Completed core features and added entity naming properties to the Saved Locations API wrapper (#98).
- Integrated OpenAI client dependencies supporting context-aware mapping algorithms (#61).
- Implemented custom user Roles and Access Permission control structures (#28).
- Developed specialized moderation systems for comment flag reporting (#27) and an admin reporting panel (#33).
- Added persistent user search history logging and deletion capabilities (#29).
- Created a `isReactedByCurrentUser` flag layer mapping interactive evaluation feedback (#35).
- Expanded spatial resources with dedicated multi-image location array uploads (#30).

### Fixed

- Patched database relation criteria handling inside the core Saved Location API pipelines (#98).
- Resolved login validation handler crashes and cookie parsing abnormalities.
- Repaired malformed identifier parameters bounding context interactions.

---

## [1.0.0] - 2025-05-15

### Added

- Initial project generation leveraging structured routing and clean config management env setups.
- Complete user Authentication system featuring secure multi-factor OTP over email endpoints (#8).
- Implemented key storage layers and caching frameworks utilizing PostgreSQL and Redis servers.
- Created core user profile dashboards (#13), system navigation algorithms (#17), and comments architecture (#15).
- Integrated PostGIS spatial processing tooling alongside osm2pgsql pipeline configurations (#48).
- Set up object storage file uploads via MinIO connection APIs (#20) alongside WebP media resolution compressors.
- Integrated standard generic pagination modules formatting response arrays smoothly (#23).

### Fixed

- Corrected CORS response headers mapping AllowedOrigins matching secure setups.
- Patched local deployment compose scripts optimizing MinIO storage dependencies.

---

## [0.2.0] - 2025-04-10

### Added

- Created foundational REST endpoints for user Profile dashboards (#13), platform Navigation pathways (#17), and base Comments layout architectures (#15).
- Added multi-environment containerization layers via a central Dockerfile configuration.
- Set up user registration and account creation flows via a dedicated Signup API handler (#12).
- Implemented automated credential recovery workflows via a Forget Password module (#12).

### Fixed

- Standardized CORS parameter scopes via single-url mappings inside the AllowdOrigins module.
- Patched authentication domain tracking cookie scopes alongside HTTP return response codes.

---

## [0.1.0] - 2025-03-25

### Added

- Initialized core repository file systems, dependency packages, structure layouts, and environment variable logic blocks (#3).
- Integrated localized infrastructure endpoints connecting to PostgreSQL and Redis database instances (#8).
- Developed backend multi-factor Authentication structures processing custom OTP codes over automated Email delivery pipelines (#8).

### Fixed

- Corrected localized telemetry formatting structures within core initialization modules (#3).
- Resolved operational database connection exceptions thrown during system setup handshakes (#8).
