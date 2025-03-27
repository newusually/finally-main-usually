# -*- coding: gbk -*-
import datetime
import numpy as np
import talib
import pandas as pd

def parse_txt_content(content):
    lines = content.strip().split('\n')
    data = {
        'time': [],
        'eth_close': [],
        'filtered_eth_close': [],
        'a': [],
        'b': [],
        'e': [],
        'a-b': [],
        'b-e': [],
        'macd': [],
        'signal': [],
        'hist': [],
        'buy_label': [],
        'next_5_closes': [],
    }

    for line in lines:
        items = line.split(',')
        time_str = items[1].split('>>')[1]
        time_obj = datetime.datetime.strptime(time_str, '%Y-%m-%d %H:%M:%S')
        data['time'].append(time_obj)
        close = float(items[7].split('>>')[1])
        data['eth_close'].append(close)
        data['a'].append(float(items[2].split('>>')[1]))
        data['b'].append(float(items[3].split('>>')[1]))
        data['e'].append(float(items[4].split('>>')[1]))
        data['a-b'].append(float(items[5].split('>>')[1]))
        data['b-e'].append(float(items[6].split('>>')[1]))

    macd, signal, hist = talib.MACD(np.array(data['eth_close']), fastperiod=12, slowperiod=26, signalperiod=9)

    for i in range(len(lines)):
        data['macd'].append(macd[i])
        data['signal'].append(signal[i])
        data['hist'].append(hist[i])

        if i > 0 and macd[i] < -5 and macd[i] > macd[i-1] and data['a'][i] < 0.03 :
            data['buy_label'].append(1)
        else:
            data['buy_label'].append(0)

        if i + 5 < len(lines):
            next_5_closes = ','.join([str(data['eth_close'][j]) for j in range(i+1, i+6)])
        else:
            next_5_closes = ''
        data['next_5_closes'].append(next_5_closes)

        # Filter the close values with 1.0025 increase
        if i >= 20 and data['eth_close'][i] >= data['eth_close'][i-20] * 1.0025:
            data['filtered_eth_close'].append(data['eth_close'][i])
        else:
            data['filtered_eth_close'].append(None)

    df=pd.DataFrame(data)
    df.to_csv("ethdata.csv",index=0)
    # 计算相关系数
    correlation_filtered_eth_close_a = df['filtered_eth_close'].corr(df['a'])
    correlation_filtered_eth_close_b = df['filtered_eth_close'].corr(df['b'])
    correlation_filtered_eth_close_e = df['filtered_eth_close'].corr(df['e'])
    correlation_filtered_eth_close_a_b = df['filtered_eth_close'].corr(df['a-b'])
    correlation_filtered_eth_close_b_e = df['filtered_eth_close'].corr(df['b-e'])
    correlation_filtered_eth_close_macd = df['filtered_eth_close'].corr(df['macd'])

    # 打印相关系数
    print("Correlation between filtered_eth_close and a:", correlation_filtered_eth_close_a)
    print("Correlation between filtered_eth_close and b:", correlation_filtered_eth_close_b)
    print("Correlation between filtered_eth_close and e:", correlation_filtered_eth_close_e)
    print("Correlation between filtered_eth_close and a-b:", correlation_filtered_eth_close_a_b)
    print("Correlation between filtered_eth_close and b-e:", correlation_filtered_eth_close_b_e)
    print("Correlation between filtered_eth_close and macd:", correlation_filtered_eth_close_macd)

def read_data(file_path):
    with open(file_path, 'r') as file:
        content = file.read()
    parse_txt_content(content)

if __name__ == '__main__':
    file_path='15min_buylog.txt'
    read_data(file_path)