package broker

// import (
// 	"reflect"
// 	"task2/internal/logger"
// 	"task2/internal/message"
// 	"task2/internal/subscriber"
// 	"testing"
// )

// func TestBroker_MessageSubscribersByTopic(t *testing.T) {
// 	type args struct {
// 		topicID string
// 		message *message.Message
// 	}
// 	tests := []struct {
// 		name          string
// 		args          args
// 		addTopic      bool
// 		addSubscriber bool
// 		wantErr       bool
// 	}{
// 		{
// 			name: "ok",
// 			args: args{
// 				topicID: "t_id",
// 				message: &message.Message{
// 					ID:    "m_id",
// 					Topic: "t_id",
// 					Body:  "m_body",
// 					Done:  false,
// 				},
// 			},
// 			addTopic:      true,
// 			addSubscriber: true,
// 			wantErr:       false,
// 		},
// 		{
// 			name: "err - no topic",
// 			args: args{
// 				topicID: "t_id",
// 				message: &message.Message{
// 					ID:    "m_id",
// 					Topic: "t_id",
// 					Body:  "m_body",
// 					Done:  false,
// 				},
// 			},
// 			addTopic:      false,
// 			addSubscriber: false,
// 			wantErr:       true,
// 		},
// 		{
// 			name: "err - no subscribers",
// 			args: args{
// 				topicID: "t_id",
// 				message: &message.Message{
// 					ID:    "m_id",
// 					Topic: "t_id",
// 					Body:  "m_body",
// 					Done:  false,
// 				},
// 			},
// 			addTopic:      true,
// 			addSubscriber: false,
// 			wantErr:       true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			logger := logger.NewLogger()
// 			logger.StartNoopLogger()

// 			br := NewBroker(logger)

// 			if tt.addTopic {
// 				br.Topics["t_id"] = subscriber.Subscribers{}
// 			}

// 			if tt.addSubscriber {
// 				br.Topics["t_id"]["s_id"] = subscriber.NewSubscriber(logger)
// 			}

// 			if err := br.MessageSubscribersByTopic(tt.args.topicID, tt.args.message); (err != nil) != tt.wantErr {
// 				t.Errorf("Broker.MessageSubscribersByTopic() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

// func TestBroker_GetSubscribers(t *testing.T) {
// 	tests := []struct {
// 		name          string
// 		addSubscriber bool
// 		want          *subscriber.Subscribers
// 		wantErr       bool
// 	}{
// 		{
// 			name:          "ok",
// 			addSubscriber: true,
// 			want: &subscriber.Subscribers{
// 				"s_id": &subscriber.Subscriber{},
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name:          "err - no subscribers",
// 			addSubscriber: false,
// 			want:          nil,
// 			wantErr:       true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			logger := logger.NewLogger()
// 			logger.StartNoopLogger()

// 			br := NewBroker(logger)

// 			br.Subscribers = make(subscriber.Subscribers)

// 			if tt.addSubscriber {
// 				br.Subscribers["s_id"] = &subscriber.Subscriber{}
// 			}

// 			got, err := br.GetSubscribers()
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Broker.GetSubscribers() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("Broker.GetSubscribers() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestBroker_AddSubscriber(t *testing.T) {

// 	type args struct {
// 		sub *subscriber.Subscriber
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 	}{
// 		{
// 			name: "ok",
// 			args: args{
// 				sub: &subscriber.Subscriber{},
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			logger := logger.NewLogger()
// 			logger.StartNoopLogger()

// 			br := NewBroker(logger)

// 			br.AddSubscriber(tt.args.sub)
// 		})
// 	}
// }

// func TestBroker_GetSubscriber(t *testing.T) {
// 	type args struct {
// 		ID string
// 	}
// 	tests := []struct {
// 		name          string
// 		args          args
// 		addSubscriber bool
// 		want          *subscriber.Subscriber
// 		wantErr       bool
// 	}{
// 		{
// 			name: "ok",
// 			args: args{
// 				ID: "s_id",
// 			},
// 			addSubscriber: true,
// 			want: &subscriber.Subscriber{
// 				ID: "s_id",
// 			},
// 		},
// 		{
// 			name: "err - no subscriber",
// 			args: args{
// 				ID: "s_id",
// 			},
// 			addSubscriber: false,
// 			want:          nil,
// 			wantErr:       true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			logger := logger.NewLogger()
// 			logger.StartNoopLogger()

// 			br := NewBroker(logger)

// 			if tt.addSubscriber {
// 				br.Subscribers["s_id"] = &subscriber.Subscriber{
// 					ID: "s_id",
// 				}
// 			}

// 			got, err := br.GetSubscriber(tt.args.ID)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Broker.GetSubscriber() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("Broker.GetSubscriber() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestBroker_SubscribeToTopic(t *testing.T) {
// 	type args struct {
// 		sub     *subscriber.Subscriber
// 		topicID string
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
// 				sub: &subscriber.Subscriber{
// 					ID:     "s_id",
// 					Active: true,
// 				},
// 				topicID: "t_id",
// 			},
// 			addTopic: true,
// 		},
// 		{
// 			name: "err - sub inactive",
// 			args: args{
// 				sub: &subscriber.Subscriber{
// 					ID:     "s_id",
// 					Active: false,
// 				},
// 				topicID: "t_id",
// 			},
// 			addTopic: true,
// 			wantErr:  true,
// 		},
// 		{
// 			name: "err - no topic",
// 			args: args{
// 				sub: &subscriber.Subscriber{
// 					ID:     "s_id",
// 					Active: false,
// 				},
// 				topicID: "t_id",
// 			},
// 			addTopic: false,
// 			wantErr:  true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			logger := logger.NewLogger()
// 			logger.StartNoopLogger()

// 			br := NewBroker(logger)

// 			if tt.addTopic {
// 				br.Topics["t_id"] = subscriber.Subscribers{}
// 			}

// 			if err := br.SubscribeToTopic(tt.args.sub, tt.args.topicID); (err != nil) != tt.wantErr {
// 				t.Errorf("Broker.SubscribeToTopic() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

// func TestBroker_GetTopicsAndSubscribers(t *testing.T) {
// 	tests := []struct {
// 		name                  string
// 		addTopicAndsubscriber bool
// 		want                  map[string]subscriber.Subscribers
// 		wantErr               bool
// 	}{
// 		{
// 			name:                  "ok",
// 			addTopicAndsubscriber: true,
// 			want: map[string]subscriber.Subscribers{
// 				"t_id": {
// 					"s_id": &subscriber.Subscriber{
// 						ID: "s_id",
// 					},
// 				},
// 			},
// 		},
// 		{
// 			name:                  "err - no topics",
// 			addTopicAndsubscriber: false,
// 			want:                  nil,
// 			wantErr:               true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			logger := logger.NewLogger()
// 			logger.StartNoopLogger()

// 			br := NewBroker(logger)

// 			if tt.addTopicAndsubscriber {
// 				br.Topics["t_id"] = subscriber.Subscribers{
// 					"s_id": &subscriber.Subscriber{
// 						ID: "s_id",
// 					},
// 				}
// 			}

// 			got, err := br.GetTopicsAndSubscribers()
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Broker.GetTopicsAndSubscribers() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("Broker.GetTopicsAndSubscribers() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestBroker_AddTopic(t *testing.T) {
// 	type args struct {
// 		topicID string
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
// 				topicID: "t_id",
// 			},
// 		},
// 		{
// 			name: "err - topic id empty",
// 			args: args{
// 				topicID: "",
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "err - topic already exists",
// 			args: args{
// 				topicID: "t_id",
// 			},
// 			addTopic: true,
// 			wantErr:  true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			logger := logger.NewLogger()
// 			logger.StartNoopLogger()

// 			br := NewBroker(logger)

// 			if tt.addTopic {
// 				br.Topics["t_id"] = subscriber.Subscribers{}
// 			}

// 			if err := br.AddTopic(tt.args.topicID); (err != nil) != tt.wantErr {
// 				t.Errorf("Broker.AddTopic() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }
