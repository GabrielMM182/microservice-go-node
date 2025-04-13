const amqp = require('amqplib');
const logger = require('../utils/logger');
const { processMessage } = require('./messageProcessor');

let channel, connection;

async function setupRabbitMQ() {
  try {
    connection = await amqp.connect(process.env.RABBITMQ_URL);
    channel = await connection.createChannel();

    // Configurar exchange
    await channel.assertExchange(process.env.RABBITMQ_EXCHANGE, 'topic', {
      durable: true
    });

    // Configurar fila
    const queue = await channel.assertQueue(process.env.RABBITMQ_QUEUE, {
      durable: true
    });

    // Bind da fila com a exchange
    await channel.bindQueue(
      queue.queue,
      process.env.RABBITMQ_EXCHANGE,
      process.env.RABBITMQ_ROUTING_KEY
    );

    // Configurar consumidor
    await channel.consume(queue.queue, async (msg) => {
      if (msg) {
        try {
          logger.info('Received message from queue');
          await processMessage(msg.content.toString());
          channel.ack(msg);
        } catch (error) {
          logger.error('Error processing message:', error);
          channel.nack(msg, false, false);
        }
      }
    });

    logger.info('RabbitMQ connection established');
  } catch (error) {
    logger.error('Error connecting to RabbitMQ:', error);
    throw error;
  }
}

process.on('SIGINT', async () => {
  try {
    await channel?.close();
    await connection?.close();
  } catch (error) {
    logger.error('Error closing RabbitMQ connection:', error);
  }
  process.exit(0);
});

module.exports = { setupRabbitMQ };