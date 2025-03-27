# -*- coding: gbk -*-

import base64
import hashlib
import hmac
import json
import os
import time
import urllib.parse
from datetime import datetime

import pandas as pd
import requests
import requests as r

import okx.Account_api as Account
import okx.Market_api as Market
import okx.Trade_api as Trade
from userinfo import User


class MVC:

    def sender(sendtexts):
        try:
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

            params = {
                'sign': sign,

                'timestamp': timestamp
            }
            text_data = {
                "msgtype": "text",
                "text": {
                    "content": str(
                        datetime.now()) + '--->>>我们是守护者，也是一群时刻对抗危险和疯狂的可怜虫 ！^_^     -->> ' + sendtexts
                }
            }
            roboturl = 'https://oapi.dingtalk.com/robot/send?access_token=f8195c9e4ad6da4427d67e80dffed5d07ecaca1d1e79462fb5c0a9c6b12e90f2'
            r.post(roboturl, data=json.dumps(text_data), params=params, headers=headers)
        except:
            pass

    def getsymbollist():
        t = time.time()
        # 原始时间数据
        # print (int(t))                  #秒级时间戳
        # print (int(round(t * 1000)))    #毫秒级时间戳
        # print (int(round(t * 1000000))) #微秒级时间戳
        tt = str((int(t * 1000)))
        ttt = str((int(round(t * 1000000))))
        headers = {
            'authority': 'www.okx.com',
            'timeout': '10000',
            'x-cdn': 'https://static.okx.com',
            'devid': 'c5ccc0b9-af31-436e-9345-3f4c9a7a65fc',
            'accept-language': 'zh-CN',
            'user-agent': 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.87 Safari/537.36 SE 2.X MetaSr 1.0',
            'accept': 'application/json',
            'x-utc': '8',
            'sec-fetch-dest': 'empty',
            'app-type': 'web',
            'sec-fetch-site': 'same-origin',
            'sec-fetch-mode': 'cors',
            'referer': 'https://www.okx.com/trade-spot/btc-usdt',
            'cookie': 'locale=zh_CN; defaultLocale=zh_CN; _gcl_au=1.1.2606549.' + str(
                tt) + '; _ga=GA1.2.2137624954.' + str(tt) + '; _gid=GA1.2.1766984091.' + str(
                tt) + '; _gat_UA-35324627-3=1; amp_56bf9d=mUJNvZMqgezu0vXrzFkdTp...1fttm7de6.1fttm7de6.0.0.0',
        }

        params = (
            ('t', str(tt)),
            ('instType', 'SWAP'),
        )

        response = requests.get('https://www.okx.com/priapi/v5/public/simpleProduct', headers=headers, params=params)

        # Note: original query string below. It seems impossible to parse and
        # reproduce query strings 100% accurately so the one below is given
        # in case the reproduced version is not "correct".
        # response = requests.get('https://www.okx.com/priapi/v5/market/tickers?t='+str(tt)+'^&instType=SWAP', headers=headers)
        symbollist = []
        symbols = pd.DataFrame(eval(json.dumps(response.json()))['data'])['instId']
        for symbol in symbols:
            symbol = symbol.split('-')[0].split('-')[0] + '-USDT-SWAP'
            symbollist.append(symbol)
        return list(set(symbollist))

    # 查询最新价格
    def getlastprice(api_key, secret_key, passphrase, flag, symbol):

        # market api
        marketAPI = Market.MarketAPI(api_key, secret_key, passphrase, False, flag)

        # 获取单个产品行情信息  Get Ticker
        result = marketAPI.get_ticker(symbol)
        print(eval(json.dumps(result['data'][0])))

        return eval(json.dumps(result['data'][0]))['last']

    def sellall(api_key, secret_key, passphrase, flag):
        # account api
        accountAPI = Account.AccountAPI(api_key, secret_key, passphrase, False, flag)
        result = accountAPI.get_position_risk('SWAP')
        datas = result['data'][0]['posData']
        # 获取持仓的交易品种['instId']

        if len(datas) > 0:
            # print(datas)
            # 使用循环遍历 posData 列表中的每个元素
            for item in datas:
                try:
                    time.sleep(1)
                    # 打印每个 instId
                    symbol = item['instId']
                    # 将 instId 添加到 symbollist 列表中

                    # 查看持仓信息  Get Positions
                    result = accountAPI.get_positions('SWAP', symbol)

                    # 未实现收益率
                    uplRatio = float(result['data'][0]['uplRatio'])

                    if uplRatio > 0.35 or uplRatio < -30:
                        print("symbol--->>>", symbol, "未实现收益率--->>>", uplRatio)
                        tradeAPI = Trade.TradeAPI(api_key, secret_key, passphrase, False, flag)

                        # 市价仓位全平  Close Positions
                        result = tradeAPI.close_positions(symbol, 'cross', 'long', '')
                except:
                    pass

    # 查询是否抄底 然后 读取买入日志更新最新10行数据保存
    def getuplRatio_instId(api_key, secret_key, passphrase, flag):

        # 读取文件
        # df = pd.read_csv("..\\datas\\log\\buylog.txt", header=None)

        # 获取最后10行
        # last_10_lines = df.tail(10)

        # 保存到原文件，替换源文件
        # last_10_lines.to_csv("..\\datas\\log\\buylog.txt", header=False, index=False)

        # account api
        accountAPI = Account.AccountAPI(api_key, secret_key, passphrase, False, flag)
        result = accountAPI.get_position_risk('SWAP')
        datas = result['data'][0]['posData']

        # 获取持仓的交易品种['instId']

        # 定义要存放文件的路径
        folder_path = '../datas/uplRatio/log'
        file_path = os.path.join(folder_path, 'uplRatio.txt')

        # 检查路径是否存在，不存在则创建
        if not os.path.exists(folder_path):
            os.makedirs(folder_path)

        # 现在路径存在，打开文件并清空之前内容
        with open(file_path, 'w', encoding='utf-8') as file:
            # 将日志信息写入文件
            file.write("")

        if len(datas) > 0:
            # print(datas)
            # 使用循环遍历 posData 列表中的每个元素
            for item in datas:
                try:
                    time.sleep(1)
                    # 打印每个 instId
                    symbol = item['instId']
                    # 将 instId 添加到 symbollist 列表中

                    # 查看持仓信息  Get Positions
                    result = accountAPI.get_positions('SWAP', symbol)
                    posSide = result['data'][0]['posSide']
                    # 未实现收益率
                    uplRatio = float(result['data'][0]['uplRatio'])

                    print(uplRatio)
                    if uplRatio < -0.5:

                        if -5 < uplRatio < -3:

                            MVC.buyorders(api_key, secret_key, passphrase, flag, symbol, "imr")

                        else:

                            MVC.buyorders(api_key, secret_key, passphrase, flag, symbol, "low")

                    notionalUsd = float(result['data'][0]['notionalUsd'])
                    # 保证金
                    imr = float(result['data'][0]['imr'])
                    # 量
                    pos = float(result['data'][0]['pos'])
                    # 单个量 保证金值
                    onlyimr = imr / pos
                    # 查看持仓信息  Get Positions
                    result = accountAPI.get_positions('SWAP', symbol)
                    # 开仓均价
                    avgPx = float(result['data'][0]['avgPx'])
                    # 现价
                    last = float(result['data'][0]['last'])
                    # 杠杆倍数
                    lever = float(result['data'][0]['lever'])

                    log = ("\nsymbol--->>>" + symbol + ",未实现收益率--->>>" + "{:.5f}".format(uplRatio * 100) + "%" +
                           ",现价--->>>" + "{:.5f}".format(last) + ",开仓均价--->>>" + "{:.5f}".format(
                                avgPx) + ",保证金--->>>" + "{:.5f}".format(imr) + ",杠杆倍数--->>>" + str(
                                lever) + ",总计亏损金额--->>>" + "{:.5f}".format(
                                imr * uplRatio))

                    # 现在路径存在，打开文件并在尾部追加内容，指定编码为UTF-8
                    with open(file_path, 'a', encoding='utf-8') as file:
                        # 将日志信息写入文件
                        file.write(log)


                except:
                    pass

    # 获取实时账户资金信息 每分钟查询一次
    def getcashbal(api_key, secret_key, passphrase, flag):
        # 定义要存放文件的路径
        folder_path = '../datas/uplRatio/log'
        file_path = os.path.join(folder_path, 'cashbal.txt')

        # 检查路径是否存在，不存在则创建
        if not os.path.exists(folder_path):
            os.makedirs(folder_path)

        # 现在路径存在，打开文件并清空之前内容
        with open(file_path, 'w', encoding='utf-8') as file:
            # 将日志信息写入文件
            file.write("")

        # account api
        accountAPI = Account.AccountAPI(api_key, secret_key, passphrase, False, flag)

        # 查看账户余额  Get Balance
        result = accountAPI.get_account('USDT')['data'][0]["details"][0]
        swap = accountAPI.get_position_risk('SWAP')
        posData = 0 if not swap['data'][0]['posData'] else swap['data'][0]['posData']
        posData_length = str(len(posData)) if isinstance(posData, list) else str(posData)
        # 检查每个值，如果为空或者为0，就打印0，否则打印实际的值
        disEq = 0 if not result["disEq"] else float(result["disEq"])
        upl = 0 if not result["upl"] else float(result["upl"])
        cashBal = 0 if not result["cashBal"] else float(result["cashBal"])
        mgnRatio = 0 if not result["mgnRatio"] else float(result["mgnRatio"])
        availBal = 0 if not result["availBal"] else float(result["availBal"])
        frozenBal = 0 if not result["frozenBal"] else float(result["frozenBal"])

        log = ("\n合约订单数量----->>>" + str(posData_length) + "个,  " +
               "币种折算权益----->>>" + "{:.2f}".format(disEq) + "＄,  " +
               "\n实际未结算盈亏总额：--->>>" + "{:.2f}".format(upl) + "＄,  " +
               "USDT币种余额----->>>" + "{:.2f}".format(cashBal) + "＄," +
               "\n保证金率----->>>" + "{:.2f}".format(mgnRatio * 100) + "%,  " +
               "可用余额----->>>" + "{:.2f}".format(availBal) + "＄,  " +
               "\n,币种占用金额--->>>" + "{:.2f}".format(frozenBal) + "＄" + "\n\n")

        # 现在路径存在，打开文件并在尾部追加内容，指定编码为UTF-8
        with open(file_path, 'a', encoding='utf-8') as file:
            # 将日志信息写入文件
            file.write(log)

        SendDingding.sender(log)

    # 资金历史记录 每分钟记录一次
    def getcashhistory(api_key, secret_key, passphrase, flag):

        # account api
        accountAPI = Account.AccountAPI(api_key, secret_key, passphrase, False, flag)
        # 查看账户持仓风险 GET Position_risk
        # result = accountAPI.get_position_risk('SWAP')
        # 查看账户余额  Get Balance
        result = accountAPI.get_account('USDT')['data'][0]["details"][0]
        swap = accountAPI.get_position_risk('SWAP')

        posData = 0 if not swap['data'][0]['posData'] else swap['data'][0]['posData']
        posData_length = str(len(posData)) if isinstance(posData, list) else str(posData)
        today = time.strftime("%Y-%m-%d %H:%M:%S", time.localtime())

        # 检查每个值，如果为空或者为0，就打印0，否则打印实际的值
        disEq = 0 if not result["disEq"] else float(result["disEq"])
        upl = 0 if not result["upl"] else float(result["upl"])
        cashBal = 0 if not result["cashBal"] else float(result["cashBal"])
        mgnRatio = 0 if not result["mgnRatio"] else float(result["mgnRatio"])
        frozenBal = 0 if not result["frozenBal"] else float(result["frozenBal"])

        # 美金层面币种折算权益
        with open("../datas/uplRatio/log/disEq_history.txt", 'a', encoding='utf-8') as file:
            # 将日志信息写入文件
            file.write("\n" + today + "," + "{:.2f}".format(disEq))

        # 实际未结算盈亏总额
        with open("../datas/uplRatio/log/upl_history.txt", 'a', encoding='utf-8') as file:
            # 将日志信息写入文件
            file.write("\n" + today + "," + "{:.2f}".format(upl))

        # USDT币种余额
        with open("../datas/uplRatio/log/cashBal_history.txt", 'a', encoding='utf-8') as file:
            # 将日志信息写入文件
            file.write("\n" + today + "," + "{:.2f}".format(float(cashBal)))

        # 合约订单数量
        with open("../datas/uplRatio/log/posdatacount_history.txt", 'a', encoding='utf-8') as file:
            # 将日志信息写入文件
            file.write("\n" + today + "," + str(posData_length))

        # USDT保证金金额
        with open("../datas/uplRatio/log/frozenBal_history.txt", 'a', encoding='utf-8') as file:
            # 将日志信息写入文件
            file.write("\n" + today + "," + "{:.2f}".format(frozenBal))

        # USDT保证金率
        with open("../datas/uplRatio/log/mgnRatio_history.txt", 'a', encoding='utf-8') as file:
            # 将日志信息写入文件
            file.write("\n" + today + "," + "{:.2f}".format(mgnRatio * 100) + "%")

    def ordersell(api_key, secret_key, passphrase, flag, symbol, minute):
        # account api
        accountAPI = Account.AccountAPI(api_key, secret_key, passphrase, False, flag)
        result = accountAPI.get_position_risk('SWAP')
        posData = result['data'][0]['posData']

        if (len(posData) < 50 and minute != "low") or minute == "low" or minute == "imr":
            MVC.sellorders(api_key, secret_key, passphrase, flag, symbol, minute)

    def sellorders(api_key, secret_key, passphrase, flag, symbol, minute):
        sr, dollar, dollar_eth = User.get_user_sr()
        sr1 = str(sr)
        r = 0.9965

        accountAPI = Account.AccountAPI(api_key, secret_key, passphrase, False, flag)

        swap = accountAPI.get_position_risk('SWAP')

        # 设置杠杆倍数  Set Leverage
        accountAPI.set_leverage(instId=symbol, lever='50', mgnMode='cross')

        # 第一次买入
        tradeAPI = Trade.TradeAPI(api_key, secret_key, passphrase, False, flag)
        tradeAPI.place_order(instId=symbol, tdMode='cross', side='sell', posSide='short',
                             ordType='market', sz=sr1)

        # 设置杠杆倍数  Set Leverage
        accountAPI.set_leverage(instId=symbol, lever='50', mgnMode='cross')

        # account api
        accountAPI = Account.AccountAPI(api_key, secret_key, passphrase, False, flag)
        # 查看持仓信息  Get Positions
        result = accountAPI.get_positions('SWAP', symbol)

        time.sleep(5)

        # =====================================================================
        # print(dollar)
        # 持仓量
        if len(result['data']) > 0 and len(result['data'][0]) > 0 and len(
                result['data'][0]['notionalUsd']) > 0:
            notionalUsd = float(result['data'][0]['notionalUsd'])
            # 保证金
            imr = float(result['data'][0]['imr'])
            # 量
            pos = float(result['data'][0]['pos'])
            # 单个量 保证金值
            onlyimr = imr / pos
            # 每一美金值多少量
            onlyorder = int(float(dollar) / onlyimr)

            # 查看持仓信息  Get Positions
            result = accountAPI.get_positions('SWAP', symbol)

            if "low" == minute:
                onlyorder = int(onlyorder / 2)
            if "imr" == minute:
                onlyorder = int(onlyorder * 1.5)

            # 第2次买入
            tradeAPI = Trade.TradeAPI(api_key, secret_key, passphrase, False, flag)
            tradeAPI.place_order(instId=symbol, tdMode='cross', side='sell', posSide='short',
                                 ordType='market', sz=str(onlyorder))

            # 设置杠杆倍数  Set Leverage
            accountAPI.set_leverage(instId=symbol, lever='50', mgnMode='cross')

            time.sleep(1)

            # 查看持仓信息  Get Positions
            result = accountAPI.get_positions('SWAP', symbol)

            # 开仓均价
            avgPx = float(result['data'][0]['avgPx'])

            # 策略委托下单  Place Algo Order
            tradeAPI.place_algo_order(symbol, 'cross', 'buy', ordType='conditional',
                                      sz=sr1, posSide='long', tpTriggerPx=str(float(avgPx) * r),
                                      tpOrdPx=str(float(avgPx) * r))

            # 策略委托下单  Place Algo Order
            result = tradeAPI.place_algo_order(symbol, 'cross', 'buy', ordType='conditional',
                                               sz=str(onlyorder), posSide='long', tpTriggerPx=str(float(avgPx) * r),
                                               tpOrdPx=str(float(avgPx) * r))
            time.sleep(1)
            # 查看持仓信息  Get Positions
            result = accountAPI.get_positions('SWAP', symbol)

            # 开仓均价
            avgPx = float(result['data'][0]['avgPx'])

            # 策略委托下单  Place Algo Order
            result = tradeAPI.place_algo_order(symbol, 'cross', 'buy', ordType='conditional',
                                               sz=str(onlyorder), posSide='long', tpTriggerPx=str(float(avgPx) * r),
                                               tpOrdPx=str(float(avgPx) * r))

            time.sleep(1)

            # 查看持仓信息  Get Positions
            result = accountAPI.get_positions('SWAP', symbol)

            # 开仓均价
            avgPx = float(result['data'][0]['avgPx'])

            # 策略委托下单  Place Algo Order
            result = tradeAPI.place_algo_order(symbol, 'cross', 'buy', ordType='conditional',
                                               sz=str(onlyorder), posSide='long', tpTriggerPx=str(float(avgPx) * r),
                                               tpOrdPx=str(float(avgPx) * r))

    def orderbuy(api_key, secret_key, passphrase, flag, symbol, minute):
        # account api
        accountAPI = Account.AccountAPI(api_key, secret_key, passphrase, False, flag)
        result = accountAPI.get_position_risk('SWAP')
        posData = result['data'][0]['posData']

        if (len(posData) < 15 and minute != "low") or minute == "low" or minute == "imr":
            MVC.buyorders(api_key, secret_key, passphrase, flag, symbol, minute)

    def buyorders(api_key, secret_key, passphrase, flag, symbol, minute):
        sr, dollar, dollar_eth = User.get_user_sr()
        sr1 = str(sr)
        r = 1.0035

        accountAPI = Account.AccountAPI(api_key, secret_key, passphrase, False, flag)

        swap = accountAPI.get_position_risk('SWAP')

        # 设置杠杆倍数  Set Leverage
        accountAPI.set_leverage(instId=symbol, lever='50', mgnMode='cross')

        # 第一次买入
        tradeAPI = Trade.TradeAPI(api_key, secret_key, passphrase, False, flag)
        tradeAPI.place_order(instId=symbol, tdMode='cross', side='buy', posSide='long',
                             ordType='market', sz=sr1)

        # 设置杠杆倍数  Set Leverage
        accountAPI.set_leverage(instId=symbol, lever='50', mgnMode='cross')

        # account api
        accountAPI = Account.AccountAPI(api_key, secret_key, passphrase, False, flag)
        # 查看持仓信息  Get Positions
        result = accountAPI.get_positions('SWAP', symbol)

        time.sleep(5)

        # =====================================================================
        # print(dollar)
        # 持仓量
        if len(result['data']) > 0 and len(result['data'][0]) > 0 and len(
                result['data'][0]['notionalUsd']) > 0:
            notionalUsd = float(result['data'][0]['notionalUsd'])
            # 保证金
            imr = float(result['data'][0]['imr'])
            # 量
            pos = float(result['data'][0]['pos'])
            # 单个量 保证金值
            onlyimr = imr / pos
            # 每一美金值多少量
            onlyorder = int(float(dollar) / onlyimr)

            # 查看持仓信息  Get Positions
            result = accountAPI.get_positions('SWAP', symbol)

            if "low" == minute:
                onlyorder = int(onlyorder / 2)
            if "imr" == minute:
                onlyorder = int(onlyorder * 1.5)

            # 第2次买入
            tradeAPI = Trade.TradeAPI(api_key, secret_key, passphrase, False, flag)
            tradeAPI.place_order(instId=symbol, tdMode='cross', side='buy', posSide='long',
                                 ordType='market', sz=str(onlyorder))

            # 设置杠杆倍数  Set Leverage
            accountAPI.set_leverage(instId=symbol, lever='50', mgnMode='cross')

            time.sleep(1)

            # 查看持仓信息  Get Positions
            result = accountAPI.get_positions('SWAP', symbol)

            # 开仓均价
            avgPx = float(result['data'][0]['avgPx'])

            # 策略委托下单  Place Algo Order
            tradeAPI.place_algo_order(symbol, 'cross', 'sell', ordType='conditional',
                                      sz=sr1, posSide='long', tpTriggerPx=str(float(avgPx) * r),
                                      tpOrdPx=str(float(avgPx) * r))

            # 策略委托下单  Place Algo Order
            result = tradeAPI.place_algo_order(symbol, 'cross', 'sell', ordType='conditional',
                                               sz=str(onlyorder), posSide='long', tpTriggerPx=str(float(avgPx) * r),
                                               tpOrdPx=str(float(avgPx) * r))
            time.sleep(1)
            # 查看持仓信息  Get Positions
            result = accountAPI.get_positions('SWAP', symbol)

            # 开仓均价
            avgPx = float(result['data'][0]['avgPx'])

            # 策略委托下单  Place Algo Order
            result = tradeAPI.place_algo_order(symbol, 'cross', 'sell', ordType='conditional',
                                               sz=str(onlyorder), posSide='long', tpTriggerPx=str(float(avgPx) * r),
                                               tpOrdPx=str(float(avgPx) * r))

            time.sleep(1)

            # 查看持仓信息  Get Positions
            result = accountAPI.get_positions('SWAP', symbol)

            # 开仓均价
            avgPx = float(result['data'][0]['avgPx'])

            # 策略委托下单  Place Algo Order
            result = tradeAPI.place_algo_order(symbol, 'cross', 'sell', ordType='conditional',
                                               sz=str(onlyorder), posSide='long', tpTriggerPx=str(float(avgPx) * r),
                                               tpOrdPx=str(float(avgPx) * r))


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
