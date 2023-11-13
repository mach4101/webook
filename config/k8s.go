package config

//
// var Config = config{
// 	DB: DBConfig{
// 		DSN: "root:root@tcp(webook-mysql:11309)/webook",
// 	},
// 	Redis: RedisConfig{
// 		Addr: "webook-redis:11479",
// 	},
// }

// var Config = config{
// 	DB: DBConfig{
// 		DSN: "root:root@tcp(localhost:30002)/webook",
// 	},
// 	Redis: RedisConfig{
// 		Addr: "localhost:30003",
// 	},
// }

var Config = config{
	DB: DBConfig{
		DSN: "root:root@tcp(localhost:13316)/webook",
	},
	Redis: RedisConfig{
		Addr: "localhost:6379",
	},
}
