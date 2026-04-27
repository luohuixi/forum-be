module forum-agent

replace forum => ../../

replace forum-user => ../user

replace forum-chat => ../chat

replace forum-feed => ../feed

replace forum-ocr => ../ocr

replace forum-post => ../post

go 1.24.7

require (
	forum v0.0.0-00010101000000-000000000000
	github.com/cloudwego/eino v0.7.13
	github.com/cloudwego/eino-ext/components/embedding/dashscope v0.0.0-20260416081055-0ebab92e14f2
	github.com/cloudwego/eino-ext/components/indexer/es8 v0.0.0-20260422030438-09c530946953
	github.com/cloudwego/eino-ext/components/model/openai v0.1.13
	github.com/cloudwego/eino-ext/components/retriever/es8 v0.0.0-20260422030438-09c530946953
	github.com/cloudwego/eino-ext/components/tool/mcp v0.0.8
	github.com/elastic/go-elasticsearch/v8 v8.16.0
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/mark3labs/mcp-go v0.48.0
	github.com/samber/lo v1.53.0
	github.com/segmentio/kafka-go v0.4.38
	github.com/spf13/viper v1.12.0
	go.uber.org/zap v1.21.0
	gorm.io/gorm v1.25.12
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/BurntSushi/toml v1.6.0 // indirect
	github.com/bahlo/generic-list-go v0.2.0 // indirect
	github.com/bmatcuk/doublestar/v4 v4.6.1 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/bytedance/gopkg v0.1.3 // indirect
	github.com/bytedance/sonic v1.14.1 // indirect
	github.com/bytedance/sonic/loader v0.3.0 // indirect
	github.com/casbin/casbin/v2 v2.103.0 // indirect
	github.com/casbin/gorm-adapter/v3 v3.7.2 // indirect
	github.com/casbin/govaluate v1.3.0 // indirect
	github.com/casbin/redis-watcher/v2 v2.2.0 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/cloudwego/base64x v0.1.6 // indirect
	github.com/cloudwego/eino-ext/libs/acl/openai v0.1.17 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/eino-contrib/jsonschema v1.0.3 // indirect
	github.com/elastic/elastic-transport-go/v8 v8.9.0 // indirect
	github.com/evanphx/json-patch v0.5.2 // indirect
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/glebarez/go-sqlite v1.22.0 // indirect
	github.com/glebarez/sqlite v1.11.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-redis/redis v6.15.9+incompatible // indirect
	github.com/go-redis/redis/v8 v8.11.5 // indirect
	github.com/go-sql-driver/mysql v1.9.2 // indirect
	github.com/golang-sql/civil v0.0.0-20220223132316-b832511892a9 // indirect
	github.com/golang-sql/sqlexp v0.1.0 // indirect
	github.com/google/jsonschema-go v0.4.2 // indirect
	github.com/goph/emperror v0.17.2 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.13.0 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.1 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/pgtype v1.12.0 // indirect
	github.com/jackc/pgx/v4 v4.17.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/klauspost/cpuid/v2 v2.2.10 // indirect
	github.com/magiconair/properties v1.8.6 // indirect
	github.com/mailru/easyjson v0.9.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/meguminnnnnnnnn/go-openai v0.1.2 // indirect
	github.com/microsoft/go-mssqldb v1.8.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/ncruces/go-strftime v0.1.9 // indirect
	github.com/nikolalohinski/gonja v1.5.3 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/slongfield/pyfmt v0.0.0-20220222012616-ea85ff4c361f // indirect
	github.com/spf13/afero v1.8.2 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/subosito/gotenv v1.3.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/wk8/go-ordered-map/v2 v2.1.8 // indirect
	github.com/yargevad/filepathx v1.0.0 // indirect
	github.com/yosida95/uritemplate/v3 v3.0.2 // indirect
	go.opentelemetry.io/otel v1.29.0 // indirect
	go.opentelemetry.io/otel/metric v1.29.0 // indirect
	go.opentelemetry.io/otel/trace v1.29.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/arch v0.15.0 // indirect
	golang.org/x/crypto v0.39.0 // indirect
	golang.org/x/exp v0.0.0-20250408133849-7e4ce0ab07d0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	gopkg.in/ini.v1 v1.66.4 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/driver/mysql v1.5.7 // indirect
	gorm.io/driver/postgres v1.4.5 // indirect
	gorm.io/driver/sqlserver v1.5.4 // indirect
	gorm.io/plugin/dbresolver v1.5.3 // indirect
	modernc.org/libc v1.62.1 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.9.1 // indirect
	modernc.org/sqlite v1.37.0 // indirect
)
