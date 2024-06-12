package algorithm

func Allocate(points map[string][]float64) map[string]int{
	addr2shard := make(map[string]int)
	for addr,v := range(points) {
		max := 0.0
		for shard,point := range(v) {
			if(point>max) {
				max = point
				addr2shard[addr] = shard
			}
		}
	}
	return addr2shard
}