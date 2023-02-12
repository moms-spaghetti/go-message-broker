package protocols

import (
	"encoding/json"
	"reflect"
	"task2/internal/broker"
	"task2/internal/logger"
	"task2/internal/message"
	"task2/internal/subscriber"
	"testing"
)

func Test_protocolHelpers_buildJsonResponse(t *testing.T) {
	type args struct {
		err  error
		data interface{}
	}
	tests := []struct {
		name  string
		args  args
		want1 jsonResponse
	}{
		{
			name: "ok no error",
			args: args{
				err:  nil,
				data: "hello world",
			},
			want1: jsonResponse{
				Err:    "",
				Status: 200,
				Data:   "hello world",
			},
		},
		{
			name: "topic not exist",
			args: args{
				err:  broker.ErrTopicNotExist,
				data: "",
			},
			want1: jsonResponse{
				Err:    "topic does not exist",
				Status: 404,
				Data:   nil,
			},
		},
		{
			name: "message has not body",
			args: args{
				err:  message.ErrMessageBodyEmpty,
				data: "",
			},

			want1: jsonResponse{
				Err:    "message body empty",
				Status: 400,
				Data:   nil,
			},
		},
		{
			name: "err method forbidden",
			args: args{
				err:  subscriber.ErrSubNoMessages,
				data: "",
			},

			want1: jsonResponse{
				Err:    "no messages for topic",
				Status: 404,
				Data:   nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := logger.NewLogger()
			logger.StartNoopLogger()

			got1 := BuildJsonResponse(tt.args.err, tt.args.data, logger)

			// convert back to jsonResponse else checking bytes vs bytes
			var assertResponse jsonResponse

			if err := json.Unmarshal(got1, &assertResponse); err != nil {
				t.Errorf("json.Unmarshal assertResponse err = %v", err)

				return
			}

			if !reflect.DeepEqual(assertResponse, tt.want1) {
				t.Errorf("protocolHelpers.buildJsonResponse() assertResponse = %v, want %v", assertResponse, tt.want1)
			}
		})
	}
}
