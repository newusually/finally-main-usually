# -*- coding: gbk -*-
import base64
import hashlib
import hmac
import json
import os
import time
import urllib
import warnings

import pandas as pd
import requests

from dataprice import DataPrice as data
from mvc import MVC

from userinfo import User

warnings.filterwarnings('ignore')
pd.set_option('mode.chained_assignment', None)


class Databuy:

    def getethinfo(self, minute):

        symbol = "ETH-USDT-SWAP"
        if (not os.path.exists(f'../datas/new_data/' + symbol)):
            try:
                os.makedirs((f'../datas/new_data/' + symbol), exist_ok=True)
            except:
                pass
        else:
            files = f'../datas/old_data/' + symbol + '/' + symbol + '-' + str(minute) + 'min.csv'

            if not os.path.exists(files):
                try:
                    os.makedirs((f'../datas/old_data/' + symbol), exist_ok=True)
                    data.new_symbol_isbuy(minute, symbol)
                except:
                    pass
            else:
                # 使用二进制模式打开文件，并指定忽略编码错误
                with open(files, 'rb') as f:
                    content_bytes = f.read()

                # 尝试解码内容，忽略错误
                content_decoded = content_bytes.decode('utf-8', errors='replace')
                if len(content_decoded) < 500:
                    try:
                        data.new_symbol_isbuy(minute, symbol)
                    except:
                        pass
                else:
                    data.eth_isbuy(minute, symbol)

    def getnext_onedata(self, symbol, minute):

        # 读取csv文件
        df = pd.read_csv(f'../datas/old_data/' + symbol + '/' + symbol + '-' + str(minute) + 'min.csv')

    def getbuyinfo(self, symbollist, minute):

        for symbol in symbollist:

            if (not os.path.exists(f'../datas/new_data/' + symbol)):
                try:
                    os.makedirs((f'../datas/new_data/' + symbol), exist_ok=True)
                except:
                    pass
            else:
                files = f'../datas/old_data/' + symbol + '/' + symbol + '-' + str(minute) + 'min.csv'

                if not os.path.exists(files):

                    os.makedirs((f'../datas/old_data/' + symbol), exist_ok=True)
                    data.new_symbol_isbuy(minute, symbol)

                else:
                    # 使用二进制模式打开文件，并指定忽略编码错误
                    with open(files, 'rb') as f:
                        content_bytes = f.read()

                    # 尝试解码内容，忽略错误
                    content_decoded = content_bytes.decode('utf-8', errors='replace')
                    if len(content_decoded) < 500:

                        data.new_symbol_isbuy(minute, symbol)

                    else:
                        data.eth_isbuy(minute, symbol)
                        isbuy = self.getnext_onedata(symbol, minute)

    def buysell(self, minute):

        api_key, secret_key, passphrase, flag = User.get_userinfo()

        MVC.getuplRatio_instId(api_key, secret_key, passphrase, flag)


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
