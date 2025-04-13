const { S3Client, GetObjectCommand } = require('@aws-sdk/client-s3');
const { SESClient, SendEmailCommand } = require('@aws-sdk/client-ses');
const nodemailer = require('nodemailer');
const logger = require('../utils/logger');

const s3Client = new S3Client({
  region: process.env.AWS_REGION
});

const sesClient = new SESClient({
  region: process.env.AWS_REGION
});

async function downloadFromS3(bucket, key) {
  try {
    const command = new GetObjectCommand({
      Bucket: bucket,
      Key: key
    });

    const response = await s3Client.send(command);
    const chunks = [];

    for await (const chunk of response.Body) {
      chunks.push(chunk);
    }

    return Buffer.concat(chunks);
  } catch (error) {
    logger.error('Error downloading file from S3:', error);
    throw error;
  }
}

async function sendEmailWithSES(recipientEmail, csvContent) {
  try {
    const transporter = nodemailer.createTransport({
      SES: { ses: sesClient, aws: { SendEmailCommand } }
    });

    await transporter.sendMail({
      from: process.env.SES_FROM_EMAIL,
      to: recipientEmail,
      subject: 'Relatório de Tarefas',
      text: 'Segue em anexo seu relatório de tarefas.',
      attachments: [
        {
          filename: 'relatorio-tarefas.csv',
          content: csvContent
        }
      ]
    });

    logger.info('Email sent successfully');
  } catch (error) {
    logger.error('Error sending email:', error);
    throw error;
  }
}

async function processMessage(messageContent) {
  try {
    const message = JSON.parse(messageContent);
    const { s3_bucket, s3_key, recipient_email } = message;

    if (!s3_bucket || !s3_key || !recipient_email) {
      throw new Error('Missing required fields in message');
    }

    logger.info('Downloading file from S3...');
    const csvContent = await downloadFromS3(s3_bucket, s3_key);

    logger.info('Sending email with CSV attachment...');
    await sendEmailWithSES(recipient_email, csvContent);

    logger.info('Message processed successfully');
  } catch (error) {
    logger.error('Error processing message:', error);
    throw error;
  }
}

module.exports = { processMessage };