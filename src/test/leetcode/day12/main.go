package main

import "fmt"

func intToRoman(num int) string {
	money := []int{1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1}
	str := []string{"M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"}
	//list := []string{
	//	1000: "M",
	//	900:  "CM",
	//	500:  "D",
	//	400:  "CD",
	//	100:  "C",
	//	90:   "XC",
	//	50:   "L",
	//	40:   "XL",
	//	10:   "X",
	//	9:    "IX",
	//	5:    "V",
	//	4:    "IV",
	//	1:    "I",
	//}
	var res string
	for k, v := range money {
		for num/v >= 1 {
			res += str[k]
			num -= v
		}
	}
	return res
}

func intToRoman2(num int) string {

	d := [4][]string{
		[]string{"", "I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX"},
		[]string{"", "X", "XX", "XXX", "XL", "L", "LX", "LXX", "LXXX", "XC"},
		[]string{"", "C", "CC", "CCC", "CD", "D", "DC", "DCC", "DCCC", "CM"},
		[]string{"", "M", "MM", "MMM"},
	}
	return d[3][num/1000] +
		d[2][num/100%10] +
		d[1][num/10%10] +
		d[0][num%10]
}

func main() {
	fmt.Println(intToRoman2(1874))
}
