#!/usr/bin/env bash

# Define some colors to use for output
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

COMPOSE_FILE=./docker-compose.yaml

# Kill and remove any running containers
cleanup () {
  echo "Clean remaining containers"
  docker-compose -f ${COMPOSE_FILE} stop
  docker-compose -f ${COMPOSE_FILE} rm -f
}

# Catch unexpected failures, do cleanup and output an error message
trap 'cleanup ; printf "${RED}Tests Failed For Unexpected Reasons${NC}\n"'\
  HUP INT QUIT PIPE TERM

# Clean the previous build
cleanup

# Start dependencies
if ! docker-compose -f ${COMPOSE_FILE} run --rm start-dependencies
then
  printf "${RED}Docker Compose Failed${NC}\n"
  exit -1
fi

# Start Dkron
if ! docker-compose -f ${COMPOSE_FILE} run --rm start-services
then
  printf "${RED}Docker Compose Failed${NC}\n"
  exit -1
fi

# Run tests
if ! docker-compose -f ${COMPOSE_FILE} up --abort-on-container-exit integration-tests
then
  printf "${RED}Docker Compose Failed${NC}\n"
  exit -1
fi

# wait for the test service to complete and grab the exit code
TEST_EXIT_CODE=$(docker wait integration-tests)

echo "Script exited with code: $TEST_EXIT_CODE"

# Inspect the output of the test and display respective message
if [ -z ${TEST_EXIT_CODE+x} ] || [ "$TEST_EXIT_CODE" -ne 0 ] ; then
  printf "${RED}Tests Failed${NC} - Exit Code: $TEST_EXIT_CODE\n"
else
  printf "${GREEN}Tests Passed${NC}\n"
fi

# Clean the remaining containers
cleanup

# Exit the script with the same code as the test service code
exit ${TEST_EXIT_CODE}
