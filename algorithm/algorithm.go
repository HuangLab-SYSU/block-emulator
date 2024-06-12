package algorithm

import "blockEmulator/utils"

var flag = 1

func MigrationAlgorithm(before map[string]bool, shardID int) map[string]int {
	out := make(map[string]int)
	// for k := range before {
	// 	if k == "c418969d5f8948d9a40465f0a432d510b7e80b36" {
	// 		if shardID == 0 && flag == 1 {
	// 			out[k] = 1 - shardID
	// 			flag++
	// 		}
	// 		break
	// 	}
	// }
	// return out

	for i := 0; i < 184379; i++ {
		if i%10 == 0 {
			var addr string
			if i%20==0 {
				addr = utils.Int2hexString(i)
			}else {
				addr = utils.Int2hexString(i+5)
			}
			if before[addr] {
				out[addr] = 1- shardID
			}
			
		}
	}
	return out


}

func Algorithm2(old map[string]int, shardID int) map[string]int {

	new := old
	// for k := range new {
	// 	// if k == "c418969d5f8948d9a40465f0a432d510b7e80b36" {
	// 	// 	new[k]=1
	// 	// }else if k == "25eaff5b179f209cf186b1cdcbfa463a69df4c45" {
	// 	// 	new[k]=0
	// 	// }
	// 	new[k] = 1-new[k]

	// }

	for i := 0; i < 184379; i++ {
		if i%10 == 0 {
			if i%20==0 {
				addr := utils.Int2hexString(i)
				new[addr] = 1 - new[addr]
			}else {
				addr := utils.Int2hexString(i+5)
				new[addr] = 1 - new[addr]
			}
			
		}
	}

	return new
}





func Pagerank(graph map[string]map[string]int, addrs []string, addr2shard map[string]int, numbda float64, iters, shard_num int) map[string][]float64 {

	// 每个账户在各个分片的得分, 初始都为0
	points := make(map[string][]float64)
	for _, addr := range addrs {
		points[addr] = make([]float64, shard_num)
	}

	// 每个分片账户的数量
	shard_size := make([]int, shard_num)
	for _, shard := range addr2shard {
		shard_size[shard]++
	}

	var w float64
	// 迭代次数iters
	for iter := 0; iter<iters; iter++ {
		for _, addr := range addrs {
			for shard := 0; shard<shard_num; shard++ {

				if addr2shard[addr] == shard {
					w = 1/float64(shard_size[shard])
				}else {
					w = 0
				}

				points[addr][shard] = (1-numbda)*w + numbda*sum(graph[addr], points, shard)
			}
		}
	}
	return points
}


func sum(out map[string]int, points map[string][]float64, shard int) float64{
	total := 0.0
	count := 0
	for i,val := range(out) {
		if val != 0 {
			total += float64(val)*points[i][shard]
			count += val
		}
	}
	return total/float64(count)
}