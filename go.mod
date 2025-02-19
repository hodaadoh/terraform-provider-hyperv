module github.com/taliesins/terraform-provider-hyperv

go 1.13

require (
	github.com/Azure/go-ntlmssp v0.0.0-20200615164410-66371956d46c // indirect
	github.com/ChrisTrenkamp/goxpath v0.0.0-20190607011252-c5096ec8773d // indirect
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/dylanmei/iso8601 v0.1.0
	github.com/fatih/color v1.10.0 // indirect
	github.com/gofrs/uuid v4.0.0+incompatible // indirect
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/google/go-cmp v0.5.5 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.0 // indirect
	github.com/hashicorp/go-uuid v1.0.2 // indirect
	github.com/hashicorp/hcl/v2 v2.9.0 // indirect
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.4.4
	github.com/hashicorp/terraform-plugin-test/v2 v2.1.2 // indirect
	github.com/hashicorp/yamux v0.0.0-20200609203250-aecfd211c9ce // indirect
	github.com/jolestar/go-commons-pool/v2 v2.1.1
	github.com/masterzen/simplexml v0.0.0-20190410153822-31eea3082786 // indirect
	github.com/masterzen/winrm v0.0.0-20201030141608-56ca5c5f2380
	github.com/mitchellh/copystructure v1.1.1 // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/oklog/run v1.1.0 // indirect
	github.com/segmentio/ksuid v1.0.3
	github.com/stretchr/testify v1.6.1 // indirect
	github.com/vmihailenco/msgpack v4.0.4+incompatible // indirect
	golang.org/x/crypto v0.0.0-20210220033148-5ea612d1eb83 // indirect
	golang.org/x/net v0.0.0-20210226172049-e18ecbb05110 // indirect
	golang.org/x/sys v0.0.0-20210304195927-444254391f19 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20210303154014-9728d6b83eeb // indirect
	google.golang.org/grpc v1.36.0 // indirect
)

// BUAK: Replace reference of upstream provider to our custom code
replace github.com/taliesins/terraform-provider-hyperv => github.com/hodaadoh/terraform-provider-hyperv v1.0.3
