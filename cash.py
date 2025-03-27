# -*- coding: gbk -*-
from mvc import MVC
from userinfo import User

if __name__ == '__main__':
    api_key, secret_key, passphrase, flag = User.get_userinfo()

    #获取实时账户资金信息 每分钟查询一次
    MVC.getcashbal(api_key, secret_key, passphrase, flag)

