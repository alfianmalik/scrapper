package main

import (
	"bytes"
	"testing"
)

func Test_realMain(t *testing.T) {
	type args struct {
		products []string
	}
	tests := []struct {
		name   string
		args   args
		expect string
	}{
		// TODO: Add test cases.
		{"Zjiang ZJ-5890T", args{[]string{"Zjiang ZJ-5890T"}}, "Finished https://www.tokopedia.com/search?st=product&q=Zjiang+ZJ-5890T"},
		{"Sennheiser HD 202", args{[]string{"Sennheiser HD 202"}}, "Finished https://www.tokopedia.com/search?st=product&q=Sennheiser+HD+202"},
		{"Blank", args{[]string{""}}, "Finished https://www.tokopedia.com/search?st=product&q="},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output bytes.Buffer
			realMain(&output, tt.args.products)
			if tt.expect != output.String() {
				t.Errorf("got %s but expected %s", output.String(), tt.expect)
			}
		})
	}
}
