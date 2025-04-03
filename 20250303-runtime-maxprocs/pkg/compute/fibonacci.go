package compute

// Fibonacci 计算斐波那契数列，模拟CPU密集型任务
func Fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return Fibonacci(n-1) + Fibonacci(n-2)
}

// FibonacciMatrix 使用矩阵法计算斐波那契数列，复杂度O(log n)
func FibonacciMatrix(n int) int {
	if n <= 1 {
		return n
	}

	a, b := 1, 1
	p, q := 0, 1
	for i := n; i > 0; i >>= 1 {
		if i&1 == 1 {
			p, q = p*a+q*b, p*b+q*(a+b)
		}
		a, b = a*a+b*b, b*b+(a+a)*b
	}
	return p
}

// MatrixMultiply 矩阵乘法，也是CPU密集型操作
func MatrixMultiply(a, b [][]int) [][]int {
	if len(a) == 0 || len(a[0]) == 0 || len(b) == 0 || len(b[0]) == 0 || len(a[0]) != len(b) {
		return nil
	}

	rows, cols := len(a), len(b[0])
	result := make([][]int, rows)
	for i := range result {
		result[i] = make([]int, cols)
	}

	// 执行矩阵乘法
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			for k := 0; k < len(b); k++ {
				result[i][j] += a[i][k] * b[k][j]
			}
		}
	}

	return result
}
