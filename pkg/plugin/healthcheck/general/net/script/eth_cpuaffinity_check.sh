#!/bin/bash
# 中断CPU亲和性检测脚本
# 功能：解析/proc/interrupts中的中断号，并显示对应的CPU核心分布

output_result_json() {
  result=$1
  result_msg=$2
  result_json="{\"result\":\"${result}\",\"result_msg\":\"${result_msg}\"}"
  echo "OUTPUT_RESULT $result_json OUTPUT_END"
  exit "$result"
}

do_main() {
  irq_list=$(awk 'NF>=4 {print $1 "\t" $(NF) "\t"  $(NF-1) "\t" $(NF-2) "\t" $(NF-3) "\t" $(NF-4)}' /proc/interrupts|grep -i -E 'input|eth')

  #echo "中断号 | CPU掩码(HEX) | CPU核心分布"
  #echo "----------------------------------------"
  #printf "%-8s %-12s %-20s\n" "IRQ" "SMP_Affinity" "CPU Cores"

  # 用于存储所有唯一的CPU核心编号
  declare -A unique_cores
  IFS=$'\n'
  for line in ${irq_list};
  do
      irq=$(echo $line|awk '{print $1}'|awk -F ':' '{print $1}')
      affinity_file="/proc/irq/${irq}/smp_affinity"

      # 读取亲和性掩码（十六进制）
      hex_mask=$(tr -d ' ' < "$affinity_file" | tr '[:lower:]' '[:upper:]'|awk -F ',' '{print $1}')

      # 十六进制转二进制（补前导零至8位，支持最多64核）
      bin_mask=$(echo "obase=2; ibase=16; $hex_mask" | bc | awk '{printf "%08d\n", $0}')

      # 计算绑定的CPU核心列表
      cpu_cores=()
      for ((i=0; i<${#bin_mask}; i++)); do
          if [ "${bin_mask:$i:1}" -eq 1 ]; then
              cpu_cores+=("$i")
              # 使用关联数组记录唯一的CPU核心编号
              unique_cores[$i]=1
          fi
      done

      # 格式化输出
      #printf "%-8s %-12s %-20s\n" "$irq" "$hex_mask" "${cpu_cores[@]:-未绑定}"
  done

  # 计算唯一的CPU核心数量
  core_count=${#unique_cores[@]}

  # 判断绑核数量是否超过4
  if [ $core_count -gt 4 ]; then
      output_result_json "0" "网卡多队列绑核检测通过：当前绑核数为$core_count"
  else
      output_result_json "1" "网卡多队列绑核检测异常：当前绑核数为${core_count}."
  fi
}

do_main