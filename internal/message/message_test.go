package message

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestUnit_MarshalData(t *testing.T) {
	type fields struct {
		data []byte
	}
	type args struct {
		data interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "positive",
			fields: fields{
				data: make([]byte, 16),
			},
			args: args{
				data: struct {
					Data string `json:"data" validate:"required"`
				}{Data: "test"},
			},
			want:    `{"data":"test"}`,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New(tt.fields.data, nil)
			err := m.MarshalData(tt.args.data)
			assert.JSONEq(t, tt.want, string(m.GetData()), "MarshalData() = %v, want %v", string(m.GetData()), tt.want)
			tt.wantErr(t, err, "MarshalData() error")
		})
	}
}

func TestUnit_UnmarshalData(t *testing.T) {
	type fields struct {
		data      []byte
		validator *validator.Validate
	}
	type args struct {
		data interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *struct {
			Data string `json:"data" validate:"required"`
		}
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "positive",
			fields: fields{
				data:      []byte(`{"data":"test"}`),
				validator: validator.New(),
			},
			args: args{
				data: &struct {
					Data string `json:"data" validate:"required"`
				}{},
			},
			want: &struct {
				Data string `json:"data" validate:"required"`
			}{Data: "test"},
			wantErr: assert.NoError,
		},
		{
			name: "invalid JSON",
			fields: fields{
				data:      []byte(`invalid json`),
				validator: validator.New(),
			},
			args: args{
				data: &struct {
					Data string `json:"data" validate:"required"`
				}{},
			},
			want: &struct {
				Data string `json:"data" validate:"required"`
			}{Data: ""},
			wantErr: assert.Error,
		},
		{
			name: "invalid struct",
			fields: fields{
				data:      []byte(`{"data": 123}`),
				validator: validator.New(),
			},
			args: args{
				data: &struct {
					Data string `json:"data" validate:"required"`
				}{},
			},
			want: &struct {
				Data string `json:"data" validate:"required"`
			}{Data: ""},
			wantErr: assert.Error,
		},
		{
			name: "missing required field",
			fields: fields{
				data:      []byte(`{"missingField": "test"}`),
				validator: validator.New(),
			},
			args: args{
				data: &struct {
					Data string `json:"data" validate:"required"`
				}{},
			},
			want: &struct {
				Data string `json:"data" validate:"required"`
			}{Data: ""},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New(tt.fields.data, tt.fields.validator)
			err := m.UnmarshalData(tt.args.data)
			assert.Equalf(t, tt.want, tt.args.data, "UnmarshalData() = %v, want %v", tt.args.data, tt.want)
			tt.wantErr(t, err, "UnmarshalData() error")
		})
	}
}
