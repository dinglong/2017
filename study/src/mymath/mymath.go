package mymath

func Add(a, b int) int {
    return a + b
}
func Max(a, b int) (ret int) {
    ret = a
    if b > a {
        ret = b
    }
    return
}
