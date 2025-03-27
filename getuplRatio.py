# -*- coding: utf-8 -*-
from mvc import MVC
from userinfo import User

if __name__ == '__main__':
    api_key, secret_key, passphrase, flag = User.get_userinfo()

    MVC.getuplRatio_instId(api_key, secret_key, passphrase, flag)
