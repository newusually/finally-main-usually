# -*- coding: gbk -*-
import base64
import hashlib
import hmac
import json
import time
import urllib
from datetime import datetime

import numpy as np
import pandas as pd
import requests as r


class DataPrice:
    def new_symbol_isbuy(minute, symbol):

        for i in range(3):
            try:

                if (minute == '1'):
                    minute = '1m'
                if (minute == '3'):
                    minute = '3m'
                if (minute == '5'):
                    minute = '5m'
                if (minute == '15'):
                    minute = '15m'

                t = time.time()

                # print (t)                       #原始时间数据
                # print (int(t))                  #秒级时间戳
                # print (int(round(t * 1000)))    #毫秒级时间戳
                # print (int(round(t * 1000000))) #微秒级时间戳
                tt = str((int(t * 1000)))
                ttt = str((int(round(t * 1000000))))

                # time.sleep(int(minute)/10)

                # ===获取close数据

                headers = {
                    'authority': 'www.okx.com',
                    'timeout': '1',
                    'x-cdn': 'https://static.okx.com',
                    'devid': '6ec23520-a48b-41f1-b35e-5dea795c61b8',
                    'accept-language': 'zh-CN',
                    'user-agent': 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.87 Safari/537.36 SE 2.X MetaSr 1.0',
                    'accept': 'application/json',
                    'x-utc': '8',
                    'sec-fetch-dest': 'empty',
                    'app-type': 'web',
                    'sec-fetch-site': 'same-origin',
                    'sec-fetch-mode': 'cors',
                    'referer': 'https://www.okx.com/trade-swap/' + symbol,
                    'cookie': 'locale=zh_CN; defaultLocale=zh_CN; _gcl_au=1.1.1514314517.' + str(
                        tt) + '; _ga=GA1.2.1025788009.' + str(tt) + '; _gid=GA1.2.1077289716.' + str(
                        tt) + '; amp_56bf9d=zTmKdiXyRK-5EUgHM2Qg_x...1fp5jebfd.1fp5jgo7d.2.0.2',
                }

                params = (
                    ('instId', symbol),
                    ('bar', minute),
                    ('after', ''),
                    ('limit', '1500'),
                    ('t', str(ttt)),
                )

                response = r.get('https://www.okx.com/priapi/v5/market/candles', headers=headers, params=params)

                if response.cookies.get_dict():  # 保持cookie有效
                    s = r.session()
                    c = r.cookies.RequestsCookieJar()  # 定义一个cookie对象
                    c.set('cookie-name', 'cookie-value')  # 增加cookie的值
                    s.cookies.update(c)  # 更新s的cookie
                    s.get(url='https://www.okx.com/priapi/v5/market/candles?instId=' + symbol + '&bar=' + str(
                        minute) + '&after=&limit=1500&t=' + tt, headers=headers)
                # print(eval(json.dumps(response.json()))['data'])
                new_df = pd.DataFrame(eval(json.dumps(response.json()))['data'])

                response.close()
                time.sleep(1)
                # print(new_df)
                df = pd.DataFrame()
                df['date'] = new_df[0]
                df['open'] = new_df[1]
                df['high'] = new_df[2]
                df['low'] = new_df[3]
                df['close'] = new_df[4]
                df['vol'] = new_df[5]

                # new_df.columns = ['date', 'open', 'high', 'low', 'close', 'vol', 'p', 'pp']
                datelist = []

                for timestamp in df['date']:
                    date = datetime.fromtimestamp(int(timestamp) / 1000)
                    date = date.strftime('%Y-%m-%d %H:%M:%S')
                    datelist.append(date)
                df['date'] = datelist
                # df['date'] = pd.to_datetime(df['date'], format='mixed')
                df['vol'] = df['vol'].astype('float')
                df['close'] = df['close'].astype('float')
                # print(new_df)
                df.sort_values(by=['date'], axis=0, ascending=True, inplace=True)

                if (minute == '1m'):
                    minute = '1'
                if (minute == '3m'):
                    minute = '3'
                if (minute == '5m'):
                    minute = '5'
                if (minute == '15m'):
                    minute = '15'
                df.to_csv(
                    f'../datas/new_data/' + symbol + '/' + symbol + '-' + str(minute) + 'min.csv', index=False)
                df.to_csv(
                    f'../datas/old_data/' + symbol + '/' + symbol + '-' + str(minute) + 'min.csv', index=False)
                return df

            except:
                time.sleep(0.5)
                continue

    def eth_isbuy(minute, symbol):
        close = 0
        for i in range(1):
            try:

                if (minute == '1'):
                    minute = '1m'
                if (minute == '3'):
                    minute = '3m'
                if (minute == '5'):
                    minute = '5m'
                if (minute == '15'):
                    minute = '15m'

                t = time.time()

                # print (t)                       #原始时间数据
                # print (int(t))                  #秒级时间戳
                # print (int(round(t * 1000)))    #毫秒级时间戳
                # print (int(round(t * 1000000))) #微秒级时间戳
                tt = str((int(t * 1000)))
                ttt = str((int(round(t * 1000000))))

                # time.sleep(int(minute) / 10)

                # ===获取close数据

                headers = {
                    'authority': 'www.okx.com',
                    'timeout': '1',
                    'x-cdn': 'https://static.okx.com',
                    'devid': '6ec23520-a48b-41f1-b35e-5dea795c61b8',
                    'accept-language': 'zh-CN',
                    'user-agent': 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.87 Safari/537.36 SE 2.X MetaSr 1.0',
                    'accept': 'application/json',
                    'x-utc': '8',
                    'sec-fetch-dest': 'empty',
                    'app-type': 'web',
                    'sec-fetch-site': 'same-origin',
                    'sec-fetch-mode': 'cors',
                    'referer': 'https://www.okx.com/trade-swap/' + symbol,
                    'cookie': 'locale=zh_CN; defaultLocale=zh_CN; _gcl_au=1.1.1514314517.' + str(
                        tt) + '; _ga=GA1.2.1025788009.' + str(tt) + '; _gid=GA1.2.1077289716.' + str(
                        tt) + '; amp_56bf9d=zTmKdiXyRK-5EUgHM2Qg_x...1fp5jebfd.1fp5jgo7d.2.0.2',
                }

                params = (
                    ('instId', symbol),
                    ('bar', minute),
                    ('after', ''),
                    ('limit', '1500'),
                    ('t', str(ttt)),
                )

                response = r.get('https://www.okx.com/priapi/v5/market/candles', headers=headers, params=params)

                if response.cookies.get_dict():  # 保持cookie有效
                    s = r.session()
                    c = r.cookies.RequestsCookieJar()  # 定义一个cookie对象
                    c.set('cookie-name', 'cookie-value')  # 增加cookie的值
                    s.cookies.update(c)  # 更新s的cookie
                    s.get(url='https://www.okx.com/priapi/v5/market/candles?instId=' + symbol + '&bar=' + str(
                        minute) + '&after=&limit=1500&t=' + tt, headers=headers)
                # print(eval(json.dumps(response.json()))['data'])
                new_df = pd.DataFrame(eval(json.dumps(response.json()))['data'])

                response.close()
                time.sleep(1)
                # print(new_df)
                df = pd.DataFrame()
                df['date'] = new_df[0]
                df['open'] = new_df[1]
                df['high'] = new_df[2]
                df['low'] = new_df[3]
                df['close'] = new_df[4]
                df['vol'] = new_df[5]

                # new_df.columns = ['date', 'open', 'high', 'low', 'close', 'vol', 'p', 'pp']
                datelist = []

                for timestamp in df['date']:
                    date = datetime.fromtimestamp(int(timestamp) / 1000)
                    date = date.strftime('%Y-%m-%d %H:%M:%S')
                    datelist.append(date)
                df['date'] = datelist
                # df['date'] = pd.to_datetime(df['date'], format='mixed')
                df['vol'] = df['vol'].astype('float')
                df['close'] = df['close'].astype('float')
                # print(new_df)
                df.sort_values(by=['date'], axis=0, ascending=True, inplace=True)

                if (minute == '1m'):
                    minute = '1'
                if (minute == '3m'):
                    minute = '3'
                if (minute == '5m'):
                    minute = '5'
                if (minute == '15m'):
                    minute = '15'

                df.to_csv(
                    f'../datas/new_data/' + symbol + '/' + symbol + '-' + str(minute) + 'min.csv', index=False)

                # 保存数据 变换数据格式
                old_df = pd.read_csv(f'../datas/old_data/' + symbol + '/' + symbol + '-' + str(minute) + 'min.csv')[
                         :-100]
                old_df[['date', 'close', 'open', 'high', 'low', 'vol']].to_csv(
                    f'../datas/old_data/' + symbol + '/' + symbol + '-' + str(minute) + 'min.csv', index=0)
                old_df = pd.read_csv(f'../datas/old_data/' + symbol + '/' + symbol + '-' + str(minute) + 'min.csv')

                old_df['date'] = pd.to_datetime(old_df['date'])
                # old_df['date'] = pd.to_datetime(df['date'], format='mixed')
                new_df = pd.read_csv(f'../datas/new_data/' + symbol + '/' + symbol + '-' + str(minute) + 'min.csv')

                df = pd.DataFrame()

                # 确保在合并新旧数据后立即对数据进行排序，并且重置索引
                df = pd.concat([old_df, new_df], axis=0)
                df['date'] = pd.to_datetime(df['date'], format='%Y-%m-%d %H:%M:%S')
                df.sort_values(by=['date'], axis=0, ascending=True, inplace=True)
                df.reset_index(drop=True, inplace=True)  # 重置索引

                # 转换数据类型
                df['vol'] = df['vol'].astype('float')
                df['open'] = df['open'].astype('float')
                df['close'] = df['close'].astype('float')
                df['high'] = df['high'].astype('float')
                df['low'] = df['low'].astype('float')

                df.drop_duplicates(subset=['date'], keep='first', inplace=True)

                # 保存数据之前再次确保数据是按日期正序排列的
                df = df[['date', 'close', 'open', 'high', 'low', 'vol']]
                df.to_csv(f'../datas/old_data/{symbol}/{symbol}-{minute}min.csv', index=False)
                # df = getfulldata(df).dropna()
                df.to_csv(
                    f'../datas/old_data/' + symbol + '/' + symbol + '-' + str(minute) + 'min.csv', index=0)
                if symbol == 'ETH-USDT-SWAP':
                    # date=df['date'].values[-1]
                    close = df["close"].values[-1]
                    # print(symbol,minute,close)
                    # print("================================================")
                    if len(close) < 3:
                        SendDingding.sender(
                            "close null!!!--->>>" + "，\n我们是守护者，也是一群时刻对抗危险和疯狂的可怜虫！！！")

                break
            except:
                time.sleep(0.5)
                continue

        return str(close)


# 发钉钉的类先声明
class SendDingding:
    def sender(txt):
        headers = {
            'Content-Type': 'application/json'
        }
        timestamp = str(round(time.time() * 1000))
        secret = "SEC050a3b2c9e5d8d0c777bbdd61270676a8bdad3608b36a086d70e95b712ad2db0"
        secret_enc = secret.encode('utf-8')
        string_to_sign = '{}\n{}'.format(timestamp, secret)
        string_to_sign_enc = string_to_sign.encode('utf-8')
        hmac_code = hmac.new(secret_enc, string_to_sign_enc, digestmod=hashlib.sha256).digest()
        sign = urllib.parse.quote_plus(base64.b64encode(hmac_code))
        today = time.strftime("%Y-%m-%d %H:%M:%S", time.localtime())

        sendtexts = "本地时间： " + today + "--->>>" + txt + "，\n" + "，\n我们是守护者，也是一群时刻对抗危险和疯狂的可怜虫！！！"

        params = {
            'sign': sign,

            'timestamp': timestamp
        }
        text_data = {
            "msgtype": "text",
            "text": {
                "content": sendtexts
            }
        }

        roboturl = 'https://oapi.dingtalk.com/robot/send?access_token=f8195c9e4ad6da4427d67e80dffed5d07ecaca1d1e79462fb5c0a9c6b12e90f2'
        r = requests.post(roboturl, data=json.dumps(text_data), params=params, headers=headers)
