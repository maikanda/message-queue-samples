package producer;

import org.apache.pulsar.client.api.Producer;
import org.apache.pulsar.client.api.PulsarClient;
import org.apache.pulsar.client.api.PulsarClientException;
import org.apache.pulsar.client.api.Schema;

public class App {
    public static Producer<String> producer;
    public static PulsarClient client;

    public App() throws PulsarClientException {
         client = PulsarClient.builder()
                .serviceUrl("pulsar://pulsar:6650")
                .build();
         producer = client.newProducer(Schema.STRING)
                 .topic("persistent://my-tenant/my-namespace/my-topic")
                 .enableBatching(false)
                 .create();
    }

    public static void main(String[] args) throws PulsarClientException {
        new App();
        System.out.println("Start publishing messages to Pulsar");

        while (true) {
            try {
                for (int i = 0; i < 3; i ++) {
                    producer.newMessage()
                            .key(String.format("key-%d", i))
                            .value("Hello world")
                            .send();
                }
                Thread.sleep(10000);
            } catch (Exception e) {
                System.out.println("Error: " + e);
                producer.close();
                client.close();
                break;
            }
        }
    }
}
