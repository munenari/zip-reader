package util

import (
	"fmt"
	"testing"
)

func TestGzipPath(t *testing.T) {
	v := `\\examplenas.local\directory1\somewhere\tmp\archives\長い日本語が混じったZIPファイル名.zip`
	//    YKCsFdXfw3mTqevqi8vm6VPWXvSa4ua2MfCgvhf3JYYXPwYbdY72uENEZuc2Fpy8Es5fnueaPq7oKC5n8RVm9Vi1Vv9aSgUpcbV1o7YKfy3oGbBbwP1itaJhTogrKLoKnveMq9XvtNGPxLHNkS5wPdHDBxVFELtwQsY2SN8t8ugxGUHkWXedBaZm
	res, err := GzipPath(v)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
	res2, err := UnGzipPath(res)
	if err != nil {
		t.Fatal(err)
	}
	if res2 != v {
		t.Error("unexpected result", res2)
	}
	// t.Error()
}
