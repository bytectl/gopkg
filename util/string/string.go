package string

func RemoveDuplicate(strs []string) []string {
	var result []string
	flagMap := map[string]struct{}{}

	for _, str := range strs {
		if _, ok := flagMap[str]; !ok {
			flagMap[str] = struct{}{}
			result = append(result, str)
		}
	}
	return result
}
