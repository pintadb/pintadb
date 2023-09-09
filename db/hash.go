package db

// Hash function using the FNV-1a algorithm
// https://en.wikipedia.org/wiki/Fowler%E2%80%93Noll%E2%80%93Vo_hash_function
func FowlerNollVo32(s string) uint32 {
	const (
		fnvOffset32 uint32 = 2166136261
		fnvPrime32  uint32 = 16777619
	)

	hash := fnvOffset32
	for _, c := range s {
		hash ^= uint32(c)
		hash *= fnvPrime32
	}

	return hash
}
