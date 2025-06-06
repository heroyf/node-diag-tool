#!/bin/bash

output_result_json() {
  result=$1
  result_msg=$2
  result_json="{\"result\":\"${result}\",\"result_msg\":\"${result_msg}\"}"
  echo "OUTPUT_RESULT $result_json OUTPUT_END"
  exit "$result"
}

check_ethtool(){
    # 获取网卡统计信息
    output=$(ip -s -s link show $1)

    # 使用awk解析数据
    awk_result=$(echo "$output" | awk '
    /RX:/ {
        getline
        rx_errors = $3
        rx_dropped = $4
    }
    /RX errors:/ {
        getline
        while ($0 ~ /^[[:space:]]+[0-9]/) {
            if ($1 == "length") crc = $2
            getline
        }
    }
    /TX:/ {
        getline
        tx_errors = $3
        tx_dropped = $4
    }
    END {
        printf "%d %d %d %d %d",crc, rx_errors, rx_dropped, tx_errors, tx_dropped
    }')

    value1=`echo "$awk_result"|awk '{print $1}'`
    if [[ "$value1" -gt 100 ]]; then
       output_result_json "1" "丢包分析检测异常[CRC错包/$value1]"
    fi

    value2=`echo "$awk_result"|awk '{print $2}'`
    if [[ "$value2" -gt 100 ]]; then
       output_result_json "1" "丢包分析检测异常[RX错包/$value2]"
    fi

    value3=`echo "$awk_result"|awk '{print $3}'`
    if [[ "$value3" -gt 100 ]]; then
       output_result_json "1" "丢包分析检测异常[RX丢包/$value3]"
    fi

    value4=`echo "$awk_result"|awk '{print $4}'`
    if [[ "$value4" -gt 100 ]]; then
       output_result_json "1" "丢包分析检测异常[TX错包/$value4]"
    fi

    value5=`echo "$awk_result"|awk '{print $5}'`
    if [[ "$value5" -gt 100 ]]; then
       output_result_json "1" "丢包分析检测异常[TX丢包/$value5]"
    fi
}

do_main() {
  check_ethtool "eth1"

  echo "---eth1 丢包统计---"
  ip -s -s link show eth1

  eth0_count=`ip a|grep -i eth0|wc -l`
  if [[ "$eth0_count" -gt 0 ]]; then
      echo "---eth0 丢包统计---"
      ip -s -s link show eth0

      check_ethtool "eth0"
  fi

  output_result_json "0" "未发现明显丢包"
}

do_main
