#!/bin/bash

output_result_json() {
  result=$1
  result_msg=$2
  result_json="{\"result\":\"${result}\",\"result_msg\":\"${result_msg}\"}"
  echo "OUTPUT_RESULT $result_json OUTPUT_END"
  exit "$result"
}

do_main() {
  # 展示明细文件
  iptables-save |grep -i -E 'drop|reject|forward|return|redirect|random'|grep -i -E -v "KUBE-|DOCKER"
  output_result_json "0" "success"
}

do_main
