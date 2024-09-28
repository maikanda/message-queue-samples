package producer;

import io.nats.client.*;
import io.nats.client.api.PublishAck;
import io.nats.client.impl.NatsMessage;

import java.io.IOException;
import java.nio.charset.StandardCharsets;

public class App {
    public static Connection natsCon;
    public static JetStream js;

    public App() throws IOException, InterruptedException {
        natsCon = Nats.connect("nats://nats:4222");
        js = natsCon.jetStream();
    }

    public static void main(String[] args) throws IOException, InterruptedException {
        new App();
        System.out.println("Start publishing messages to NATS");

        while (true) {
            try {
                String msgStr = "hello world" + System.currentTimeMillis();
                Message msg = NatsMessage.builder()
                        .subject("sub.my-topic")
                        .data(msgStr.getBytes(StandardCharsets.UTF_8))
                        .build();
                PublishOptions po = PublishOptions.builder()
                        .stream("stream-1")
                        .expectedStream("stream-1")
                        .build();
                PublishAck res = js.publish(msg, po);
                System.out.println("Published message: " + res.toString());
                Thread.sleep(10000);
            } catch (Exception e) {
                System.out.println("Error: " + e);
                natsCon.close();
                break;
            }
        }
    }
}
