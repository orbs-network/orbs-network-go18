package builders

func EmptyPayloads(num int) [][]byte {
	res := [][]byte{}
	for i := 0; i < num; i++ {
		res = append(res, []byte{})
	}
	return res
}
