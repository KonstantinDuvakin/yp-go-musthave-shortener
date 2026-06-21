package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinksStorage_GetLink(t *testing.T) {
	type fields struct {
		Links map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		arg    string
		want   string
	}{
		{
			name: "Get link",
			fields: fields{
				Links: map[string]string{
					"key": "value",
				},
			},
			arg:  "key",
			want: "value",
		},
		{
			name: "Get link from empty links map",
			fields: fields{
				Links: nil,
			},
			arg:  "empty",
			want: "",
		},
		{
			name: "Get nonexistent link",
			fields: fields{
				Links: map[string]string{
					"key": "value",
				},
			},
			arg:  "nonexistent",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ls := &LinksStorage{
				Links: tt.fields.Links,
			}

			assert.Equal(t, tt.want, ls.GetLink(tt.arg))
		})
	}
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func TestLinksStorage_GetOrCreateLink(t *testing.T) {
	type fields struct {
		Links         map[string]string
		RevertedLinks map[string]string
	}

	tests := []struct {
		name   string
		fields fields
		arg    []byte
		exist  bool
		want   string
	}{
		{
			name: "Get link",
			fields: fields{
				Links: map[string]string{
					"key": "value",
				},
				RevertedLinks: map[string]string{
					"value": "key",
				},
			},
			arg:   []byte("value"),
			exist: true,
			want:  "key",
		},
		{
			name: "Create link",
			fields: fields{
				Links:         nil,
				RevertedLinks: nil,
			},
			arg:   []byte("value"),
			exist: false,
		},
		{
			name: "Create link with non nil links",
			fields: fields{
				Links: map[string]string{
					"key": "value",
				},
				RevertedLinks: map[string]string{
					"value": "key",
				},
			},
			arg:   []byte("value1"),
			exist: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ls := &LinksStorage{
				Links:         tt.fields.Links,
				RevertedLinks: tt.fields.RevertedLinks,
			}

			mapLenBefore := len(ls.Links)

			link := ls.GetOrCreateLink(tt.arg)

			mapLenAfter := len(ls.Links)

			if tt.exist {
				assert.Equal(t, tt.want, link)
				assert.Equal(t, mapLenBefore, mapLenAfter, "должна добавиться ровно одна запись")
			} else {
				assert.Len(t, link, 8)
				for _, c := range link {
					assert.Contains(t, letters, string(c), "символ %q не из алфавита", string(c))
				}

				got, ok := ls.Links[link]

				assert.True(t, ok, "ключ не добавлен в мапу")
				assert.Equal(t, string(tt.arg), got, "по ключу лежит не то значение")

				assert.Equal(t, mapLenBefore+1, mapLenAfter, "должна добавиться ровно одна запись")
			}
		})
	}
}

func TestLinksStorage_saveLink(t *testing.T) {
	type fields struct {
		Links map[string]string
	}
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name       string
		fields     fields
		args       []args
		wantLength int
	}{
		{
			name: "Save link",
			fields: fields{
				Links: map[string]string{},
			},
			args: []args{
				{
					key:   "key",
					value: "value",
				},
			},
			wantLength: 1,
		},
		{
			name: "Save link with nil storage",
			fields: fields{
				Links: nil,
			},
			args: []args{
				{
					key:   "key",
					value: "value",
				},
			},
			wantLength: 1,
		},
		{
			name: "Save link two times with same key",
			fields: fields{
				Links: map[string]string{},
			},
			args: []args{
				{
					key:   "key",
					value: "value",
				},
				{
					key:   "key",
					value: "value1",
				},
			},
			wantLength: 1,
		},
		{
			name: "Save link two times with different key",
			fields: fields{
				Links: map[string]string{},
			},
			args: []args{
				{
					key:   "key",
					value: "value",
				},
				{
					key:   "key1",
					value: "value1",
				},
			},
			wantLength: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ls := &LinksStorage{
				Links: tt.fields.Links,
			}

			for _, arg := range tt.args {
				ls.saveLink(arg.key, arg.value)
			}

			assert.Len(t, ls.Links, tt.wantLength)

			last := tt.args[len(tt.args)-1]
			assert.Equal(t, last.value, ls.Links[last.key])
		})
	}
}
