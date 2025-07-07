package tg

import (
	"fmt"
	"strings"
)

func (tg *TgController) HandlePair(pair string) string {
	pair = strings.ToUpper(pair)
	if strings.HasSuffix(pair, "/USDT:USDT") {
		return pair
	}
	return fmt.Sprintf("%s/USDT:USDT", pair)
}
