module feed-client

replace forum => ../../

replace forum-feed => ../../microservice/feed

replace forum-user => ../../microservice/user

go 1.18

require (
	forum v0.0.0-00010101000000-000000000000
	forum-feed v0.0.0-00010101000000-000000000000
	github.com/opentracing/opentracing-go v1.2.0
)
