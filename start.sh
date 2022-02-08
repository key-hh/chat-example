!bin

cd docker/kafka-docker

sudo docker-compose -f docker-compose-single-broker.yml up -d


# kafka test
cd /Users/kakao/test/kafka/kafka_2.13-2.8.1/bin

kafka-topics.sh --create --zookeeper localhost:2181 --replication-factor 1 --partitions 1 --topic test_topic

kafka-console-producer.sh --topic test_topic --broker-list localhost:9092

kafka-console-consumer.sh --topic test_topic --bootstrap-server localhost:9092 --from-beginning