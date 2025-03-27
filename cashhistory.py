# -*- coding: gbk -*-
from mvc import MVC
from userinfo import User

if __name__ == '__main__':
    api_key, secret_key, passphrase, flag = User.get_userinfo()


    ##金历史记录 每分钟记录一次
    MVC.getcashhistory(api_key, secret_key, passphrase, flag)
