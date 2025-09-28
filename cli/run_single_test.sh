#!/bin/bash
cd /host-home/erik/Projects/ddx/cli
export DDX_TEST_MODE=1
go test -v -run "TestConfigCommand/get_specific_config_value" .