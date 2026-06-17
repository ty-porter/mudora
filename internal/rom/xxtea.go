package rom

// xxteaDecrypt decrypts an 8-byte block (2 × uint32 LE words) using the given 16-byte key.
// Uses 32 rounds, matching the z3randomizer ASM implementation.
func xxteaDecrypt(block []byte, key [4]uint32) [8]byte {
	v0 := uint32(block[0]) | uint32(block[1])<<8 | uint32(block[2])<<16 | uint32(block[3])<<24
	v1 := uint32(block[4]) | uint32(block[5])<<8 | uint32(block[6])<<16 | uint32(block[7])<<24

	const delta = uint32(0x9E3779B9)

	sum := uint32(0xC6EF3720) // delta * 32, wrapped to uint32
	y := v0
	for sum != 0 {
		e := (sum >> 2) & 3
		// p=1
		z := v0
		v1 -= xxteaMX(sum, y, z, 1, e, key)
		y = v1
		// p=0
		z = v1
		v0 -= xxteaMX(sum, y, z, 0, e, key)
		y = v0
		sum -= delta
	}

	return [8]byte{
		byte(v0), byte(v0 >> 8), byte(v0 >> 16), byte(v0 >> 24),
		byte(v1), byte(v1 >> 8), byte(v1 >> 16), byte(v1 >> 24),
	}
}

func xxteaMX(sum, y, z uint32, p, e uint32, key [4]uint32) uint32 {
	return ((z>>5 ^ y<<2) + (y>>3 ^ z<<4)) ^ ((sum ^ y) + (key[(p&3)^e] ^ z))
}
