package p2591

func DistMoney(money int, children int) int {
	if children == 1 && money == 4 {
		return -1
	}

	money -= children
	if money < 0 {
		return -1
	}

	cnt := 0

	for money >= 7 && cnt < children {
		cnt++
		money -= 7
	}

	if money > 0 {
		if cnt < children {
			if money == 3 && cnt == children-1 {
				cnt--
			}
		} else {
			cnt--
		}
	}

	return cnt
}
