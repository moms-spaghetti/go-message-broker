package subscriber

// import (
// 	"reflect"
// 	"task2/internal/logger"
// 	"task2/internal/message"
// 	"testing"
// )

// func TestSubscriber_GetTopics(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		addTopic bool
// 		want     []string
// 		wantErr  bool
// 	}{
// 		{
// 			name:     "ok",
// 			addTopic: true,
// 			want:     []string{"t_id"},
// 		},
// 		{
// 			name:    "err - no topics",
// 			want:    nil,
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			logger := logger.NewLogger()
// 			logger.StartNoopLogger()

// 			s := NewSubscriber(logger)

// 			if tt.addTopic {
// 				s.Topics["t_id"] = []*message.Message{}
// 			}

// 			got, err := s.GetTopics()
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Subscriber.GetTopics() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("Subscriber.GetTopics() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestSubscriber_SubscribeToTopic(t *testing.T) {
// 	l := logger.NewLogger()
// 	l.StartNoopLogger()

// 	type args struct {
// 		topicID string
// 		logger  *logger.Logger
// 	}
// 	tests := []struct {
// 		name     string
// 		addTopic bool
// 		args     args
// 		wantErr  bool
// 	}{
// 		{
// 			name: "ok ",
// 			args: args{
// 				topicID: "t_id",
// 				logger:  l,
// 			},
// 		},
// 		{
// 			name:     "ok ",
// 			addTopic: true,
// 			args: args{
// 				topicID: "t_id",
// 				logger:  l,
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := NewSubscriber(l)

// 			if tt.addTopic {
// 				s.Topics["t_id"] = []*message.Message{}
// 			}

// 			if err := s.SubscribeToTopic(tt.args.topicID, tt.args.logger); (err != nil) != tt.wantErr {
// 				t.Errorf("Subscriber.SubscribeToTopic() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

// func TestSubscriber_GetMessages(t *testing.T) {
// 	type args struct {
// 		topicID string
// 	}
// 	tests := []struct {
// 		name     string
// 		args     args
// 		addTopic bool
// 		want     []*message.Message
// 		wantErr  bool
// 	}{
// 		{
// 			name: "ok",
// 			args: args{
// 				topicID: "t_id",
// 			},
// 			addTopic: true,
// 			want: []*message.Message{
// 				{
// 					ID: "m_id",
// 				},
// 			},
// 		},
// 		{
// 			name: "err - topic doesn't exist",
// 			args: args{
// 				topicID: "t_id_1",
// 			},
// 			addTopic: true,
// 			want:     nil,
// 			wantErr:  true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			logger := logger.NewLogger()
// 			logger.StartNoopLogger()

// 			s := NewSubscriber(logger)

// 			if tt.addTopic {
// 				s.Topics["t_id"] = []*message.Message{
// 					{
// 						ID: "m_id",
// 					},
// 				}
// 			}

// 			got, err := s.GetMessages(tt.args.topicID)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Subscriber.GetMessages() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("Subscriber.GetMessages() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestSubscriber_GetMessageCount(t *testing.T) {
// 	type args struct {
// 		topicID string
// 	}
// 	tests := []struct {
// 		name     string
// 		args     args
// 		addTopic bool
// 		want     int
// 	}{
// 		{
// 			name: "ok",
// 			args: args{
// 				topicID: "t_id",
// 			},
// 			addTopic: true,
// 			want:     1,
// 		},
// 		{
// 			name: "ok - no messages",
// 			args: args{
// 				topicID: "t_id",
// 			},
// 			addTopic: false,
// 			want:     0,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			logger := logger.NewLogger()
// 			logger.StartNoopLogger()

// 			s := NewSubscriber(logger)

// 			if tt.addTopic {
// 				s.Topics["t_id"] = []*message.Message{
// 					{
// 						ID: "m_id",
// 					},
// 				}
// 			}

// 			if got := s.GetMessageCount(tt.args.topicID); got != tt.want {
// 				t.Errorf("Subscriber.GetMessageCount() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestSubscriber_GetNextMessage(t *testing.T) {
// 	type args struct {
// 		topicID string
// 	}
// 	tests := []struct {
// 		name     string
// 		args     args
// 		addTopic bool
// 		want     message.Message
// 	}{
// 		{
// 			name: "ok",
// 			args: args{
// 				topicID: "t_id",
// 			},
// 			addTopic: true,
// 			want: message.Message{
// 				ID: "m_id_0",
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			logger := logger.NewLogger()
// 			logger.StartNoopLogger()

// 			s := NewSubscriber(logger)

// 			if tt.addTopic {
// 				s.Topics["t_id"] = []*message.Message{
// 					{
// 						ID: "m_id_0",
// 					},
// 					{
// 						ID: "m_id_1",
// 					},
// 				}
// 			}

// 			if got := s.GetNextMessage(tt.args.topicID); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("Subscriber.GetNextMessage() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestSubscriber_DeleteMessage(t *testing.T) {
// 	l := logger.NewLogger()
// 	l.StartNoopLogger()

// 	type args struct {
// 		topicID   string
// 		messageID string
// 		logger    *logger.Logger
// 	}
// 	tests := []struct {
// 		name     string
// 		args     args
// 		addTopic bool
// 		wantErr  bool
// 	}{
// 		{
// 			name: "ok",
// 			args: args{
// 				topicID:   "t_id",
// 				messageID: "m_id",
// 				logger:    l,
// 			},
// 			addTopic: true,
// 		},
// 		{
// 			name: "ok",
// 			args: args{
// 				topicID:   "t_id",
// 				messageID: "m_id",
// 				logger:    l,
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := NewSubscriber(l)

// 			if tt.addTopic {
// 				s.Topics["t_id"] = []*message.Message{
// 					{
// 						ID: "m_id",
// 					},
// 				}
// 			}

// 			if err := s.DeleteMessage(tt.args.topicID, tt.args.messageID, tt.args.logger); (err != nil) != tt.wantErr {
// 				t.Errorf("Subscriber.DeleteMessage() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

// func TestSubscriber_ValidateTopicAndMessageCount(t *testing.T) {
// 	type args struct {
// 		topicID string
// 	}
// 	tests := []struct {
// 		name       string
// 		args       args
// 		addTopic   bool
// 		addMessage bool
// 		wantErr    bool
// 	}{
// 		{
// 			name: "ok",
// 			args: args{
// 				topicID: "t_id",
// 			},
// 			addTopic:   true,
// 			addMessage: true,
// 		},
// 		{
// 			name: "err - topic not exist",
// 			args: args{
// 				topicID: "t_id_1",
// 			},
// 			addTopic: true,
// 			wantErr:  true,
// 		},
// 		{
// 			name: "err - no messages",
// 			args: args{
// 				topicID: "t_id",
// 			},
// 			addTopic:   true,
// 			addMessage: false,
// 			wantErr:    true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			logger := logger.NewLogger()
// 			logger.StartNoopLogger()

// 			s := NewSubscriber(logger)

// 			if tt.addTopic {
// 				s.Topics["t_id"] = []*message.Message{}
// 			}

// 			if tt.addMessage {
// 				s.Topics["t_id"] = append(s.Topics["t_id"], &message.Message{
// 					ID: "m_id",
// 				})
// 			}

// 			if err := s.ValidateTopicAndMessageCount(tt.args.topicID); (err != nil) != tt.wantErr {
// 				t.Errorf("Subscriber.ValidateTopicAndMessageCount() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }
