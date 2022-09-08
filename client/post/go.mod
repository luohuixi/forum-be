module post-client

replace forum => ../../

replace forum-post => ../../microservice/post

go 1.18

require (
	forum v0.0.0-00010101000000-000000000000
	forum-post v0.0.0-00010101000000-000000000000
	github.com/opentracing/opentracing-go v1.2.0
)
