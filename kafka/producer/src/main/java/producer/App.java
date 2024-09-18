package producer;

import org.apache.kafka.clients.producer.KafkaProducer;
import org.apache.kafka.clients.producer.ProducerRecord;
import org.apache.kafka.clients.producer.RecordMetadata;
import org.apache.kafka.common.serialization.IntegerSerializer;
import org.apache.kafka.common.serialization.StringSerializer;

import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.Properties;
import java.util.concurrent.Future;

public class App {
    public static KafkaProducer<Integer, String> producer;

    public App() throws UnknownHostException {
        // Initialize Producer
        Properties config = new Properties();
        config.put("client.id", InetAddress.getLocalHost().getHostName());
        config.put("bootstrap.servers", "kafka:9092");
        config.put("key.serializer", IntegerSerializer.class.getName());
        config.put("value.serializer", StringSerializer.class.getName());
        config.put("acks", "all");
        producer = new KafkaProducer<>(config);
    }

    public static void main(String[] args) throws UnknownHostException {
        new App();
        System.out.println("Start publishing messages to Kafka");

        while (true) {
            try {
                final ProducerRecord<Integer, String> record = new ProducerRecord<>("test-topic", 1, "{}");
                producer.send(record, (metadata, e) -> System.out.printf("Message sent to topic: %s timestamp: %s topicPartition %s\n", metadata.topic(), metadata.timestamp(), metadata.partition()));
                Thread.sleep(10000);
            } catch (Exception e) {
                System.out.println("Error: " + e);
                producer.close();
                break;
            }
        }
    }
}
