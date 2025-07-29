//gomodjail:confined
module github.com/containerd/nerdctl/v2

go 1.24.0

toolchain go1.24.5

require (
	github.com/GrigoryEvko/gozstd v1.22.1
	github.com/Masterminds/semver/v3 v3.4.0
	github.com/Microsoft/go-winio v0.6.2
	github.com/Microsoft/hcsshim v0.13.0
	github.com/compose-spec/compose-go/v2 v2.8.1 //gomodjail:unconfined
	github.com/containerd/accelerated-container-image v1.3.0
	github.com/containerd/cgroups/v3 v3.0.5 //gomodjail:unconfined
	github.com/containerd/console v1.0.5 //gomodjail:unconfined
	github.com/containerd/containerd/api v1.9.0
	github.com/containerd/containerd/v2 v2.1.3 //gomodjail:unconfined
	github.com/containerd/continuity v0.4.5 //gomodjail:unconfined
	github.com/containerd/errdefs v1.0.0
	github.com/containerd/fifo v1.1.0 //gomodjail:unconfined
	github.com/containerd/go-cni v1.1.12 //gomodjail:unconfined
	github.com/containerd/imgcrypt/v2 v2.0.1 //gomodjail:unconfined
	github.com/containerd/log v0.1.0
	github.com/containerd/nerdctl/mod/tigron v0.0.0
	github.com/containerd/nydus-snapshotter v0.15.2 //gomodjail:unconfined
	github.com/containerd/platforms v1.0.0-rc.1 //gomodjail:unconfined
	github.com/containerd/stargz-snapshotter v0.17.0 //gomodjail:unconfined
	github.com/containerd/stargz-snapshotter/estargz v0.17.0 //gomodjail:unconfined
	github.com/containerd/stargz-snapshotter/ipfs v0.17.0 //gomodjail:unconfined
	github.com/containerd/typeurl/v2 v2.2.3
	github.com/containernetworking/cni v1.3.0 //gomodjail:unconfined
	github.com/containernetworking/plugins v1.7.1 //gomodjail:unconfined
	github.com/coreos/go-iptables v0.8.0 //gomodjail:unconfined
	github.com/coreos/go-systemd/v22 v22.5.0
	github.com/cyphar/filepath-securejoin v0.4.1 //gomodjail:unconfined
	github.com/distribution/reference v0.6.0
	github.com/docker/cli v28.3.2+incompatible //gomodjail:unconfined
	github.com/docker/docker v28.3.2+incompatible //gomodjail:unconfined
	github.com/docker/go-connections v0.5.0
	github.com/docker/go-units v0.5.0
	github.com/fahedouch/go-logrotate v0.3.0 //gomodjail:unconfined
	github.com/fatih/color v1.18.0 //gomodjail:unconfined
	github.com/fluent/fluent-logger-golang v1.10.0
	github.com/fsnotify/fsnotify v1.9.0 //gomodjail:unconfined
	github.com/go-viper/mapstructure/v2 v2.4.0
	github.com/ipfs/go-cid v0.5.0
	github.com/klauspost/compress v1.18.0
	github.com/mattn/go-isatty v0.0.20 //gomodjail:unconfined
	github.com/moby/sys/mount v0.3.4
	github.com/moby/sys/signal v0.7.1
	github.com/moby/sys/user v0.4.0 //gomodjail:unconfined
	github.com/moby/sys/userns v0.1.0 //gomodjail:unconfined
	github.com/moby/term v0.5.2 //gomodjail:unconfined
	github.com/muesli/cancelreader v0.2.2 //gomodjail:unconfined
	github.com/opencontainers/go-digest v1.0.0
	github.com/opencontainers/image-spec v1.1.1
	github.com/opencontainers/runtime-spec v1.2.1
	github.com/pelletier/go-toml/v2 v2.2.4
	github.com/rootless-containers/bypass4netns v0.4.2 //gomodjail:unconfined
	github.com/rootless-containers/rootlesskit/v2 v2.3.5 //gomodjail:unconfined
	github.com/shirou/gopsutil/v4 v4.25.6
	github.com/spf13/cobra v1.9.1 //gomodjail:unconfined
	github.com/spf13/pflag v1.0.7 //gomodjail:unconfined
	github.com/vishvananda/netlink v1.3.1 //gomodjail:unconfined
	github.com/vishvananda/netns v0.0.5 //gomodjail:unconfined
	github.com/yuchanns/srslog v1.1.0
	go.uber.org/mock v0.5.2
	go.yaml.in/yaml/v3 v3.0.4
	golang.org/x/crypto v0.40.0
	golang.org/x/net v0.42.0
	golang.org/x/sync v0.16.0 //gomodjail:unconfined
	golang.org/x/sys v0.34.0 //gomodjail:unconfined
	golang.org/x/term v0.33.0 //gomodjail:unconfined
	golang.org/x/text v0.27.0
	gotest.tools/v3 v3.5.2
	tags.cncf.io/container-device-interface v1.0.1 //gomodjail:unconfined
)

require (
	github.com/Azure/go-ansiterm v0.0.0-20250102033503-faa5f7b0171c // indirect
	github.com/cilium/ebpf v0.19.0 // indirect
	github.com/containerd/errdefs/pkg v0.3.0 // indirect
	github.com/containerd/go-runc v1.1.0 // indirect
	github.com/containerd/plugin v1.0.0 // indirect
	github.com/containerd/ttrpc v1.2.7 // indirect
	github.com/containers/ocicrypt v1.2.1 // indirect
	github.com/creack/pty v1.1.24 // indirect
	github.com/djherbis/times v1.6.0 // indirect
	github.com/docker/docker-credential-helpers v0.9.3 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-jose/go-jose/v4 v4.1.1 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20241129210726-2c02b8208cf8 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/klauspost/cpuid/v2 v2.3.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-shellwords v1.0.12 // indirect
	github.com/miekg/pkcs11 v1.1.1 // indirect
	github.com/minio/sha256-simd v1.0.1 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/moby/docker-image-spec v1.3.1 // indirect
	github.com/moby/locker v1.0.1 // indirect
	github.com/moby/sys/mountinfo v0.7.2 // indirect
	github.com/moby/sys/sequential v0.6.0 // indirect
	github.com/moby/sys/symlink v0.3.0 // indirect
	github.com/mr-tron/base58 v1.2.0 // indirect
	github.com/multiformats/go-base32 v0.1.0 // indirect
	github.com/multiformats/go-base36 v0.2.0 // indirect
	github.com/multiformats/go-multiaddr v0.16.0 // indirect
	github.com/multiformats/go-multibase v0.2.0 // indirect
	github.com/multiformats/go-multihash v0.2.3 // indirect
	github.com/multiformats/go-varint v0.0.7 // indirect
	github.com/opencontainers/runtime-tools v0.9.1-0.20221107090550-2e043c6bd626 // indirect
	github.com/opencontainers/selinux v1.12.0 // indirect
	github.com/petermattis/goid v0.0.0-20250721140440-ea1c0173183e // indirect
	github.com/philhofer/fwd v1.2.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/santhosh-tekuri/jsonschema/v6 v6.0.2 // indirect
	github.com/sasha-s/go-deadlock v0.3.5 // indirect
	//gomodjail:unconfined
	github.com/sirupsen/logrus v1.9.3
	github.com/smallstep/pkcs7 v0.2.1 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/stefanberger/go-pkcs11uri v0.0.0-20230803200340-78284954bff6 // indirect
	//gomodjail:unconfined
	github.com/syndtr/gocapability v0.0.0-20200815063812-42c35b437635 // indirect
	github.com/tinylib/msgp v1.3.0 // indirect
	github.com/vbatts/tar-split v0.12.1 // indirect
	github.com/xhit/go-str2duration/v2 v2.1.0 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.62.0 // indirect
	go.opentelemetry.io/otel v1.37.0 // indirect
	go.opentelemetry.io/otel/metric v1.37.0 // indirect
	go.opentelemetry.io/otel/trace v1.37.0 // indirect
	golang.org/x/exp v0.0.0-20250718183923-645b1fa84792 // indirect
	golang.org/x/mod v0.26.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250721164621-a45f3dfb1074 // indirect
	//gomodjail:unconfined
	google.golang.org/grpc v1.74.2 // indirect
	//gomodjail:unconfined
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	lukechampine.com/blake3 v1.4.1 // indirect
	sigs.k8s.io/yaml v1.6.0 // indirect
	tags.cncf.io/container-device-interface/specs-go v1.0.0 // indirect
)

require github.com/stretchr/testify v1.10.0

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/ebitengine/purego v0.8.4 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/lufia/plan9stats v0.0.0-20250317134145-8bc96cf8fc35 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/power-devops/perfstat v0.0.0-20240221224432-82ca36839d55 // indirect
	github.com/tklauser/go-sysconf v0.3.15 // indirect
	github.com/tklauser/numcpus v0.10.0 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	go.yaml.in/yaml/v2 v2.4.2 // indirect
)

replace github.com/containerd/nerdctl/mod/tigron v0.0.0 => ./mod/tigron

replace github.com/containerd/stargz-snapshotter => github.com/GrigoryEvko/stargz-snapshotter v0.17.3
