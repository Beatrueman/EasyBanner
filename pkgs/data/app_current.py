import logging
from flask import Flask, jsonify, request
import subprocess
import datetime

app = Flask(__name__)

# log
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s',  handlers=[logging.StreamHandler()])

@app.route('/execute', methods=['GET'])
def execute_command():
    # 获取当前时间的小时数，格式化为日志时间格式，例如："25/Aug/2024:15"
    current_time = datetime.datetime.now().strftime("%d/%b/%Y:%H")

    # 构建并执行命令
    command = f"cat /var/log/mirror/access.log | grep '{current_time}' | awk '{{print $1}}' | sort | uniq -c | sort -rn | head -10"
    result = subprocess.run(command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)

    if result.returncode != 0:
        return jsonify({'error': result.stderr.decode('utf-8')}), 500

    # 解析命令输出
    results = []
    lines = result.stdout.decode('utf-8').split('\n')
    for line in lines:
        if line.strip():  # 过滤掉空行
            count, ip = line.split(maxsplit=1)
            results.append({'ip': ip, 'count': int(count)})

    return jsonify(results)


# 对 IP 进行封禁
@app.route('/ban', methods=['POST'])
def ban_ips():
    # 从请求的JSON主体中获取IP列表
    ip_data_list = request.json
    if not ip_data_list:
        logging.warning("Missing IP data in request body!")
        return jsonify({"error": "Missing IP data in request body"}), 400

    results = []
    for ip_data in ip_data_list:
        ip = ip_data.get("ip")
        count = ip_data.get("count", 0)
        if count >= 250:
            try:
                # 对每个IP执行curl命令
                curl_command = f"curl -s '198.18.114.2:8080/update?cidr={ip}&ban_type=0&ban_time=2592000'"
                curl_output = subprocess.check_output(curl_command, shell=True, text=True)
                # 检查curl命令的输出
                if f"Successfully added {ip} to banned list" in curl_output or f'{ip} have been updated' in curl_output:
                    logging.info(f"Succeed to ban {ip}")
                    results.append({"ip": ip, "count": count, "status": "success"})
                else:
                    # 添加日志
                    logging.error(f"Failed to ban {ip}: cur_output: {curl_output}")
                    results.append({"ip": ip, "count": count, "status": "failed"})
            except subprocess.CalledProcessError as e:
                logging.error(f"Failed to ban {ip}: {e}")
                results.append({"ip": ip, "count": count, "status": "failed", "error": str(e)})
        else:
            logging.info(f"Skipped IP: {ip}, count: {count}")
            results.append({"ip": ip, "count": count, "status": "skipped"})

    return jsonify(results)
if __name__ == '__main__':
    app.run(host='172.20.14.2', port=9521)
