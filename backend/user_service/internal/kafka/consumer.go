package kafka_listener

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	handlers "user_service/internal/handlers"

	"github.com/rs/zerolog/log"
	kafka "github.com/segmentio/kafka-go"
)

func getKafkaReader(kafkaURL string, topics []string, groupID string) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		GroupID:     groupID,
		GroupTopics: topics,
		MinBytes:    10e3, // 10KB
		MaxBytes:    10e6, // 10MB
	})
}

func RunKafkaListener(db *sql.DB) {
	kafkaURL := fmt.Sprintf("%v:%v", os.Getenv("KAFKA_HOST"), os.Getenv("KAFKA_PORT"))
	topics := []string{"user-created", "user-deleted"}
	groupID := "1"

	reader := getKafkaReader(kafkaURL, topics, groupID)

	defer reader.Close()

	log.Info().Msg("Start consuming kafka topic")

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatal().Err(err).Msg("kafka message reading fatal")
		}

		messageInfo := fmt.Sprintf(
			"message at topic:%v partition:%v offset:%v	%s = %s\n",
			m.Topic,
			m.Partition,
			m.Offset,
			string(m.Key),
			string(m.Value),
		)
		log.Info().Msg(messageInfo)

		switch m.Topic {
		case "user-created":
			log.Info().Msg("user-created message received")

			err = handlers.AddUser(db, m.Value)
			if err != nil {
				log.Error().Msg("user-created processing failed")
			}
		case "user-deleted":
			log.Info().Msg("user-deleted message received")

			err = handlers.DeleteUser(db, m.Value)
			if err != nil {
				log.Error().Msg("user-deleted processing failed")
			}
		default:
			log.Error().Msg("topic undefined")
		}
	}
}
