# -*- coding: gbk -*-
from mvc import MVC
from userinfo import User

if __name__ == '__main__':
    api_key, secret_key, passphrase, flag = User.get_userinfo()

    MVC.sellall(api_key, secret_key, passphrase, flag)
