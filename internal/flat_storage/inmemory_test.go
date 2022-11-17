package flat_storage

import (
	"github.com/kontsevoye/rentaflat/internal/flat_parser"
	"testing"
)

func TestInMemoryStorage_Has(t *testing.T) {
	type fields struct {
		flats map[string]flat_parser.Flat
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Existing item check",
			fields: fields{
				flats: map[string]flat_parser.Flat{
					"1": {},
				},
			},
			args: args{"1"},
			want: true,
		},
		{
			name: "Non existing item check",
			fields: fields{
				flats: map[string]flat_parser.Flat{
					"1": {},
				},
			},
			args: args{"2"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &InMemoryStorage{
				flats: tt.fields.flats,
			}
			if got := s.Has(tt.args.id); got != tt.want {
				t.Errorf("Has() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInMemoryStorage_Store(t *testing.T) {
	type fields struct {
		flats map[string]flat_parser.Flat
	}
	type args struct {
		flat flat_parser.Flat
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "New item in empty flat_storage",
			fields: fields{
				flats: map[string]flat_parser.Flat{},
			},
			args: args{
				flat_parser.Flat{
					Id: "1",
				},
			},
			want: 1,
		},
		{
			name: "Override existing item",
			fields: fields{
				flats: map[string]flat_parser.Flat{
					"1": {},
				},
			},
			args: args{
				flat_parser.Flat{
					Id: "1",
				},
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &InMemoryStorage{
				flats: tt.fields.flats,
			}
			s.Store(tt.args.flat)
			if got := s.Count(); got != tt.want {
				t.Errorf("Store() = %v, want %v", got, tt.want)
			}
		})
	}
}
