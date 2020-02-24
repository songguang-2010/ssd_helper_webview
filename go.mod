module ssd_helper_webview

go 1.12

// github.com/karalabe/xgo v0.0.0-20191115072854-c5ccff8648a7 // indirect
require (
	controller v0.0.0
	github.com/spf13/viper v1.6.1 // indirect
	github.com/zserge/webview v0.0.0-20191103184548-1a9ebffc2601 // indirect
	golang.org/x/tools/gopls v0.3.2 // indirect
	// github.com/zserge/webview v0.0.0-20191103184548-1a9ebffc2601
	lib/config v0.0.0
	lib/file v0.0.0
	lib/logwrap v0.0.0
	lib/response v0.0.0
	lib/route v0.0.0
	lib/serror v0.0.0
	middleware v0.0.0
	model/aos v0.0.0
	model/misc v0.0.0
	model/order v0.0.0
	model/sku v0.0.0
	model/stat v0.0.0
	model/tps v0.0.0
	register v0.0.0
	src.techknowlogick.com/xgo v0.0.0-20191206145604-980bc3ce3f09 // indirect
)

replace (
	cloud.google.com/go => github.com/googleapis/google-cloud-go v0.49.0
	cloud.google.com/go/bigquery => github.com/googleapis/google-cloud-go/bigquery v1.3.0
	cloud.google.com/go/datastore => github.com/googleapis/google-cloud-go/datastore v1.0.0
	cloud.google.com/go/pubsub => github.com/googleapis/google-cloud-go/pubsub v1.1.0
	cloud.google.com/go/storage => github.com/googleapis/google-cloud-go/storage v1.4.0
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20191202143827-86a70503ff7e
	golang.org/x/exp => github.com/golang/exp v0.0.0-20191129062945-2f5052295587
	golang.org/x/image => github.com/golang/image v0.0.0-20191009234506-e7c1f5e7dbb8
	golang.org/x/lint => github.com/golang/lint v0.0.0-20191125180803-fdd1cda4f05f
	golang.org/x/mobile => github.com/golang/mobile v0.0.0-20191130191448-5c0e7e404af8
	golang.org/x/mod => github.com/golang/mod v0.1.0
	golang.org/x/net => github.com/golang/net v0.0.0-20191204025024-5ee1b9f4859a
	golang.org/x/oauth2 => github.com/golang/oauth2 v0.0.0-20191202225959-858c2ad4c8b6
	golang.org/x/sync => github.com/golang/sync v0.0.0-20190911185100-cd5d95a43a6e
	golang.org/x/sys => github.com/golang/sys v0.0.0-20191204072324-ce4227a45e2e
	//  golang.org/x/sys => github.com/golang/sys latest
	golang.org/x/text => github.com/golang/text v0.3.2
	golang.org/x/time => github.com/golang/time v0.0.0-20191024005414-555d28b269f0
	golang.org/x/tools => github.com/golang/tools v0.0.0-20191204011308-9611592c72f6
	google.golang.org/api => github.com/googleapis/google-api-go-client v0.14.0
	google.golang.org/appengine => github.com/golang/appengine v1.6.5
	google.golang.org/genproto => github.com/googleapis/go-genproto v0.0.0-20191203220235-3fa9dbf08042
	google.golang.org/grpc => github.com/grpc/grpc-go v1.25.1
)

replace register v0.0.0 => ./register

replace controller v0.0.0 => ./controller

replace middleware v0.0.0 => ./middleware

replace lib/response v0.0.0 => ./lib/response

replace lib/route v0.0.0 => ./lib/route

replace lib/config v0.0.0 => ./lib/config

replace lib/rabbitmq v0.0.0 => ./lib/rabbitmq

replace model/sku v0.0.0 => ./model/sku

replace model/order v0.0.0 => ./model/order

replace model/aos v0.0.0 => ./model/aos

replace model/tps v0.0.0 => ./model/tps

replace model/misc v0.0.0 => ./model/misc

replace model/stat v0.0.0 => ./model/stat

replace lib/log v0.0.0 => ./lib/log

replace lib/logwrap v0.0.0 => ./lib/logwrap

replace lib/file v0.0.0 => ./lib/file

replace lib/serror v0.0.0 => ./lib/serror

replace lib/model v0.0.0 => ./lib/model

replace lib/coroutine v0.0.0 => ./lib/coroutine

replace lib/stime v0.0.0 => ./lib/stime
