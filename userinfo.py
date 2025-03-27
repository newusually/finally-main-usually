# -*- coding: utf-8 -*-

import json


class User:
    # ��ȡ�û�API��Ϣ f'.'
    def get_userinfo():
        with open(f'../datas/api.json', 'r', encoding='utf-8') as f:
            obj = json.loads(f.read())

        api_key = obj['api_key']
        secret_key = obj['secret_key']
        passphrase = obj['passphrase']

        # flag��ʵ����ģ���̵��л����� flag is the key parameter which can help you to change between demo and real trading.
        # flag = '1'  # ģ���� demo trading
        flag = '0'  # ʵ�� real trading

        return api_key, secret_key, passphrase, flag

    # ��ȡ�û����˻�API��Ϣ f'./subAccount_api_key/'
    def get_subAccount_userinfo():
        subAccount_api_key = []
        subAccount_secret_key = []
        subAccount_passphrase = []

        for i in range(5):
            with open(f'../datas/subAccount_api_key/api0' + str(i + 1) + '.json', 'r', encoding='utf-8') as f:
                obj = json.loads(f.read())

            api_key = obj['api_key']
            secret_key = obj['secret_key']
            passphrase = obj['passphrase']

            subAccount_api_key.append(api_key)
            subAccount_secret_key.append(secret_key)
            subAccount_passphrase.append(passphrase)

            # flag��ʵ����ģ���̵��л����� flag is the key parameter which can help you to change between demo and real trading.
            # flag = '1'  # ģ���� demo trading
            flag = '0'  # ʵ�� real trading

        return subAccount_api_key, subAccount_secret_key, subAccount_passphrase, flag

    # ��ȡ�û��������� f'./
    def get_user_sr():
        with open(f'../datas/sr.json', 'r', encoding='utf-8') as f:
            obj = json.loads(f.read())

        sr = obj['sr']
        dollar = obj['dollar']
        dollar_eth = obj['dollar_eth']
        return sr, dollar, dollar_eth
