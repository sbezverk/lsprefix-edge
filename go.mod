module github.com/sbezverk/lsprefix-edge

go 1.15

require (
	github.com/Shopify/sarama v1.27.2
	github.com/arangodb/go-driver v0.0.0-20201202080739-c41c94f2de00
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/sbezverk/gobmp v0.0.1-beta.0.20210125022733-d2f3a99d56ad
	github.com/sbezverk/topology v0.0.0-20201218172603-92c874d9602a
)

replace (
	github.com/sbezverk/gobmp => ../gobmp
	github.com/sbezverk/topology => ../topology
)
