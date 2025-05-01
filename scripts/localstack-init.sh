#!/bin/bash

apt install jq -y

# create products table
awslocal dynamodb create-table \
  --table-name products \
  --attribute-definitions AttributeName=id,AttributeType=S \
  --key-schema AttributeName=id,KeyType=HASH \
  --billing-mode PAY_PER_REQUEST \
  --region us-west-2

# create orders table
awslocal dynamodb create-table \
  --table-name orders \
  --attribute-definitions AttributeName=id,AttributeType=S \
  --key-schema AttributeName=id,KeyType=HASH \
  --billing-mode PAY_PER_REQUEST \
  --region us-west-2

# create SNS topic
awslocal sns create-topic --name orders-topic

# create SQS
awslocal sqs create-queue --queue-name products-queue

# get the queue url
queue_url=$(awslocal sqs get-queue-url --queue-name products-queue | jq -r '.QueueUrl')

#get the queue arn
queue_arn=$(awslocal sqs get-queue-attributes \
  --queue-url $queue_url \
  --attribute-name QueueArn | jq -r '.Attributes.QueueArn')

# get the topic arn
topic_arn=$(awslocal sns list-topics | jq -r '.Topics[0].TopicArn')
# subscribe the queue to the topic
awslocal sns subscribe \
  --topic-arn $topic_arn \
  --protocol sqs \
  --notification-endpoint $queue_arn