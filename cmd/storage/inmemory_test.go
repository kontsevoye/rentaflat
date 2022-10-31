package storage

import (
	"github.com/kontsevoye/rentaflat/cmd/parser"
	"reflect"
	"testing"
)

func TestInMemoryStorage_Has(t *testing.T) {
	type fields struct {
		flats map[string]wrappedFlat
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
				flats: map[string]wrappedFlat{
					"1": {
						Flat:  parser.Flat{},
						isNew: true,
					},
				},
			},
			args: args{"1"},
			want: true,
		},
		{
			name: "Non existing item check",
			fields: fields{
				flats: map[string]wrappedFlat{
					"1": {
						Flat:  parser.Flat{},
						isNew: true,
					},
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
		flats map[string]wrappedFlat
	}
	type args struct {
		flats []parser.Flat
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "New item in empty storage",
			fields: fields{
				flats: map[string]wrappedFlat{},
			},
			args: args{
				[]parser.Flat{
					{
						Id: "1",
					},
				},
			},
			want: 1,
		},
		{
			name: "Override existing item",
			fields: fields{
				flats: map[string]wrappedFlat{
					"1": {},
				},
			},
			args: args{
				[]parser.Flat{
					{
						Id: "1",
					},
				},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &InMemoryStorage{
				flats: tt.fields.flats,
			}
			if got := s.Store(tt.args.flats); got != tt.want {
				t.Errorf("Store() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInMemoryStorage_GetAllNew(t *testing.T) {
	type fields struct {
		flats map[string]wrappedFlat
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]parser.Flat
	}{
		{
			name: "No new flats in storage",
			fields: fields{
				flats: map[string]wrappedFlat{
					"1": {
						Flat:  parser.Flat{},
						isNew: false,
					},
				},
			},
			want: map[string]parser.Flat{},
		},
		{
			name: "New flats in storage",
			fields: fields{
				flats: map[string]wrappedFlat{
					"1": {
						Flat:  parser.Flat{},
						isNew: true,
					},
				},
			},
			want: map[string]parser.Flat{
				"1": parser.Flat{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &InMemoryStorage{
				flats: tt.fields.flats,
			}
			if got := s.GetAllNew(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllNew() = %v, want %v", got, tt.want)
			}
			if got := len(s.GetAllNew()); got != 0 {
				t.Errorf("Second call len(GetAllNew()) = %v, want %v", got, 0)
			}
		})
	}
}
