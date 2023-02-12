package message

import (
	"reflect"
	"task2/internal/logger"
	"testing"
)

func TestNewMessage(t *testing.T) {
	l := logger.NewLogger()
	l.StartNoopLogger()

	type args struct {
		topic  string
		body   interface{}
		logger *logger.Logger
	}
	tests := []struct {
		name    string
		args    args
		want    Message
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				topic:  "t_id",
				body:   "m_body",
				logger: l,
			},
			want: Message{
				ID:    "m_id",
				Topic: "t_id",
				Body:  "m_body",
				Done:  false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMessage(tt.args.topic, tt.args.body, tt.args.logger)

			if got.ID != "" {
				got.ID = "m_id"
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("NewMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
