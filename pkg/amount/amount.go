// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package amount

// FenToYuan converts price from fen (int) to yuan (float64)
func FenToYuan(fen int) float64 {
	return float64(fen) / 100.0
}

// YuanToFen converts price from yuan (float64) to fen (int)
func YuanToFen(yuan float64) int {
	return int(yuan * 100)
}
