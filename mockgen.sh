# Repository
mockery --disable-version-string --case underscore --name=DBRepository --dir ./internal/repository --output ./mocks/repository
mockery --disable-version-string --case underscore --name=RedisRepository --dir ./internal/repository --output ./mocks/repository
mockery --disable-version-string --case underscore --name=MessagingRepository --dir ./internal/repository --output ./mocks/repository
mockery --disable-version-string --case underscore --name=NotificationRepository --dir ./internal/repository --output ./mocks/repository

# Service
mockery --disable-version-string --case underscore --name=Service --dir ./internal/service --output ./mocks/service