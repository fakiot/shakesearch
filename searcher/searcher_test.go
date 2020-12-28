package searcher

import (
	"bytes"
	"index/suffixarray"
	"reflect"
	"strings"
	"testing"
)

const testText = `HORATIO.
Here, sweet lord, at your service.

HAMLET.
Horatio, thou art e’en as just a man
As e’er my conversation cop’d withal.

HORATIO.
O my dear lord.

HAMLET.
Nay, do not think I flatter;
For what advancement may I hope from thee,
That no revenue hast, but thy good spirits
To feed and clothe thee? Why should the poor be flatter’d?
No, let the candied tongue lick absurd pomp,
And crook the pregnant hinges of the knee
Where thrift may follow fawning. Dost thou hear?
Since my dear soul was mistress of her choice,
And could of men distinguish, her election
Hath seal’d thee for herself. For thou hast been
As one, in suffering all, that suffers nothing,
A man that Fortune’s buffets and rewards
Hast ta’en with equal thanks. And bles’d are those
Whose blood and judgment are so well co-mingled
That they are not a pipe for Fortune’s finger
To sound what stop she please. Give me that man
That is not passion’s slave, and I will wear him
In my heart’s core, ay, in my heart of heart,
As I do thee. Something too much of this.
There is a play tonight before the King.
One scene of it comes near the circumstance
Which I have told thee, of my father’s death.
I prythee, when thou see’st that act a-foot,
Even with the very comment of thy soul
Observe mine uncle. If his occulted guilt
Do not itself unkennel in one speech,
It is a damned ghost that we have seen;
And my imaginations are as foul
As Vulcan’s stithy. Give him heedful note;
For I mine eyes will rivet to his face;
And after we will both our judgments join
In censure of his seeming.`

// Test_binarySearch unit test for binary search
func Test_binarySearch(t *testing.T) {
	type args struct {
		t int
		p []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "Match",
			args: args{
				t: 5,
				p: []int{2, 4, 6, 8},
			},
			want: []int{4, 6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := binarySearch(tt.args.t, tt.args.p); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("binarySearch() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestSearcher_Search search unit test
func TestSearcher_Search(t *testing.T) {
	type fields struct {
		CompleteWorks string
		SuffixArray   *suffixarray.Index
	}
	type args struct {
		query string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Match",
			fields: fields{
				CompleteWorks: testText,
				SuffixArray:   suffixarray.New(bytes.ToUpper([]byte(testText))),
			},
			args: args{
				query: "hAmLeT",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Searcher{
				CompleteWorks: tt.fields.CompleteWorks,
				SuffixArray:   tt.fields.SuffixArray,
			}
			got := s.Search(tt.args.query)
			for _, v := range got {
				// must find at least one
				if !strings.Contains(strings.ToLower(v), strings.ToLower(tt.args.query)) {
					t.Errorf("Searcher.Search() = %v", got)
				}
			}
		})
	}
}

// TestSearcher_Load load file unit test
func TestSearcher_Load(t *testing.T) {
	type fields struct {
		CompleteWorks string
		SuffixArray   *suffixarray.Index
	}
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Self",
			args: args{
				filename: "searcher_test.go",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Searcher{
				CompleteWorks: tt.fields.CompleteWorks,
				SuffixArray:   tt.fields.SuffixArray,
			}
			if err := s.Load(tt.args.filename); (err != nil) != tt.wantErr {
				t.Errorf("Searcher.Load() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
