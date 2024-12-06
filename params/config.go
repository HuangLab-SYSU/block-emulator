package params

type ChainConfig struct {
	ChainID           int
	NodeID            string
	ShardID           string
	Shard_num         int
	Malicious_num     int    // per shard
	Path              string // input file path
	Block_interval    int    // millisecond
	MaxBlockSize      int
	Relay_interval    int
	MaxRelayBlockSize int
	MinRelayBlockSize int
	Inject_speed      int // tx count per second
	Reconfig_interval int
	// 有待优化
	Reconfig_time int
	Relay_time    int
}

var (
	ClientAddr = "127.0.0.1:8200"
	NodeTable  = map[string]map[string]string{
		"S0": {
			"N0": "127.0.0.1:8010",
			"N1": "127.0.0.1:8011",
			"N2": "127.0.0.1:8012",
			"N3": "127.0.0.1:8013",
		},
		"S1": {
			"N0": "127.0.0.1:8014",
			"N1": "127.0.0.1:8015",
			"N2": "127.0.0.1:8016",
			"N3": "127.0.0.1:8017",
		},
		// // 中心分片
		// "SC": {
		// 	"N0": "127.0.0.1:8018",
		// 	"N1": "127.0.0.1:8019",
		// 	"N2": "127.0.0.1:8020",
		// 	"N3": "127.0.0.1:8021",
		// },

		"S2": {
			"N0": "127.0.0.1:8018",
			"N1": "127.0.0.1:8019",
			"N2": "127.0.0.1:8020",
			"N3": "127.0.0.1:8021",
		},
		"S3": {
			"N0": "127.0.0.1:8022",
			"N1": "127.0.0.1:8023",
			"N2": "127.0.0.1:8024",
			"N3": "127.0.0.1:8025",
		},
		// 中心分片
		"SC": {
			"N0": "127.0.0.1:8026",
			"N1": "127.0.0.1:8027",
			"N2": "127.0.0.1:8028",
			"N3": "127.0.0.1:8029",
		},
	}

	ShardTable = map[string]int{
		"S0": 0,
		"S1": 1,
		// "SC": 2,
		"S2": 2,
		"S3": 3,
		"SC": 4,
	}
	ShardTableInt2Str = map[int]string{
		0: "S0",
		1: "S1",
		// 2: "SC",
		2: "S2",
		3: "S3",
		4: "SC",
	}

	Config = &ChainConfig{
		ChainID:           77,
		Block_interval:    3000,
		MaxBlockSize:      500,
		MaxRelayBlockSize: 500,
		MinRelayBlockSize: 1,
		Inject_speed:      200,
		Relay_interval:    500,
		Reconfig_interval: 10000,
		Reconfig_time:     10000,
		Relay_time:        2000,
	}

	Init_addrs = []string{
		"171382ed4571b1084bb5963053203c237dba6da9",
		"2185bf3bfda43894efdcc1a3f4a99a7f160bc123",
		"374be1a1d1ac0ff350dc9d0a0be3d059c7082791",
		"42fd9ff72a798780c0dffc68f89b64ba240240dd",
	}
	Init_balance string = "100000000000000000000000000000000000000000000" //40个0
)
