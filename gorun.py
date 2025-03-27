# -*- coding: gbk -*-
import sys
from datetime import datetime

from mvc import MVC
from userinfo import User

symbol = str(sys.argv[1])
minute = str(sys.argv[2])

if __name__ == '__main__':

    if (minute == '1m'):
        minute = '1'
    if (minute == '3m'):
        minute = '3'
    if (minute == '5m'):
        minute = '5'
    if (minute == '15m'):
        minute = '15'

    print("minute-->>", minute, "symbol-->>", symbol, "start-->>", datetime.now())

    api_key, secret_key, passphrase, flag = User.get_userinfo()
    MVC.orderbuy(api_key, secret_key, passphrase, flag, symbol, minute)
