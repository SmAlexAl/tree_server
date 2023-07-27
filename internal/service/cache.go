package service

func GetChildren(res []string, id string, tree map[string][]string) []string {
	for _, val := range tree[id] {
		if _, ok := tree[val]; ok {
			res = GetChildren(res, val, tree)
		}
		res = append(res, val)
	}

	return res
}
