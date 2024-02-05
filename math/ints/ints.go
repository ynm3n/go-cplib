package ints

func ExtGcd(a, b int) (d, x, y int) {
	if b == 0 {
		return a, 1, 0
	}
	d, x, y = ExtGcd(b, a%b)
	return d, y, x - (a/b)*y
}

func PowMod(a, n, m int) int {
	res := 1 % m
	for b := a % m; n > 0; {
		if n&1 > 0 {
			res *= b
			res %= m
		}
		n >>= 1
		b *= b
		b %= m
	}
	return res
}
